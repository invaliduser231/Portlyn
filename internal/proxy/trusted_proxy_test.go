package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestURLIgnoresSpoofedForwardedProto(t *testing.T) {
	manager := NewManager(nil, nil, nil, nil, nil, nil, nil, ManagerOptions{})
	req := httptest.NewRequest(http.MethodGet, "http://app.example.com/private", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	req.Header.Set("X-Forwarded-Proto", "https")

	if got := manager.requestURL(req); got != "http://app.example.com/private" {
		t.Fatalf("expected spoofed proto to be ignored, got %q", got)
	}
}

func TestRequestURLTrustsForwardedProtoFromConfiguredProxy(t *testing.T) {
	manager := NewManager(nil, nil, nil, nil, nil, nil, nil, ManagerOptions{TrustedProxyCIDRs: []string{"10.0.0.0/8"}})
	req := httptest.NewRequest(http.MethodGet, "http://app.example.com/private", nil)
	req.RemoteAddr = "10.1.2.3:12345"
	req.Header.Set("X-Forwarded-Proto", "https")

	if got := manager.requestURL(req); got != "https://app.example.com/private" {
		t.Fatalf("expected trusted forwarded proto, got %q", got)
	}
}
