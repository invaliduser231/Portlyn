package config

import "testing"

func TestLoadRejectsMissingSecretsOutsideDevMode(t *testing.T) {
	t.Setenv("ALLOW_INSECURE_DEV_MODE", "false")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_SIGNING_SECRET", "")
	t.Setenv("SESSION_BRIDGE_SECRET", "")
	t.Setenv("OIDC_STATE_SECRET", "")
	t.Setenv("MFA_ENCRYPTION_SECRET", "")
	t.Setenv("CSRF_SECRET", "")
	t.Setenv("DATA_ENCRYPTION_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected Load to fail when required secrets are missing outside dev mode")
	}
}

func TestLoadGeneratesSecretsInDevMode(t *testing.T) {
	t.Setenv("ALLOW_INSECURE_DEV_MODE", "true")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_SIGNING_SECRET", "")
	t.Setenv("SESSION_BRIDGE_SECRET", "")
	t.Setenv("OIDC_STATE_SECRET", "")
	t.Setenv("MFA_ENCRYPTION_SECRET", "")
	t.Setenv("CSRF_SECRET", "")
	t.Setenv("DATA_ENCRYPTION_SECRET", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected Load to succeed in insecure dev mode, got error: %v", err)
	}
	if cfg.JWTSecret == "" || cfg.JWTSigningSecret == "" || cfg.DataEncryptionSecret == "" {
		t.Fatal("expected generated non-empty secrets in insecure dev mode")
	}
}
