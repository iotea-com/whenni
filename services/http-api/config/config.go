package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/libs/legacy/engine/observability"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
	clickhouseService "github.com/iotea-com/iotea/services/http-api/services/clickhouse"
	postgresService "github.com/iotea-com/iotea/services/http-api/services/prisma"
	ioteaSmtp "github.com/iotea-com/iotea/services/http-api/services/smtp"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Build time variables
const (
	BuildVersion = "{{BUILD_VERSION}}"
	BuildCommit  = "{{BUILD_COMMIT}}"
	BuildDirty   = "{{BUILD_DIRTY}}"
	BuildTime    = "{{BUILD_TIME}}"
	BuildCreator = "{{BUILD_CREATOR}}"
)

// Global app config
var Env environment.Env
var EnvConf *EnvConfig
var VaultConf *VaultConfig
var SecretsClient secrets.SecretsClient

// EnvConfig is the configuration we get from the environment variables
type EnvConfig struct {
	ApiVaultUsername string `mapstructure:"API_VAULT_USERNAME"`
	ApiVaultPassword string `mapstructure:"API_VAULT_PASSWORD"`
}

// VaultConfig is the configuration we get from the Vault secrets
type VaultConfig struct {
	// Logging
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// API server port
	ApiPort int `mapstructure:"API_SERVER_PORT"`

	// Orchestrator gRPC end point
	EngineGrpcServiceUrl string `mapstructure:"ENGINE_GRPC_SERVER_URL"`

	// Database
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	JwtSecret   string `mapstructure:"JWT_SECRET"`

	// Observability
	PlatformOtelCollectorEndpoint string `mapstructure:"PLATFORM_OTEL_COLLECTOR_ENDPOINT"`

	// Dev Environments
	DevenvGrpcServerUrl string `mapstructure:"DEVENV_GRPC_SERVER_URL"`

	// Runtime Metrics Collector
	CollectorGrpcRuntimeUrl string `mapstructure:"COLLECTOR_GRPC_RUNTIME_URL"`

	// Clickhouse
	ClickhouseHost     string `mapstructure:"CLICKHOUSE_HOST"`
	ClickhousePort     int    `mapstructure:"CLICKHOUSE_PORT"`
	ClickhouseUsername string `mapstructure:"CLICKHOUSE_USERNAME"`
	ClickhousePassword string `mapstructure:"CLICKHOUSE_PASSWORD"`
	ClickhouseDatabase string `mapstructure:"CLICKHOUSE_DATABASE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// SMTP
	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtpPort     int    `mapstructure:"SMTP_PORT"`
	SmtpUser     string `mapstructure:"SMTP_USER"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
}

func init() {
	Env = environment.GetFromEnvVar()
	EnvConf, VaultConf = load(Env)
}

func load(env environment.Env) (envConfig *EnvConfig, vaultConfig *VaultConfig) {
	var err error

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Vault address from the environment variables
	vaultAddress := os.Getenv("VAULT_ADDRESS")

	if env == environment.Test {
		VaultConf = &VaultConfig{
			JwtSecret: "super-secret-jwt-token-with-at-least-32-characters-long",
		}
		return
	} else if env == environment.Development || env == environment.Local {
		viper.BindEnv("API_VAULT_USERNAME")
		viper.BindEnv("API_VAULT_PASSWORD")

		// Put the read config into our app config struct
		if err := viper.Unmarshal(&envConfig); err != nil {
			panic(err)
		}

		// In a local setup, we use UserPass auth method to authenticate
		// with vault, and get a token that way. This token is valid for the
		// duration of the app
		SecretsClient, err = secrets.NewClient(vaultAddress, envConfig.ApiVaultUsername, envConfig.ApiVaultPassword, "")
		if err != nil {
			panic(fmt.Sprintf("Could not connect to Vault, %v", err))
		}
	} else if env == environment.Staging || env == environment.Production {
		// In a production setup, Vault is setup to "inject" a token into
		// the container this application is running on, so we don't need
		// to authenticate
		SecretsClient, err = secrets.NewClient(vaultAddress, "", "", secrets.VaultTokenFilePath)
		if err != nil {
			panic(fmt.Sprintf("Could not connect to Vault, %v", err))
		}
	} else {
		panic(fmt.Sprintf("Unknown or empty environment env var: %v", env))
	}

	// Load secrets from Vault into the VaultConfig struct
	vaultConfig, err = loadVaultSecrets()
	if err != nil {
		panic(fmt.Sprintf("Error loading Vault secrets: %v", err))
	}

	// Print the info banner only for the main process, since fiber http framework
	// spawns multiple processes if prefork is enabled
	if os.Getppid() <= 1 {
		info := environment.BuildInfo{
			ServiceName: "HTTP API",
			Version:     BuildVersion,
			LogLevel:    vaultConfig.LogLevel,
			CommitHash:  BuildCommit,
			BuildTime:   BuildTime,
			Dirty:       BuildDirty,
			Creator:     BuildCreator,
			ExtraFields: map[string]any{
				"Vault Address": vaultAddress,
				"API Port":      vaultConfig.ApiPort,
			},
		}
		fmt.Println(environment.PrintInfoBanner(info))
	}

	// If we are running in a development environment, we need to set the engine grpc server url to the local address.
	// This is because the orchestrator is running in the same container as the http-api.
	//
	// In all other environments, the orchestrator is running in a different container,
	// so we don't need to do this.
	if env == environment.Development {
		// Strip the port from the engine grpc server url
		engineGrpcServerUrl := strings.Split(vaultConfig.EngineGrpcServiceUrl, ":")
		vaultConfig.EngineGrpcServiceUrl = fmt.Sprintf("127.0.0.1:%s", engineGrpcServerUrl[1])
	}

	// There's a bug where the prisma client only picks up the DATABASE_URL
	// environment variable if it's set before the prisma client is instantiated.
	// So we set it here.
	os.Setenv("DATABASE_URL", vaultConfig.DatabaseUrl)

	// Initialize the DB client
	err = initDbClient(vaultConfig.DatabaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Error initializing database client: %v", err))
	}

	// Initialize the Clickhouse client
	err = initClickhouseClient(vaultConfig.ClickhouseHost, vaultConfig.ClickhousePort, vaultConfig.ClickhouseUsername, vaultConfig.ClickhousePassword, vaultConfig.ClickhouseDatabase)
	if err != nil {
		panic(fmt.Sprintf("Error initializing Clickhouse client: %v", err))
	}

	// Initialize the SMTP client
	err = initSmtpClient(vaultConfig.SmtpHost, vaultConfig.SmtpPort, vaultConfig.SmtpUser, vaultConfig.SmtpPassword)
	if err != nil {
		panic(fmt.Sprintf("Error initializing SMTP client: %v", err))
	}

	// Setup logging
	zerolog.SetGlobalLevel(observability.LogLevel(vaultConfig.LogLevel).Level())

	return
}

func loadVaultSecrets() (*VaultConfig, error) {
	// Retrieve the secrets from Vault
	secrets, err := SecretsClient.ReadKeyValue(context.Background(), "config/http-api")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secrets from Vault: %w", err)
	}

	// Safely assign each value, with a fallback to empty string or handle nil cases
	vaultConfig := &VaultConfig{
		LogLevel:                      SecretsClient.GetStringFromMap(secrets, "LOG_LEVEL"),
		ApiPort:                       SecretsClient.GetIntFromMap(secrets, "API_SERVER_PORT"),
		EngineGrpcServiceUrl:          SecretsClient.GetStringFromMap(secrets, "ENGINE_GRPC_SERVER_URL"),
		DatabaseUrl:                   SecretsClient.GetStringFromMap(secrets, "DATABASE_URL"),
		JwtSecret:                     SecretsClient.GetStringFromMap(secrets, "JWT_SECRET"),
		PlatformOtelCollectorEndpoint: SecretsClient.GetStringFromMap(secrets, "PLATFORM_OTEL_COLLECTOR_ENDPOINT"),
		CollectorGrpcRuntimeUrl:       SecretsClient.GetStringFromMap(secrets, "COLLECTOR_GRPC_RUNTIME_URL"),
		DevenvGrpcServerUrl:           SecretsClient.GetStringFromMap(secrets, "DEVENV_GRPC_SERVER_URL"),
		ClickhouseHost:                SecretsClient.GetStringFromMap(secrets, "CLICKHOUSE_HOST"),
		ClickhousePort:                SecretsClient.GetIntFromMap(secrets, "CLICKHOUSE_PORT"),
		ClickhouseUsername:            SecretsClient.GetStringFromMap(secrets, "CLICKHOUSE_USERNAME"),
		ClickhousePassword:            SecretsClient.GetStringFromMap(secrets, "CLICKHOUSE_PASSWORD"),
		ClickhouseDatabase:            SecretsClient.GetStringFromMap(secrets, "CLICKHOUSE_DATABASE"),
		RedisHost:                     SecretsClient.GetStringFromMap(secrets, "REDIS_HOST"),
		RedisPort:                     SecretsClient.GetIntFromMap(secrets, "REDIS_PORT"),
		RedisUsername:                 SecretsClient.GetStringFromMap(secrets, "REDIS_USERNAME"),
		RedisPassword:                 SecretsClient.GetStringFromMap(secrets, "REDIS_PASSWORD"),
		SmtpHost:                      SecretsClient.GetStringFromMap(secrets, "SMTP_HOST"),
		SmtpPort:                      SecretsClient.GetIntFromMap(secrets, "SMTP_PORT"),
		SmtpUser:                      SecretsClient.GetStringFromMap(secrets, "SMTP_USER"),
		SmtpPassword:                  SecretsClient.GetStringFromMap(secrets, "SMTP_PASSWORD"),
	}

	return vaultConfig, nil
}

func initDbClient(databaseUrl string) error {
	// If running 'nx test', use a default mock client
	environment := os.Getenv("ENVIRONMENT")
	if environment == "test" {
		client, _, _ := db.NewMock()
		postgresService.Client = client
		return nil
	}

	// Init prisma client
	client := db.NewClient()

	// Connect
	if err := client.Prisma.Connect(); err != nil {
		return fmt.Errorf("could not connect to the database: %s", err)
	}

	postgresService.Client = client
	return nil
}

func initClickhouseClient(databaseHost string, databasePort int, databaseUsername string, databasePassword string, databaseName string) error {
	// If running 'nx test', use a default mock client
	// environment := os.Getenv("ENVIRONMENT")
	// if environment == "test" {
	// 	client, _, _ := db.NewMock()
	// 	appsecrets.Client = client
	// 	return nil
	// }

	// Connect to ClickHouse
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr:     []string{fmt.Sprintf("%s:%d", databaseHost, databasePort)},
		Protocol: clickhouse.Native,
		Auth: clickhouse.Auth{
			Database: databaseName,
			Username: databaseUsername,
			Password: databasePassword,
		},
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	})
	if err != nil {
		return fmt.Errorf("error connecting to ClickHouse: %s", err)
	}
	defer conn.Close()

	clickhouseService.Conn = conn
	return nil
}

func initSmtpClient(smtpHost string, smtpPort int, smtpUser string, smtpPassword string) error {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "test" {
		smtpClient := ioteaSmtp.NewMockSmtpClient()
		ioteaSmtp.SmtpClient = smtpClient
		return nil
	}

	if environment == "development" || environment == "local" {
		smtpClient := ioteaSmtp.NewSmtpClient(smtpHost, smtpPort, nil, nil)
		ioteaSmtp.SmtpClient = smtpClient
		return nil
	}

	smtpClient := ioteaSmtp.NewSmtpClient(smtpHost, smtpPort, &smtpUser, &smtpPassword)
	ioteaSmtp.SmtpClient = smtpClient
	return nil
}
