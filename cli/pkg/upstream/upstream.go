package upstream

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

type ChartRelease struct {
	Owner   string
	Repo    string
	Version string
	URL     string
}

func GetGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc)
	}
	return github.NewClient(nil)
}

func FetchLatestRelease(owner, repo string) (*ChartRelease, error) {
	client := GetGitHubClient()
	ctx := context.Background()

	rel, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release for %s/%s: %w", owner, repo, err)
	}

	url := findChartAsset(rel)
	if url == "" {
		return nil, fmt.Errorf("no chart asset found for %s/%s", owner, repo)
	}

	return &ChartRelease{
		Owner:   owner,
		Repo:    repo,
		Version: rel.GetTagName(),
		URL:     url,
	}, nil
}

func findChartAsset(rel *github.RepositoryRelease) string {
	keywords := []string{"helm", "chart", ".tgz", "tar.gz"}

	for _, asset := range rel.Assets {
		name := asset.GetName()
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(name), keyword) {
				return asset.GetBrowserDownloadURL()
			}
		}
	}

	return ""
}

func GetUpstreamCharts() ([]ChartRelease, error) {
	upstreams := []struct {
		owner string
		repo  string
	}{
		{"envoyproxy", "gateway"},
		{"envoyproxy", "ai-gateway-helm"},
		{"envoyproxy", "ai-gateway-crds-helm"},
		{"envoyproxy", "ai-gateway"},
	}

	var charts []ChartRelease
	var errors []string

	for _, up := range upstreams {
		chart, err := FetchLatestRelease(up.owner, up.repo)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		charts = append(charts, *chart)
	}

	if len(errors) > 0 {
		return charts, fmt.Errorf("errors fetching upstream charts:\n%s", strings.Join(errors, "\n"))
	}

	return charts, nil
}
