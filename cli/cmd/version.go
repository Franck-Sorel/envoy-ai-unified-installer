package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/helm"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/upstream"
)

var (
	cliVersion = "dev"
	gitCommit  = "unknown"
	buildTime  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and upstream component versions",
	Long: `Display the version of the envoy-ai-installer CLI and the versions
of upstream components (Envoy Gateway, AI Gateway, etc.) that will be installed.`,
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Println("📦 envoy-ai-installer Version Information")
	fmt.Println()
	fmt.Printf("  CLI Version:    %s\n", cliVersion)
	fmt.Printf("  Git Commit:     %s\n", gitCommit)
	fmt.Printf("  Build Time:     %s\n", buildTime)
	fmt.Println()

	helmCmd := helm.NewHelmCommand(false)
	helmVersion, err := helmCmd.Version()
	if err == nil {
		fmt.Printf("  Helm Version:   %s", helmVersion)
	}

	fmt.Println("\n📋 Upstream Component Versions")
	fmt.Println()

	charts, err := upstream.GetUpstreamCharts()
	if err != nil {
		fmt.Printf("  ⚠️  Could not fetch upstream versions: %v\n", err)
		return nil
	}

	for _, chart := range charts {
		fmt.Printf("  %s/%s:  %s\n", chart.Owner, chart.Repo, chart.Version)
	}

	return nil
}

func SetVersionInfo(version, commit, buildTime string) {
	cliVersion = version
	gitCommit = commit
	buildTime = buildTime
}
