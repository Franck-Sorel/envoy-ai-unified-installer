package main

import (
	"fmt"
	"os"

	"github.com/franck-sorel/envoy-ai-unified-installer/cmd"
)

var (
	version   = "0.1.0"
	gitCommit = "dev"
	buildTime = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, gitCommit, buildTime)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
