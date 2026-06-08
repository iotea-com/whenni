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
	// Environment for testing with mock services, stubs, and fixtures.
	Test Env = "test"

	// Environment for local development.
	Development Env = "development"

	// Environment for built container images deployed to a k8s cluster.
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
