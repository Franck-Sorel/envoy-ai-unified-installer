package helm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type HelmOptions struct {
	DryRun     bool
	Namespace  string
	Values     []string
	Version    string
	ChartRepo  string
}

type HelmCommand struct {
	dryRun bool
	output io.Writer
}

func NewHelmCommand(dryRun bool) *HelmCommand {
	return &HelmCommand{
		dryRun: dryRun,
		output: os.Stdout,
	}
}

func (h *HelmCommand) Execute(args ...string) error {
	if h.dryRun {
		fmt.Printf("[DRY-RUN] helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = h.output
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("helm command failed: %w", err)
	}

	return nil
}

func (h *HelmCommand) ExecuteOutput(args ...string) (string, error) {
	if h.dryRun {
		fmt.Printf("[DRY-RUN] helm %s\n", strings.Join(args, " "))
		return "", nil
	}

	cmd := exec.Command("helm", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("helm command failed: %w", err)
	}

	return out.String(), nil
}

func (h *HelmCommand) RepoAdd(name, url string) error {
	return h.Execute("repo", "add", name, url, "--force-update")
}

func (h *HelmCommand) RepoUpdate() error {
	return h.Execute("repo", "update")
}

func (h *HelmCommand) Install(releaseName, chart, namespace string, opts *HelmOptions) error {
	args := []string{"upgrade", "--install", releaseName, chart}

	args = append(args, "-n", namespace, "--create-namespace")

	if opts.Version != "" {
		args = append(args, "--version", opts.Version)
	}

	for _, v := range opts.Values {
		args = append(args, "-f", v)
	}

	if opts.DryRun {
		args = append(args, "--dry-run", "--debug")
	}

	return h.Execute(args...)
}

func (h *HelmCommand) Uninstall(releaseName, namespace string) error {
	if h.dryRun {
		fmt.Printf("[DRY-RUN] helm uninstall %s -n %s\n", releaseName, namespace)
		return nil
	}

	cmd := exec.Command("helm", "uninstall", releaseName, "-n", namespace)
	cmd.Stdout = h.output
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (h *HelmCommand) GetValues(releaseName, namespace string) (string, error) {
	return h.ExecuteOutput("get", "values", releaseName, "-n", namespace)
}

func (h *HelmCommand) List(namespace string) (string, error) {
	return h.ExecuteOutput("list", "-n", namespace)
}

func (h *HelmCommand) Version() (string, error) {
	return h.ExecuteOutput("version", "--short")
}

func ValidateHelmInstalled() error {
	cmd := exec.Command("helm", "version", "--short")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("helm is not installed or not in PATH: %w", err)
	}
	return nil
}
