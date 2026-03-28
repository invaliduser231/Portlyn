package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portlyn/internal/config"
	"portlyn/internal/domain"
)

func TestHandleReadyzReturnsOKWhenDependenciesHealthy(t *testing.T) {
	server := &Server{
		cfg:    config.Config{AppVersion: "test", ProxyHTTPAddr: ":80", ProxyHTTPSAddr: ":443"},
		logger: slog.Default(),
		healthChecks: []HealthCheck{
			NewNamedHealthCheck("database", func(_ context.Context) error { return nil }),
			NewNamedHealthCheck("redis", func(_ context.Context) error { return nil }),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	server.handleReadyz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %#v", body["status"])
	}
}

func TestHandleReadyzReturnsServiceUnavailableWhenDependencyFails(t *testing.T) {
	server := &Server{
		cfg:    config.Config{AppVersion: "test"},
		logger: slog.Default(),
		healthChecks: []HealthCheck{
			NewNamedHealthCheck("database", func(_ context.Context) error { return errors.New("db down") }),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	server.handleReadyz(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var body struct {
		Status string            `json:"status"`
		Readyz []StatusCondition `json:"readyz"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Status != "error" {
		t.Fatalf("expected overall error, got %q", body.Status)
	}
	if len(body.Readyz) != 1 || body.Readyz[0].Level != HealthLevelError {
		t.Fatalf("expected failed dependency details, got %#v", body.Readyz)
	}
}

func TestHandleLivezAlwaysReturnsOK(t *testing.T) {
	server := &Server{
		cfg:    config.Config{AppVersion: "test"},
		logger: slog.Default(),
	}

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec := httptest.NewRecorder()

	server.handleLivez(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHTTPTargetsHealthCheckPassesWhenTargetResponds(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead && r.Method != http.MethodGet {
			t.Fatalf("unexpected method %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	check := NewHTTPTargetsHealthCheck("service_targets", func(context.Context) ([]HTTPHealthTarget, error) {
		return []HTTPHealthTarget{{Name: "app", URL: upstream.URL}}, nil
	}, upstream.Client())

	if err := check.Check(context.Background()); err != nil {
		t.Fatalf("expected target to be healthy, got %v", err)
	}
}

func TestHTTPTargetsHealthCheckFailsWhenTargetReturnsServerError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer upstream.Close()

	check := NewHTTPTargetsHealthCheck("service_targets", func(context.Context) ([]HTTPHealthTarget, error) {
		return []HTTPHealthTarget{{Name: "app", URL: upstream.URL}}, nil
	}, upstream.Client())

	err := check.Check(context.Background())
	if err == nil {
		t.Fatal("expected target check to fail")
	}
	if got := err.Error(); got != "app: returned status 502" {
		t.Fatalf("unexpected error: %s", got)
	}
}

func TestServiceTargetsFromDomainBuildsDisplayName(t *testing.T) {
	targets := ServiceTargetsFromDomain([]domain.Service{
		{
			Name:      "Genusswerk",
			Path:      "/",
			TargetURL: "http://host.docker.internal:3333",
			Domain:    domain.Domain{Name: "genusswerk.local"},
		},
	})

	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].Name != "Genusswerk (genusswerk.local/)" {
		t.Fatalf("unexpected target name %q", targets[0].Name)
	}
	if targets[0].URL != "http://host.docker.internal:3333" {
		t.Fatalf("unexpected target url %q", targets[0].URL)
	}
}

func TestEvaluateServiceHealthReturnsPendingWhenServiceNotDeployed(t *testing.T) {
	server := &Server{}

	health := server.evaluateServiceHealth(context.Background(), domain.Service{
		Name:      "Genusswerk",
		TargetURL: "http://host.docker.internal:333",
	})

	if health.Status != "pending" {
		t.Fatalf("expected pending, got %q", health.Status)
	}
	if health.Error != "" {
		t.Fatalf("expected no error, got %q", health.Error)
	}
}

func TestEvaluateServiceHealthReturnsUnhealthyWhenTargetIsDown(t *testing.T) {
	server := &Server{}
	now := time.Now()

	health := server.evaluateServiceHealth(context.Background(), domain.Service{
		Name:           "Genusswerk",
		TargetURL:      "http://127.0.0.1:333",
		LastDeployedAt: &now,
	})

	if health.Status != "unhealthy" {
		t.Fatalf("expected unhealthy, got %q", health.Status)
	}
	if health.Error == "" {
		t.Fatal("expected an error message for unhealthy target")
	}
}
