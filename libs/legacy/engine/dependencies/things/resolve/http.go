package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	sqldb "github.com/ongruent/gruent/db/sqlc"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/secrets"
)

const (
	// This network is the one spawn up with the docker compose file
	// in the local environment at deploy/engine/local/docker-compose.yaml
	LocalEnvDockerNetworkName = "gruent_network"
)

func resolveHttpServerDependency(
	env environment.Env,
	thing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = secretsClient
	_ = spaceId

	// Unmarshal the thing so that we can process it
	var server things.HttpServer
	err := json.Unmarshal(thing.Attributes, &server)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	if thing.Internal {
		switch env {
		case environment.Development:
			if err := handleDevEnvServer(&server); err != nil {
				return fmt.Errorf("could not handle local server setup: %v", err)
			}
		// case environment.Local:
		// 	if err := handleDockerServer(&server); err != nil {
		// 		return fmt.Errorf("could not handle docker container setup: %v", err)
		// 	}
		// case environment.Production, environment.Staging:
		// 	if err := handleK8sServer(&server); err != nil {
		// 		return fmt.Errorf("could not handle k8s server setup: %v", err)
		// 	}
		default:
			return fmt.Errorf("unknown environment: %v", env)
		}
	}

	// Serialize the updated server configuration
	updatedAttributes, err := json.Marshal(server)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}
	thing.Attributes = updatedAttributes

	// Set the corresponding config field to this server
	nodeConfig[key] = server

	return nil
}

func handleDevEnvServer(server *things.HttpServer) error {
	// update server port
	port, err := getFreePort()
	if err != nil {
		return fmt.Errorf("could not find a free port: %v", err)
	}
	server.Port = port // Update the server's port with the free port

	// update host, protocol, and paths
	server.Host = "localhost"
	server.Protocol = "http"
	server.Paths = []string{"/"}

	return nil
}

func handleDockerServer(server *things.HttpServer) error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	containerID, err := getOwnContainerId()
	if err != nil {
		return fmt.Errorf("failed to get own container ID: %v", err)
	}

	// Inspect our own container
	container, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %v", err)
	}

	// Get the container's IP from the runtime_network
	var containerIP string
	if networkSettings, exists := container.NetworkSettings.Networks[LocalEnvDockerNetworkName]; exists {
		containerIP = networkSettings.IPAddress
	} else if networkSettings, exists := container.NetworkSettings.Networks["bridge"]; exists {
		containerIP = networkSettings.IPAddress
	} else {
		return fmt.Errorf("container not connected to expected network")
	}

	// Check if this is a runtime container and get its port from labels
	if portStr, exists := container.Config.Labels["com.iotea.port"]; exists {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port in container labels: %v", err)
		}
		server.Port = port
	} else {
		// Not a runtime container, use default port.
		// This could happen when the orchestrator is trying to check if
		// dependencies can be resolved
		server.Port = 10000
	}

	// Update the server configuration
	server.Host = containerIP
	server.Protocol = "http"
	server.Paths = []string{"/"}

	return nil
}

func handleK8sServer(server *things.HttpServer) error {
	// Get the pod's IP address
	podIP, err := getPodIP()
	if err != nil {
		return fmt.Errorf("failed to get pod IP: %v", err)
	}

	// Find a free port within the pod
	port, err := getFreePort()
	if err != nil {
		return fmt.Errorf("could not find a free port: %v", err)
	}

	// Update the server configuration
	server.Host = podIP
	server.Port = port
	server.Protocol = "http"
	server.Paths = []string{"/"}

	return nil
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func getPodIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("error getting network interfaces: %v", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("error getting addresses for interface: %v", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no IP address found")
}

func getOwnContainerId() (string, error) {
	// Try to get container ID from hostname first
	hostname, err := os.Hostname()
	if err == nil && len(hostname) == 12 {
		return hostname, nil
	}

	// Fallback to checking cgroup file
	content, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "", err
	}

	// Look for container ID in cgroup content
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		// Check for both docker and containerd patterns
		if strings.Contains(line, "docker") || strings.Contains(line, "containerd") {
			parts := strings.Split(line, "/")
			if len(parts) > 2 {
				return parts[len(parts)-1], nil
			}
		}
	}

	return "", fmt.Errorf("container ID not found in hostname or cgroup")
}
