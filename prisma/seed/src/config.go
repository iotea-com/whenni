package seed

import (
	"github.com/spf13/viper"
)

type Config struct {
	DirectUrl     string `mapstructure:"DIRECT_URL"`
	DatabaseUrl   string `mapstructure:"DATABASE_URL"`
	DatabaseToken string `mapstructure:"DATABASE_TOKEN"`
	AssetsFolder  string `mapstructure:"ASSETS_FOLDER"`

	// Vault
	VaultAddress  string `mapstructure:"VAULT_ADDRESS"`
	VaultUsername string `mapstructure:"VAULT_USERNAME"`
	VaultPassword string `mapstructure:"VAULT_PASSWORD"`
}

func loadConfig() (config *Config) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.BindEnv("DIRECT_URL")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("DATABASE_TOKEN")
	viper.BindEnv("ASSETS_FOLDER")

	viper.BindEnv("VAULT_ADDRESS")
	viper.BindEnv("VAULT_USERNAME")
	viper.BindEnv("VAULT_PASSWORD")

	// Put the read config into our app config struct
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	return
}

var Conf *Config

func init() {
	Conf = loadConfig()
}
