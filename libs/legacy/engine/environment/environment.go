package environment

import "os"

const (
	// This is the environment variable that will be used to determine the environment
	// of the application. All services should be able to read this environment variable
	// and use it to determine the environment they are running in.
	EnvVarName = "ENVIRONMENT"
)

// Env is the type for the environment of the application.
type Env string

const (
	// Test is the environment for testing, where the application is running in a
	// test environment, mock clients are used, and the application is not actually
	// running.
	Test Env = "test"

	// Development is the environment for development, where all the rules engine
	// services are running in a single container instance, and leverage the "air"
	// tool to automatically restart the container when code changes are detected.
	Development Env = "development"

	// Local is the environment for local testing, after code changes are
	// made. This will spin up individual containers for each service, and a new dedicated
	// runtime container will be spun up for each channel.
	Local Env = "local"

	// Staging is the environment for staging, where the rules engine is
	// deployed to a kubernetes cluster, and each channel is deployed as a separate
	// pod.
	Staging Env = "staging"

	// Production is the environment for production, where the rules engine is
	// deployed to a kubernetes cluster, and each channel is deployed as a separate
	// pod.
	Production Env = "production"
)

// GetFromEnvVar returns the environment from the environment variable.
func GetFromEnvVar() Env {
	return Env(os.Getenv(EnvVarName))
}

// String returns the string representation of the environment.
func (e Env) String() string {
	return string(e)
}
