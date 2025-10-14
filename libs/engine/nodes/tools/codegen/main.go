package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type NodePort struct {
	Id    string `json:"id" validate:"required"`
	Label string `json:"label" validate:"required"`
}

type NodeManifest struct {
	Name        string     `json:"name" validate:"required"`
	Label       string     `json:"label" validate:"required"`
	Version     string     `json:"version" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Category    string     `json:"category" validate:"required"`
	Inputs      []NodePort `json:"inputs" validate:"required,dive"`
	Outputs     []NodePort `json:"outputs" validate:"required,dive"`
	Runtime     string     `json:"runtime" validate:"required,oneof=wasm jvm"`
}

type NodeRef struct {
	Name     string
	Category string
	Version  string
	Schema   string // absolute path to config.schema.json
	BaseDir  string // .../plugins/<category>/<name>/<version>
}

func main() {
	// Ensure quicktype is available
	if _, err := exec.LookPath("quicktype"); err != nil {
		log.Fatalf("quicktype not found on PATH - run `npm i -g quicktype` to install")
	}

	// Run the codegen
	if err := run(); err != nil {
		log.Fatalf("codegen error: %v\n", err)
	}
}

func run() error {
	v := validator.New()

	// Discover nodes (nodes/plugins/<category>/<name>/<version>/node.json)
	root, _ := os.Getwd()
	nodesDir := filepath.Join(root, "../../plugins")

	var refs []NodeRef

	err := filepath.WalkDir(nodesDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(p) != "node.json" {
			return nil
		}

		// Read and parse the node.json file
		base := filepath.Dir(p)
		data, readErr := os.ReadFile(p)
		if readErr != nil {
			return readErr
		}

		var nodeManifest NodeManifest
		if jerr := json.Unmarshal(data, &nodeManifest); jerr != nil {
			return fmt.Errorf("parse %s: %w", p, jerr)
		}

		// Validate the node manifest
		if err := v.Struct(nodeManifest); err != nil {
			return fmt.Errorf("validate %s: %w", p, err)
		}

		schemaPath := filepath.Join(filepath.Dir(p), "config.schema.json")
		if !filepath.IsAbs(schemaPath) {
			schemaPath = filepath.Join(base, schemaPath)
		}
		refs = append(refs, NodeRef{
			Name:     toSafeName(nodeManifest.Name),
			Category: nodeManifest.Category,
			Version:  nodeManifest.Version,
			Schema:   schemaPath,
			BaseDir:  base,
		})
		return nil
	})
	if err != nil {
		return err
	}

	if len(refs) == 0 {
		return errors.New("no nodes discovered (expected nodes/<category>/<name>/<version>/node.json)")
	}

	// Generate types for each node
	targets := defaultTargets()

	for _, ref := range refs {
		log.Printf("→ Generating types for %s/%s@%s\n", ref.Category, ref.Name, ref.Version)
		if err := genOne(ref, targets.TS); err != nil {
			return err
		}
		if err := genOne(ref, targets.Go); err != nil {
			return err
		}
		if err := genOne(ref, targets.Kotlin); err != nil {
			return err
		}
		// if err := genOne(ref, targets.Rust); err != nil {
		// 	return err
		// }
	}

	// Write a registry index
	indexPath := filepath.Join(nodesDir, "registry.index.json")
	if err := writeRegistryIndex(indexPath, refs); err != nil {
		return err
	}

	log.Println("✓ Node codegen complete.")
	return nil
}

func genOne(ref NodeRef, t LangTarget) error {
	outFilename := strings.ReplaceAll(t.FilePattern, "{name}", fileBaseFor(ref))
	outPath := filepath.Join(ref.BaseDir, outFilename)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	args := []string{
		"--src", ref.Schema,
		"--src-lang", "schema",
		"--out", outPath,
		"--lang", t.Lang,
	}
	args = append(args, t.Extras...)

	cmd := exec.Command("quicktype", args...)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("quicktype %s %s: %v\nstderr: %s", ref.Name, t.Lang, err, cmd.Stderr.(*bytes.Buffer).String())
	}
	return nil
}

// Creates a relevant filename for the node, e.g. conditional_threshold_v1_0_0 (<category>_<name>_<version>)
func fileBaseFor(ref NodeRef) string {
	category := toSafeName(ref.Category)
	name := toSafeName(ref.Name)
	version := "v" + strings.ReplaceAll(ref.Version, ".", "_")
	return category + "_" + name + "_" + version
}

var safeFilenamePattern = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func toSafeName(s string) string {
	return safeFilenamePattern.ReplaceAllString(s, "")
}

func writeRegistryIndex(path string, refs []NodeRef) error {
	type idxEntry struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Manifest string `json:"manifest"`
		Schema   string `json:"schema"`
	}
	var entries []idxEntry
	for _, r := range refs {
		entries = append(entries, idxEntry{
			Name: r.Name, Version: r.Version,
			Manifest: filepath.Join(r.Name, "node.json"),
			Schema:   r.Schema,
		})
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	buf, _ := json.MarshalIndent(entries, "", "  ")
	return os.WriteFile(path, buf, 0o644)
}
