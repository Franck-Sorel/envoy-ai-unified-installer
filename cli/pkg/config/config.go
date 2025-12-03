package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	NamespaceGateway string
	NamespaceAI      string
	SkipClean        bool
	DryRun           bool
	ValuesExtra      []string
}

func Init(configPath string) error {
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			configDir := filepath.Join(home, ".envoy-ai-installer")
			viper.AddConfigPath(configDir)
		}
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("EAIG")
	viper.AutomaticEnv()

	viper.SetDefault("namespace_gateway", "envoy-gateway-system")
	viper.SetDefault("namespace_ai", "envoy-ai-gateway-system")
	viper.SetDefault("skip_clean", false)
	viper.SetDefault("dry_run", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

func Load() *Config {
	return &Config{
		NamespaceGateway: viper.GetString("namespace_gateway"),
		NamespaceAI:      viper.GetString("namespace_ai"),
		SkipClean:        viper.GetBool("skip_clean"),
		DryRun:           viper.GetBool("dry_run"),
		ValuesExtra:      viper.GetStringSlice("values_extra"),
	}
}

func SetDefaults(namespace, namespaceAI string, skipClean, dryRun bool, valuesExtra []string) {
	if namespace != "" {
		viper.Set("namespace_gateway", namespace)
	}
	if namespaceAI != "" {
		viper.Set("namespace_ai", namespaceAI)
	}
	viper.Set("skip_clean", skipClean)
	viper.Set("dry_run", dryRun)
	if len(valuesExtra) > 0 {
		viper.Set("values_extra", valuesExtra)
	}
}
