package main

type LangTarget struct {
	Lang        string // quicktype --lang
	FilePattern string // output file pattern with {name} substitution
	Extras      []string
}

type Targets struct {
	TS     LangTarget
	Go     LangTarget
	Kotlin LangTarget
	// Rust   LangTarget
}

func defaultTargets() Targets {
	return Targets{
		TS: LangTarget{
			Lang:        "ts",
			FilePattern: "{name}.ts",
		},
		Go: LangTarget{
			Lang:        "go",
			FilePattern: "{name}.go",
			Extras:      []string{"--package", "codegen"},
		},
		Kotlin: LangTarget{
			Lang:        "kotlin",
			FilePattern: "{name}.kt",
		},
		// Rust: LangTarget{
		// 	Lang:       "rust",
		// 	OutPattern: "services/runtime-rs/src/nodes/{name}_cfg.rs",
		// },
	}
}
