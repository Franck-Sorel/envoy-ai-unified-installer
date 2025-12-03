package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/config"
)

var (
	cfgFile    string
	dryRun     bool
	skipClean  bool
	verbose    bool
	namespaceGW string
	namespaceAI string
)

var rootCmd = &cobra.Command{
	Use:   "envoy-ai-installer",
	Short: "Unified installer for Envoy Gateway + Envoy AI Gateway",
	Long: `envoy-ai-installer is a production-grade CLI for installing Envoy AI Gateway.

It automatically fetches the latest upstream releases and provides
a seamless installation experience with sensible defaults and
full customization options.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(cfgFile); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", 
		"config file (default is $HOME/.envoy-ai-installer/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false,
		"simulate what would be executed without making changes")
	rootCmd.PersistentFlags().BoolVar(&skipClean, "skip-clean", false,
		"skip cleaning up previous installations")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose output")
	rootCmd.PersistentFlags().StringVar(&namespaceGW, "namespace-gateway", "envoy-gateway-system",
		"kubernetes namespace for Envoy Gateway")
	rootCmd.PersistentFlags().StringVar(&namespaceAI, "namespace-ai", "envoy-ai-gateway-system",
		"kubernetes namespace for Envoy AI Gateway")

	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))
	viper.BindPFlag("skip_clean", rootCmd.PersistentFlags().Lookup("skip-clean"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("namespace_gateway", rootCmd.PersistentFlags().Lookup("namespace-gateway"))
	viper.BindPFlag("namespace_ai", rootCmd.PersistentFlags().Lookup("namespace-ai"))

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(fmt.Sprintf("%s/.envoy-ai-installer", home))
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("EAIG")
	viper.AutomaticEnv()
}

func Execute() error {
	return rootCmd.Execute()
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
