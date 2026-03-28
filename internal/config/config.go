package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion            string
	HTTPAddr              string
	FrontendBaseURL       string
	ProxyHTTPAddr         string
	ProxyHTTPSAddr        string
	RedisURL              string
	DatabaseDriver        string
	DatabaseURL           string
	DatabasePath          string
	JWTSecret             string
	JWTIssuer             string
	TokenTTL              time.Duration
	RefreshTokenTTL       time.Duration
	LogLevel              slog.Level
	AllowedOrigins        []string
	AdminEmail            string
	AdminPassword         string
	ViewerEmail           string
	ViewerPassword        string
	TLSCertFile           string
	TLSKeyFile            string
	RedirectHTTPToHTTPS   bool
	CertificateStorageDir string
	ACMEEmail             string
	ACMECAURL             string
	ACMEEnabled           bool
	ACMELeader            bool
	ACMEPollInterval      time.Duration
	ACMERenewWithin       time.Duration
	NodeOfflineAfter      time.Duration
	OIDC                  OIDCConfig
	OTP                   OTPConfig
	RouteAuthTTL          time.Duration
	AuthRateLimit         RateLimitConfig
	AuthCacheTTL          time.Duration
	AuditBufferSize       int
	AuditBatchSize        int
	AuditFlushInterval    time.Duration
	AuditDropPolicy       string
	RouteCacheTTL         time.Duration
	RouteLocalCacheTTL    time.Duration
	RouteLocalCacheSize   int
	AllowInsecureDevMode  bool
	RequireMFAForAdmins   bool
	CSRFTokenTTL          time.Duration
	RequestBodyLimitBytes int64
}

type OIDCConfig struct {
	Enabled              bool
	IssuerURL            string
	ClientID             string
	ClientSecret         string
	RedirectURL          string
	AllowedEmailDomains  []string
	AdminRoleClaimPath   string
	AdminRoleValue       string
	DefaultProviderLabel string
	AllowEmailLinking    bool
	RequireVerifiedEmail bool
}

type OTPConfig struct {
	Enabled              bool
	TokenTTL             time.Duration
	RequestLimit         int
	RequestWindow        time.Duration
	ResponseIncludesCode bool
}

type RateLimitConfig struct {
	LoginAttempts int
	Window        time.Duration
}

type ValidationIssue struct {
	Level   string `json:"level"`
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppVersion:            getEnv("APP_VERSION", "dev"),
		HTTPAddr:              getEnv("HTTP_ADDR", ":8080"),
		FrontendBaseURL:       strings.TrimRight(getEnv("FRONTEND_BASE_URL", "http://localhost:3000"), "/"),
		ProxyHTTPAddr:         getEnv("PROXY_HTTP_ADDR", ":80"),
		ProxyHTTPSAddr:        getEnv("PROXY_HTTPS_ADDR", ":443"),
		RedisURL:              strings.TrimSpace(os.Getenv("REDIS_URL")),
		DatabaseDriver:        strings.ToLower(getEnv("DATABASE_DRIVER", "")),
		DatabaseURL:           strings.TrimSpace(os.Getenv("DATABASE_URL")),
		DatabasePath:          getEnv("DATABASE_PATH", "portlyn.db"),
		JWTSecret:             getEnv("JWT_SECRET", "change-me-in-production"),
		JWTIssuer:             getEnv("JWT_ISSUER", "portlyn"),
		TokenTTL:              getEnvDuration("TOKEN_TTL", 24*time.Hour),
		RefreshTokenTTL:       getEnvDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
		LogLevel:              getEnvLogLevel("LOG_LEVEL", slog.LevelInfo),
		AllowedOrigins:        getEnvList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://127.0.0.1:3000"}),
		AdminEmail:            strings.TrimSpace(os.Getenv("ADMIN_EMAIL")),
		AdminPassword:         os.Getenv("ADMIN_PASSWORD"),
		ViewerEmail:           strings.TrimSpace(os.Getenv("VIEWER_EMAIL")),
		ViewerPassword:        os.Getenv("VIEWER_PASSWORD"),
		TLSCertFile:           strings.TrimSpace(os.Getenv("TLS_CERT_FILE")),
		TLSKeyFile:            strings.TrimSpace(os.Getenv("TLS_KEY_FILE")),
		RedirectHTTPToHTTPS:   getEnvBool("REDIRECT_HTTP_TO_HTTPS", false),
		CertificateStorageDir: getEnv("CERTIFICATE_STORAGE_DIR", "certificates"),
		ACMEEmail:             strings.TrimSpace(os.Getenv("ACME_EMAIL")),
		ACMECAURL:             strings.TrimSpace(os.Getenv("ACME_CA_URL")),
		ACMEEnabled:           getEnvBool("ACME_ENABLED", false),
		ACMELeader:            getEnvBool("ACME_LEADER", false),
		ACMEPollInterval:      getEnvDuration("ACME_POLL_INTERVAL", 1*time.Minute),
		ACMERenewWithin:       getEnvDuration("ACME_RENEW_WITHIN", 30*24*time.Hour),
		NodeOfflineAfter:      getEnvDuration("NODE_OFFLINE_AFTER", 2*time.Minute),
		OIDC: OIDCConfig{
			Enabled:              getEnvBool("OIDC_ENABLED", false),
			IssuerURL:            strings.TrimSpace(os.Getenv("OIDC_ISSUER_URL")),
			ClientID:             strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID")),
			ClientSecret:         os.Getenv("OIDC_CLIENT_SECRET"),
			RedirectURL:          strings.TrimSpace(os.Getenv("OIDC_REDIRECT_URL")),
			AllowedEmailDomains:  getEnvList("OIDC_ALLOWED_EMAIL_DOMAINS", []string{}),
			AdminRoleClaimPath:   getEnv("OIDC_ADMIN_ROLE_CLAIM_PATH", "realm_access.roles"),
			AdminRoleValue:       strings.TrimSpace(os.Getenv("OIDC_ADMIN_ROLE_VALUE")),
			DefaultProviderLabel: getEnv("OIDC_PROVIDER_LABEL", "SSO"),
			AllowEmailLinking:    getEnvBool("OIDC_ALLOW_EMAIL_LINKING", false),
			RequireVerifiedEmail: getEnvBool("OIDC_REQUIRE_VERIFIED_EMAIL", true),
		},
		OTP: OTPConfig{
			Enabled:              getEnvBool("OTP_ENABLED", true),
			TokenTTL:             getEnvDuration("OTP_TOKEN_TTL", 10*time.Minute),
			RequestLimit:         getEnvInt("OTP_REQUEST_LIMIT", 5),
			RequestWindow:        getEnvDuration("OTP_REQUEST_WINDOW", 15*time.Minute),
			ResponseIncludesCode: getEnvBool("OTP_RESPONSE_INCLUDES_CODE", false),
		},
		RouteAuthTTL: getEnvDuration("ROUTE_AUTH_TTL", 12*time.Hour),
		AuthRateLimit: RateLimitConfig{
			LoginAttempts: getEnvInt("AUTH_RATE_LIMIT_ATTEMPTS", 10),
			Window:        getEnvDuration("AUTH_RATE_LIMIT_WINDOW", 10*time.Minute),
		},
		AuthCacheTTL:          getEnvDuration("AUTH_CACHE_TTL", 1*time.Minute),
		AuditBufferSize:       getEnvInt("AUDIT_BUFFER_SIZE", 1024),
		AuditBatchSize:        getEnvInt("AUDIT_BATCH_SIZE", 128),
		AuditFlushInterval:    getEnvDuration("AUDIT_FLUSH_INTERVAL", 250*time.Millisecond),
		AuditDropPolicy:       strings.ToLower(getEnv("AUDIT_DROP_POLICY", "sync_fallback")),
		RouteCacheTTL:         getEnvDuration("ROUTE_CACHE_TTL", 30*time.Second),
		RouteLocalCacheTTL:    getEnvDuration("ROUTE_LOCAL_CACHE_TTL", 5*time.Second),
		RouteLocalCacheSize:   getEnvInt("ROUTE_LOCAL_CACHE_SIZE", 2048),
		AllowInsecureDevMode:  getEnvBool("ALLOW_INSECURE_DEV_MODE", false),
		RequireMFAForAdmins:   getEnvBool("REQUIRE_MFA_FOR_ADMINS", false),
		CSRFTokenTTL:          getEnvDuration("CSRF_TOKEN_TTL", 12*time.Hour),
		RequestBodyLimitBytes: getEnvInt64("REQUEST_BODY_LIMIT_BYTES", 1<<20),
	}

	return cfg, cfg.Validate()
}

func (cfg *Config) Validate() error {
	issues := cfg.ValidationIssues()
	for _, issue := range issues {
		if issue.Level == "error" {
			return fmt.Errorf("%s", issue.Message)
		}
	}
	return nil
}

func (cfg *Config) ValidationIssues() []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	add := func(level, field, code, message string) {
		issues = append(issues, ValidationIssue{Level: level, Field: field, Code: code, Message: message})
	}
	if cfg.JWTSecret == "" {
		add("error", "JWT_SECRET", "missing_secret", "JWT_SECRET must not be empty")
	}
	if cfg.DatabaseDriver == "" {
		if cfg.DatabaseURL != "" {
			cfg.DatabaseDriver = "postgres"
		} else {
			cfg.DatabaseDriver = "sqlite"
		}
	}
	switch cfg.DatabaseDriver {
	case "sqlite":
		if cfg.DatabasePath == "" {
			add("error", "DATABASE_PATH", "missing_database_path", "DATABASE_PATH must not be empty when DATABASE_DRIVER=sqlite")
		}
	case "postgres", "postgresql":
		cfg.DatabaseDriver = "postgres"
		if cfg.DatabaseURL == "" {
			add("error", "DATABASE_URL", "missing_database_url", "DATABASE_URL must not be empty when DATABASE_DRIVER=postgres")
		}
	default:
		add("error", "DATABASE_DRIVER", "unsupported_database_driver", fmt.Sprintf("unsupported DATABASE_DRIVER %q", cfg.DatabaseDriver))
	}
	if !cfg.AllowInsecureDevMode {
		if cfg.JWTSecret == "change-me-in-production" || len(cfg.JWTSecret) < 32 {
			add("error", "JWT_SECRET", "weak_secret", "JWT_SECRET must be unique and at least 32 characters unless ALLOW_INSECURE_DEV_MODE=true")
		}
		if cfg.OTP.ResponseIncludesCode {
			add("error", "OTP_RESPONSE_INCLUDES_CODE", "unsafe_otp_setting", "OTP_RESPONSE_INCLUDES_CODE must be false unless ALLOW_INSECURE_DEV_MODE=true")
		}
		if cfg.TokenTTL > 7*24*time.Hour {
			add("warn", "TOKEN_TTL", "long_token_ttl", "TOKEN_TTL is unusually long for an access token")
		}
	}
	if cfg.OIDC.Enabled {
		if cfg.OIDC.IssuerURL == "" || cfg.OIDC.ClientID == "" || cfg.OIDC.ClientSecret == "" || cfg.OIDC.RedirectURL == "" {
			add("error", "OIDC", "oidc_incomplete", "OIDC is enabled but OIDC_ISSUER_URL, OIDC_CLIENT_ID, OIDC_CLIENT_SECRET, and OIDC_REDIRECT_URL are incomplete")
		}
		if !strings.HasPrefix(strings.ToLower(cfg.OIDC.RedirectURL), "https://") && !cfg.AllowInsecureDevMode {
			add("error", "OIDC_REDIRECT_URL", "insecure_oidc_redirect", "OIDC_REDIRECT_URL must use https outside insecure dev mode")
		}
	}
	if cfg.OTP.Enabled && cfg.OTP.TokenTTL < 2*time.Minute {
		add("warn", "OTP_TOKEN_TTL", "short_otp_ttl", "OTP token TTL below 2 minutes may be too aggressive")
	}
	if cfg.AuthRateLimit.LoginAttempts < 3 {
		add("warn", "AUTH_RATE_LIMIT_ATTEMPTS", "tight_rate_limit", "AUTH_RATE_LIMIT_ATTEMPTS is very low and may create operator lockouts")
	}
	if cfg.RequestBodyLimitBytes <= 0 {
		add("error", "REQUEST_BODY_LIMIT_BYTES", "invalid_request_body_limit", "REQUEST_BODY_LIMIT_BYTES must be greater than zero")
	}
	if cfg.RedirectHTTPToHTTPS && cfg.ProxyHTTPSAddr == "" {
		add("error", "PROXY_HTTPS_ADDR", "missing_https_listener", "PROXY_HTTPS_ADDR must be configured when REDIRECT_HTTP_TO_HTTPS=true")
	}
	return issues
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvLogLevel(key string, fallback slog.Level) slog.Level {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		if numeric, err := strconv.Atoi(value); err == nil {
			return slog.Level(numeric)
		}
		return fallback
	}
}

func getEnvList(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
