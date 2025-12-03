package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/config"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/helm"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/upstream"
)

var (
	valuesExtra string
	withRedis   bool
	chartRepo   string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Envoy AI Gateway with upstream charts",
	Long: `Install Envoy AI Gateway by fetching the latest upstream releases.

This command implements the official 4-step installation process:
1. Clean previous installations (unless --skip-clean)
2. Install Envoy Gateway with official values
3. Install Envoy AI Gateway CRDs
4. Install Envoy AI Gateway controller

All steps support customization via flags and config files.`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringVar(&valuesExtra, "values-extra", "",
		"comma-separated list of additional values files to use")
	installCmd.Flags().BoolVar(&withRedis, "with-redis", false,
		"install Redis for rate limiting (optional)")
	installCmd.Flags().StringVar(&chartRepo, "chart-repo", "",
		"optional pre-built chart repository URL")

	viper.BindPFlag("values_extra", installCmd.Flags().Lookup("values-extra"))
	viper.BindPFlag("with_redis", installCmd.Flags().Lookup("with-redis"))
}

func runInstall(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	isDryRun := viper.GetBool("dry_run")
	isVerbose := viper.GetBool("verbose")

	fmt.Println("🚀 Envoy AI Gateway Installer")
	fmt.Printf("  Namespace (Gateway): %s\n", cfg.NamespaceGateway)
	fmt.Printf("  Namespace (AI):      %s\n", cfg.NamespaceAI)
	fmt.Printf("  Dry Run:             %v\n", isDryRun)

	if !cfg.SkipClean {
		fmt.Println("\n📋 Step 1/4: Cleaning up previous installations...")
		if err := cleanPreviousInstall(cfg, isDryRun); err != nil {
			return fmt.Errorf("cleanup failed: %w", err)
		}
	}

	helmCmd := helm.NewHelmCommand(isDryRun)

	fmt.Println("\n📋 Step 2/4: Installing Envoy Gateway...")
	if err := installEnvoyGateway(helmCmd, cfg); err != nil {
		return fmt.Errorf("failed to install Envoy Gateway: %w", err)
	}

	fmt.Println("\n📋 Step 3/4: Installing Envoy AI Gateway CRDs...")
	if err := installAIGatewayCRDs(helmCmd, cfg); err != nil {
		return fmt.Errorf("failed to install AI Gateway CRDs: %w", err)
	}

	fmt.Println("\n📋 Step 4/4: Installing Envoy AI Gateway controller...")
	if err := installAIGatewayController(helmCmd, cfg); err != nil {
		return fmt.Errorf("failed to install AI Gateway controller: %w", err)
	}

	if withRedis {
		fmt.Println("\n📦 Installing Redis for rate limiting...")
		if err := installRedis(helmCmd, cfg); err != nil {
			return fmt.Errorf("failed to install Redis: %w", err)
		}
	}

	fmt.Println("\n✅ Installation complete!")
	if isDryRun {
		fmt.Println("   This was a dry run. Use 'envoy-ai-installer install' without --dry-run to execute.")
	} else {
		fmt.Printf("   Verify installation: kubectl get pods -n %s\n", cfg.NamespaceGateway)
	}

	return nil
}

func cleanPreviousInstall(cfg *config.Config, isDryRun bool) error {
	helmCmd := helm.NewHelmCommand(isDryRun)

	releases := []struct {
		name      string
		namespace string
	}{
		{"eg", cfg.NamespaceGateway},
		{"aieg-crd", cfg.NamespaceAI},
		{"aieg", cfg.NamespaceAI},
	}

	for _, r := range releases {
		if err := helmCmd.Uninstall(r.name, r.namespace); err != nil {
			fmt.Printf("  Note: %s was not previously installed\n", r.name)
		}
	}

	return nil
}

func installEnvoyGateway(helmCmd *helm.HelmCommand, cfg *config.Config) error {
	if err := helmCmd.RepoAdd("envoyproxy", "oci://docker.io/envoyproxy"); err != nil {
		return err
	}

	if err := helmCmd.RepoUpdate(); err != nil {
		return err
	}

	valuesFile, err := fetchRemoteValuesFile(
		"https://raw.githubusercontent.com/envoyproxy/ai-gateway/main/manifests/envoy-gateway-values.yaml",
	)
	if err != nil {
		fmt.Printf("Warning: Could not fetch official values file: %v\n", err)
		valuesFile = ""
	}

	values := []string{}
	if valuesFile != "" {
		values = append(values, valuesFile)
	}

	if valuesExtra != "" {
		extraValues := strings.Split(valuesExtra, ",")
		for _, v := range extraValues {
			v = strings.TrimSpace(v)
			if v != "" {
				values = append(values, v)
			}
		}
	}

	opts := &helm.HelmOptions{
		DryRun:    false,
		Namespace: cfg.NamespaceGateway,
		Values:    values,
		Version:   "v0.0.0-latest",
	}

	return helmCmd.Install("eg", "envoyproxy/gateway-helm", cfg.NamespaceGateway, opts)
}

func installAIGatewayCRDs(helmCmd *helm.HelmCommand, cfg *config.Config) error {
	if err := helmCmd.RepoAdd("envoyproxy-ai", "oci://docker.io/envoyproxy"); err != nil {
		return err
	}

	if err := helmCmd.RepoUpdate(); err != nil {
		return err
	}

	opts := &helm.HelmOptions{
		DryRun:    false,
		Namespace: cfg.NamespaceAI,
		Values:    []string{},
		Version:   "v0.0.0-latest",
	}

	return helmCmd.Install("aieg-crd", "envoyproxy/ai-gateway-crds-helm", cfg.NamespaceAI, opts)
}

func installAIGatewayController(helmCmd *helm.HelmCommand, cfg *config.Config) error {
	if err := helmCmd.RepoAdd("envoyproxy-ai", "oci://docker.io/envoyproxy"); err != nil {
		return err
	}

	if err := helmCmd.RepoUpdate(); err != nil {
		return err
	}

	values := []string{}
	if valuesExtra != "" {
		extraValues := strings.Split(valuesExtra, ",")
		for _, v := range extraValues {
			v = strings.TrimSpace(v)
			if v != "" {
				values = append(values, v)
			}
		}
	}

	opts := &helm.HelmOptions{
		DryRun:    false,
		Namespace: cfg.NamespaceAI,
		Values:    values,
		Version:   "v0.0.0-latest",
	}

	return helmCmd.Install("aieg", "envoyproxy/ai-gateway-helm", cfg.NamespaceAI, opts)
}

func installRedis(helmCmd *helm.HelmCommand, cfg *config.Config) error {
	if err := helmCmd.RepoAdd("bitnami", "https://charts.bitnami.com/bitnami"); err != nil {
		return err
	}

	if err := helmCmd.RepoUpdate(); err != nil {
		return err
	}

	opts := &helm.HelmOptions{
		DryRun:    false,
		Namespace: cfg.NamespaceAI,
		Values:    []string{},
	}

	return helmCmd.Install("envoy-redis", "bitnami/redis", cfg.NamespaceAI, opts)
}

func fetchRemoteValuesFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch remote file: HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "envoy-ai-values-*.yaml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
