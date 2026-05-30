package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Release struct {
	Tag          string
	AssetBaseURL string
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

const defaultUserAgent = "portlyn-selfupdate"

var (
	apiClient     = &http.Client{Timeout: 30 * time.Second}
	githubAPIBase = "https://api.github.com"
)

func ResolveLatest(ctx context.Context, repo string) (Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", githubAPIBase, repo)
	return fetchRelease(ctx, url, repo)
}

func ResolveByTag(ctx context.Context, repo, tag string) (Release, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return Release{}, fmt.Errorf("tag must not be empty")
	}
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", githubAPIBase, repo, tag)
	return fetchRelease(ctx, url, repo)
}

func fetchRelease(ctx context.Context, url, repo string) (Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Release{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", defaultUserAgent)
	resp, err := apiClient.Do(req)
	if err != nil {
		return Release{}, fmt.Errorf("github api request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return Release{}, fmt.Errorf("no release found at %s", url)
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return Release{}, fmt.Errorf("github api %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var payload githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Release{}, fmt.Errorf("decode release: %w", err)
	}
	if strings.TrimSpace(payload.TagName) == "" {
		return Release{}, fmt.Errorf("release at %s has no tag_name", url)
	}
	return Release{
		Tag:          payload.TagName,
		AssetBaseURL: fmt.Sprintf("https://github.com/%s/releases/download/%s", repo, payload.TagName),
	}, nil
}
