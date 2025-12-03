package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/franck-sorel/envoy-ai-unified-installer/pkg/helm"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health and prerequisites",
	Long: `Perform a health check on your system to ensure all prerequisites
for installing Envoy AI Gateway are met.

This command verifies:
- kubectl connectivity and cluster access
- helm installation and functionality
- kubernetes namespaces
- optional components (Redis, etc.)`,
	RunE: runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("🏥 System Health Check")
	fmt.Println()

	var allHealthy = true

	if !checkKubectl() {
		allHealthy = false
	}

	if !checkHelm() {
		allHealthy = false
	}

	if !checkKubernetesConnection() {
		allHealthy = false
	}

	namespaceGW := viper.GetString("namespace_gateway")
	namespaceAI := viper.GetString("namespace_ai")

	if !checkNamespace(namespaceGW) {
		allHealthy = false
	}

	if !checkNamespace(namespaceAI) {
		allHealthy = false
	}

	if !checkRedis(namespaceAI) {
		fmt.Println("⚠️  Redis:              Not installed (optional - install with --with-redis if needed)")
	}

	fmt.Println()
	if allHealthy {
		fmt.Println("✅ All checks passed! You're ready to install Envoy AI Gateway.")
	} else {
		fmt.Println("❌ Some checks failed. Please address the issues above.")
		return fmt.Errorf("system health check failed")
	}

	return nil
}

func checkKubectl() bool {
	fmt.Print("🔍 kubectl:            ")
	if _, err := exec.LookPath("kubectl"); err != nil {
		fmt.Println("❌ NOT FOUND")
		fmt.Println("   Install kubectl: https://kubernetes.io/docs/tasks/tools/")
		return false
	}

	cmd := exec.Command("kubectl", "version", "--client", "--short")
	if output, err := cmd.Output(); err != nil {
		fmt.Println("❌ FAILED")
		return false
	} else {
		fmt.Printf("✅ %s", string(output))
	}
	return true
}

func checkHelm() bool {
	fmt.Print("🔍 Helm:               ")
	if err := helm.ValidateHelmInstalled(); err != nil {
		fmt.Println("❌ NOT FOUND")
		fmt.Println("   Install Helm: https://helm.sh/docs/intro/install/")
		return false
	}

	helmCmd := helm.NewHelmCommand(false)
	version, err := helmCmd.Version()
	if err != nil {
		fmt.Println("❌ FAILED")
		return false
	}

	fmt.Printf("✅ %s", version)
	return true
}

func checkKubernetesConnection() bool {
	fmt.Print("🔍 Kubernetes cluster: ")
	cmd := exec.Command("kubectl", "cluster-info")
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ NOT CONNECTED")
		fmt.Println("   Configure your kubeconfig or check cluster connectivity")
		return false
	}
	fmt.Println("✅ CONNECTED")
	return true
}

func checkNamespace(namespace string) bool {
	fmt.Printf("🔍 Namespace '%s':    ", namespace)
	cmd := exec.Command("kubectl", "get", "namespace", namespace)
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ NOT FOUND")
		fmt.Printf("   Will be created during installation\n")
		return true
	}
	fmt.Println("✅ EXISTS")
	return true
}

func checkRedis(namespace string) bool {
	fmt.Print("🔍 Redis:              ")

	cmd := exec.Command("kubectl", "get", "pod", "-n", namespace,
		"-l", "app=redis", "-o", "jsonpath={.items[0].metadata.name}")

	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return false
	}

	fmt.Printf("✅ Pod: %s\n", string(output))
	return true
}
