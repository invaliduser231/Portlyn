package selfupdate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var downloadClient = &http.Client{Timeout: 5 * time.Minute}

func DownloadAsset(ctx context.Context, baseURL, name string, dst io.Writer) (int64, error) {
	url := strings.TrimRight(baseURL, "/") + "/" + name
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	resp, err := downloadClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("download %s: status %d: %s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	n, err := io.Copy(dst, resp.Body)
	if err != nil {
		return n, fmt.Errorf("copy %s: %w", url, err)
	}
	return n, nil
}

func DownloadString(ctx context.Context, baseURL, name string) (string, error) {
	var sb strings.Builder
	if _, err := DownloadAsset(ctx, baseURL, name, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}
