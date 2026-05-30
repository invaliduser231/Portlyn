package selfupdate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withFakeGitHub(t *testing.T, handler http.HandlerFunc) func() {
	t.Helper()
	srv := httptest.NewServer(handler)
	original := apiClient
	apiClient = srv.Client()
	originalBase := githubAPIBase
	githubAPIBase = srv.URL
	return func() {
		apiClient = original
		githubAPIBase = originalBase
		srv.Close()
	}
}

func TestResolveLatestParsesTag(t *testing.T) {
	cleanup := withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/repos/owner/repo/releases/latest") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name":"v2.3.4","assets":[]}`))
	})
	defer cleanup()

	rel, err := ResolveLatest(context.Background(), "owner/repo")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if rel.Tag != "v2.3.4" {
		t.Fatalf("tag = %q", rel.Tag)
	}
	if !strings.HasSuffix(rel.AssetBaseURL, "/releases/download/v2.3.4") {
		t.Fatalf("base url = %q", rel.AssetBaseURL)
	}
}

func TestResolveByTagNotFound(t *testing.T) {
	cleanup := withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
	})
	defer cleanup()

	_, err := ResolveByTag(context.Background(), "owner/repo", "v9.9.9")
	if err == nil || !strings.Contains(err.Error(), "no release") {
		t.Fatalf("expected no-release error, got %v", err)
	}
}

func TestResolveByTagEmpty(t *testing.T) {
	_, err := ResolveByTag(context.Background(), "owner/repo", "  ")
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("expected empty-tag error, got %v", err)
	}
}
