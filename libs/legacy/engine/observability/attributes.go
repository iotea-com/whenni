package observability

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"runtime"

	"github.com/iotea-com/iotea/libs/engine/environment"
	"github.com/rs/zerolog/log"

	"github.com/docker/docker/client"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Refer to https://opentelemetry.io/docs/specs/semconv/resource/
func setAttributes(env environment.Env, resourceAttr ResourceAttributes) []attribute.KeyValue {
	// Initialize base resource attributes from the user
	attrs := []attribute.KeyValue{
		// Environment
		semconv.DeploymentEnvironmentKey.String(env.String()),

		// Service information
		semconv.ServiceNameKey.String(resourceAttr.ServiceName),
		semconv.ServiceInstanceIDKey.String(resourceAttr.ServiceInstance),

		// Build information
		attribute.Key("build.version").String(resourceAttr.BuildVersion),
		attribute.Key("build.time").String(resourceAttr.BuildTime),
		attribute.Key("build.commit").String(resourceAttr.BuildCommit),
		attribute.Key("build.dirty").String(resourceAttr.BuildDirty),
		attribute.Key("build.creator").String(resourceAttr.BuildCreator),
	}

	// Get host information
	hostname, err := os.Hostname()
	if err == nil {
		attrs = append(attrs, semconv.HostNameKey.String(hostname))
	}
	attrs = append(attrs, semconv.HostArchKey.String(runtime.GOARCH))

	// Add the go version
	attrs = append(attrs, attribute.Key("runtime.version").String(runtime.Version()))

	// Environment Specific Attributes
	switch env {
	case environment.Development:
		if err := setContainerAttributes(&attrs); err != nil {
			log.Warn().Msgf("Failed to set container attributes: %v", err)
		}
	case environment.Local:
		if err := setContainerAttributes(&attrs); err != nil {
			log.Warn().Msgf("Failed to set container attributes: %v", err)
		}
	case environment.Staging, environment.Production:
		// Try Kubernetes first
		if err := setK8sAttributes(&attrs); err != nil {
			log.Warn().Msgf("Failed to set Kubernetes attributes: %v", err)
		}

	default:
		log.Warn().Msgf("Unknown environment: %s", env)
	}

	return attrs
}

func setContainerAttributes(attrs *[]attribute.KeyValue) error {
	// Try to get Docker client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer dockerClient.Close()

	// Get current container ID from hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Inspect current container with timeout context
	containerInfo, err := dockerClient.ContainerInspect(ctx, hostname)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}

	// Set container attributes
	containerAttrs := []attribute.KeyValue{
		semconv.ContainerIDKey.String(containerInfo.ID),
		semconv.ContainerNameKey.String(strings.TrimPrefix(containerInfo.Name, "/")),
		attribute.Key("container.runtime").String("docker"),
	}

	// Get image information
	if containerInfo.Config != nil {
		imageParts := strings.Split(containerInfo.Config.Image, ":")
		containerAttrs = append(containerAttrs,
			semconv.ContainerImageNameKey.String(imageParts[0]),
		)
		if len(imageParts) > 1 {
			containerAttrs = append(containerAttrs,
				semconv.ContainerImageTagKey.String(imageParts[1]),
			)
		}
	}

	// Get network information
	if containerInfo.NetworkSettings != nil {
		// First try the default network
		if containerInfo.NetworkSettings.IPAddress != "" {
			containerAttrs = append(containerAttrs,
				attribute.Key("container.ip").String(containerInfo.NetworkSettings.IPAddress),
			)
		} else {
			// If not found, check all networks
			for _, network := range containerInfo.NetworkSettings.Networks {
				if network.IPAddress != "" {
					containerAttrs = append(containerAttrs,
						attribute.Key("container.ip").String(network.IPAddress),
					)
					break
				}
			}
		}
	}

	*attrs = append(*attrs, containerAttrs...)
	return nil
}

func setK8sAttributes(attrs *[]attribute.KeyValue) error {
	// Get pod name from hostname
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		return fmt.Errorf("could not determine pod name")
	}

	// Get namespace from the service account mount
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return fmt.Errorf("failed to read namespace: %w", err)
	}
	ns := string(namespace)

	// Get basic pod metadata from downward API environment variables
	// These need to be set in your pod spec:
	// https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/
	nodeName := os.Getenv("NODE_NAME")
	podIP := os.Getenv("POD_IP")

	k8sAttrs := []attribute.KeyValue{
		semconv.K8SNamespaceNameKey.String(ns),
		semconv.K8SPodNameKey.String(podName),
		attribute.Key("container.runtime").String("k8s"),
	}

	if nodeName != "" {
		k8sAttrs = append(k8sAttrs, semconv.K8SNodeNameKey.String(nodeName))
	}

	if podIP != "" {
		k8sAttrs = append(k8sAttrs, attribute.Key("k8s.pod.ip").String(podIP))
	}

	// Get container info from env vars
	if containerName := os.Getenv("CONTAINER_NAME"); containerName != "" {
		k8sAttrs = append(k8sAttrs, semconv.ContainerNameKey.String(containerName))
	}
	if image := os.Getenv("CONTAINER_IMAGE"); image != "" {
		parts := strings.Split(image, ":")
		k8sAttrs = append(k8sAttrs, semconv.ContainerImageNameKey.String(parts[0]))
		if len(parts) > 1 {
			k8sAttrs = append(k8sAttrs, semconv.ContainerImageTagKey.String(parts[1]))
		}
	}

	*attrs = append(*attrs, k8sAttrs...)
	return nil
}
