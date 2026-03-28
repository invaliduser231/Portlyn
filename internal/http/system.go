package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"portlyn/internal/domain"
	"portlyn/internal/proxy"
)

type SystemOverview struct {
	Status       string                `json:"status"`
	CheckedAt    time.Time             `json:"checked_at"`
	Runtime      SystemRuntimeOverview `json:"runtime"`
	Certificates CertificateOverview   `json:"certificates"`
	Security     SecurityOverview      `json:"security"`
	Counts       SystemCounts          `json:"counts"`
	Warnings     SystemWarnings        `json:"warnings"`
	Health       SystemHealthBreakdown `json:"health"`
}

type SystemRuntimeOverview struct {
	Version             string    `json:"version"`
	APIStatus           string    `json:"api_status"`
	DBStatus            string    `json:"db_status"`
	ProxyStatus         string    `json:"proxy_status"`
	HTTPAddr            string    `json:"http_addr"`
	ProxyHTTPAddr       string    `json:"proxy_http_addr"`
	ProxyHTTPSAddr      string    `json:"proxy_https_addr"`
	TLSEnabled          bool      `json:"tls_enabled"`
	ACMEEnabled         bool      `json:"acme_enabled"`
	ACMEChallengeTypes  []string  `json:"acme_challenge_types"`
	RedirectHTTPToHTTPS bool      `json:"redirect_http_to_https"`
	CheckedAt           time.Time `json:"checked_at"`
}

type CertificateOverview struct {
	DefaultIssuer         string   `json:"default_issuer"`
	SupportedIssuers      []string `json:"supported_issuers"`
	DNSProviderCount      int64    `json:"dns_provider_count"`
	DNSProviderTypes      []string `json:"dns_provider_types"`
	ExpiringWindowDays    int      `json:"expiring_window_days"`
	SupportsWildcard      bool     `json:"supports_wildcard"`
	SupportsMultiSAN      bool     `json:"supports_multi_san"`
	SupportsDNSChallenges bool     `json:"supports_dns_challenges"`
}

type SecurityOverview struct {
	OIDCEnabled                bool   `json:"oidc_enabled"`
	OTPEnabled                 bool   `json:"otp_enabled"`
	CSRFEnabled                bool   `json:"csrf_enabled"`
	CookieSecure               bool   `json:"cookie_secure"`
	CookieHTTPOnly             bool   `json:"cookie_http_only"`
	CookieSameSiteSession      string `json:"cookie_same_site_session"`
	CookieSameSiteRefresh      string `json:"cookie_same_site_refresh"`
	RequireMFAForAdmins        bool   `json:"require_mfa_for_admins"`
	AdminsWithoutMFA           int    `json:"admins_without_mfa"`
	NodeHeartbeatAuthMode      string `json:"node_heartbeat_auth_mode"`
	NodeMTLSEnabled            bool   `json:"node_mtls_enabled"`
	SecurityHeadersEnabled     bool   `json:"security_headers_enabled"`
	SMTPEnabled                bool   `json:"smtp_enabled"`
	SMTPConfigured             bool   `json:"smtp_configured"`
	JWTTTLSeconds              int    `json:"jwt_ttl_seconds"`
	RefreshTokenTTLSeconds     int    `json:"refresh_token_ttl_seconds"`
	AuthRateLimitAttempts      int    `json:"auth_rate_limit_attempts"`
	AuthRateLimitWindowSeconds int    `json:"auth_rate_limit_window_seconds"`
}

type SystemCounts struct {
	Services        int64 `json:"services"`
	Domains         int64 `json:"domains"`
	Certificates    int64 `json:"certificates"`
	DNSProviders    int64 `json:"dns_providers"`
	NodesOnline     int   `json:"nodes_online"`
	NodesOffline    int   `json:"nodes_offline"`
	Users           int64 `json:"users"`
	Groups          int64 `json:"groups"`
	ServiceGroups   int64 `json:"service_groups"`
	ProxyRoutes     int   `json:"proxy_routes"`
	AuthFailures24h int64 `json:"auth_failures_24h"`
}

type SystemWarnings struct {
	ExpiringCertificates []domain.Certificate `json:"expiring_certificates"`
	FailedCertificates   []domain.Certificate `json:"failed_certificates"`
	OfflineNodes         []domain.Node        `json:"offline_nodes"`
	RiskyServices        []map[string]any     `json:"risky_services"`
	Config               []map[string]any     `json:"config"`
}

type SystemHealthBreakdown struct {
	Livez    StatusCondition   `json:"livez"`
	Readyz   []StatusCondition `json:"readyz"`
	Services []StatusCondition `json:"services"`
	Cluster  []StatusCondition `json:"cluster"`
}

func (s *Server) handleSystemOverview(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()

	usersCount, err := s.users.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	groupCount, err := s.groups.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	serviceGroupCount, err := s.serviceGroups.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	serviceCount, err := s.services.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	domainCount, err := s.domains.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	certificateCount, err := s.certificates.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	dnsProviderCount, err := s.dnsProviders.Count(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	nodes, err := s.nodes.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	services, err := s.services.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	certificates, err := s.certificates.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	dnsProviders, err := s.dnsProviders.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	failedAuthCount, err := s.auditStore.CountByActionLikeSince(r.Context(), "%failed%", now.Add(-24*time.Hour))
	if err != nil {
		s.internalError(w, err)
		return
	}
	settings, _ := s.appSettings.Get(r.Context())
	if settings == nil {
		settings = &domain.AppSettings{}
	}
	users, err := s.users.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	adminsWithoutMFA := countAdminsWithoutMFA(users)
	onlineNodes, offlineNodes, offlineNodeItems := s.nodeStatusSummary(nodes)
	expiring, failed := summarizeCertificates(certificates, now)
	riskyServices := summarizeRiskyServices(services, settings)

	healthCtx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	ready := s.evaluateReadiness(healthCtx)
	serviceHealth := s.evaluateServiceStatus(healthCtx)
	clusterHealth := s.evaluateClusterStatus(healthCtx)
	summary := summarizeConditions(append(append([]StatusCondition{}, ready...), append(serviceHealth, clusterHealth...)...))
	status := overallHealthLevel(summary)
	dbStatus := HealthLevelOK
	if summary[HealthLevelError] > 0 {
		dbStatus = overallDependencyLevel(ready, "database")
	}

	writeJSON(w, http.StatusOK, SystemOverview{
		Status:    status,
		CheckedAt: now,
		Runtime: SystemRuntimeOverview{
			Version:             s.cfg.AppVersion,
			APIStatus:           overallHealthLevel(summarizeConditions(ready)),
			DBStatus:            dbStatus,
			ProxyStatus:         overallHealthLevel(summarizeConditions(serviceHealth)),
			HTTPAddr:            s.cfg.HTTPAddr,
			ProxyHTTPAddr:       s.cfg.ProxyHTTPAddr,
			ProxyHTTPSAddr:      s.cfg.ProxyHTTPSAddr,
			TLSEnabled:          s.acme != nil && s.acme.HasHTTPS(),
			ACMEEnabled:         s.cfg.ACMEEnabled,
			ACMEChallengeTypes:  []string{domain.CertificateChallengeHTTP01, domain.CertificateChallengeDNS01},
			RedirectHTTPToHTTPS: s.cfg.RedirectHTTPToHTTPS,
			CheckedAt:           now,
		},
		Certificates: CertificateOverview{
			DefaultIssuer:         domain.CertificateIssuerLetsEncryptProd,
			SupportedIssuers:      []string{domain.CertificateIssuerLetsEncryptProd, domain.CertificateIssuerLetsEncryptStaging},
			DNSProviderCount:      dnsProviderCount,
			DNSProviderTypes:      summarizeDNSProviderTypes(dnsProviders),
			ExpiringWindowDays:    14,
			SupportsWildcard:      true,
			SupportsMultiSAN:      true,
			SupportsDNSChallenges: true,
		},
		Security: SecurityOverview{
			OIDCEnabled:                settings.OIDCEnabled,
			OTPEnabled:                 settings.OTPEnabled,
			CSRFEnabled:                true,
			CookieSecure:               !s.cfg.AllowInsecureDevMode,
			CookieHTTPOnly:             true,
			CookieSameSiteSession:      "Lax",
			CookieSameSiteRefresh:      "Strict",
			RequireMFAForAdmins:        settings.RequireMFAForAdmins,
			AdminsWithoutMFA:           adminsWithoutMFA,
			NodeHeartbeatAuthMode:      "token",
			NodeMTLSEnabled:            false,
			SecurityHeadersEnabled:     true,
			SMTPEnabled:                settings.SMTPEnabled,
			SMTPConfigured:             strings.TrimSpace(settings.SMTPHost) != "" && settings.SMTPPort > 0 && strings.TrimSpace(settings.SMTPFromEmail) != "",
			JWTTTLSeconds:              int(s.cfg.TokenTTL.Seconds()),
			RefreshTokenTTLSeconds:     int(s.cfg.RefreshTokenTTL.Seconds()),
			AuthRateLimitAttempts:      s.cfg.AuthRateLimit.LoginAttempts,
			AuthRateLimitWindowSeconds: int(s.cfg.AuthRateLimit.Window.Seconds()),
		},
		Counts: SystemCounts{
			Services:        serviceCount,
			Domains:         domainCount,
			Certificates:    certificateCount,
			DNSProviders:    dnsProviderCount,
			NodesOnline:     onlineNodes,
			NodesOffline:    offlineNodes,
			Users:           usersCount,
			Groups:          groupCount,
			ServiceGroups:   serviceGroupCount,
			ProxyRoutes:     len(s.proxy.RuntimeRoutes()),
			AuthFailures24h: failedAuthCount,
		},
		Warnings: SystemWarnings{
			ExpiringCertificates: expiring,
			FailedCertificates:   failed,
			OfflineNodes:         offlineNodeItems,
			RiskyServices:        riskyServices,
			Config:               summarizeConfigWarnings(settings),
		},
		Health: SystemHealthBreakdown{
			Livez:    StatusCondition{Name: "process", Scope: "livez", Level: HealthLevelOK, Summary: "process is serving requests", Required: true, CheckedAt: now},
			Readyz:   ready,
			Services: serviceHealth,
			Cluster:  clusterHealth,
		},
	})
}

func overallDependencyLevel(items []StatusCondition, name string) string {
	for _, item := range items {
		if item.Name == name {
			return item.Level
		}
	}
	return HealthLevelOK
}

func (s *Server) nodeStatusSummary(nodes []domain.Node) (int, int, []domain.Node) {
	online := 0
	offline := 0
	offlineNodes := make([]domain.Node, 0)
	for i := range nodes {
		s.evaluateNodeStatus(&nodes[i])
		if nodes[i].Status == domain.NodeStatusOnline {
			online++
			continue
		}
		offline++
		offlineNodes = append(offlineNodes, nodes[i])
	}
	return online, offline, offlineNodes
}

func summarizeCertificates(items []domain.Certificate, now time.Time) ([]domain.Certificate, []domain.Certificate) {
	expiring := make([]domain.Certificate, 0)
	failed := make([]domain.Certificate, 0)
	limit := now.Add(14 * 24 * time.Hour)
	for _, item := range items {
		if item.Status == domain.CertificateStatusFailed {
			failed = append(failed, item)
		}
		if !item.ExpiresAt.IsZero() && item.ExpiresAt.Before(limit) {
			item.Status = domain.CertificateStatusExpiringSoon
			expiring = append(expiring, item)
		}
	}
	return expiring, failed
}

func summarizeDNSProviderTypes(items []domain.DNSProvider) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item.Type]; ok {
			continue
		}
		seen[item.Type] = struct{}{}
		out = append(out, item.Type)
	}
	return out
}

func summarizeRiskyServices(items []domain.Service, settings *domain.AppSettings) []map[string]any {
	risky := make([]map[string]any, 0)
	for _, item := range items {
		policy, method, _, inherited := proxy.EffectiveAccessForService(item)
		score, reasons := serviceRiskAssessment(item, policy.AccessMode, method, item.AccessMethodConfig, settings)
		if score == "low" || len(reasons) == 0 {
			continue
		}
		entry := map[string]any{
			"id":               item.ID,
			"name":             item.Name,
			"domain_name":      domain.ServiceHost(item),
			"path":             item.Path,
			"access_mode":      policy.AccessMode,
			"access_method":    method,
			"risk_score":       score,
			"reasons":          reasons,
			"use_group_policy": item.UseGroupPolicy,
		}
		if inherited != nil {
			entry["inherited_from_group"] = inherited.Name
		}
		risky = append(risky, entry)
	}
	return risky
}

func summarizeConfigWarnings(settings *domain.AppSettings) []map[string]any {
	items := make([]map[string]any, 0)
	if settings.OIDCEnabled && (strings.TrimSpace(settings.OIDCIssuerURL) == "" || strings.TrimSpace(settings.OIDCClientID) == "" || strings.TrimSpace(settings.OIDCClientSecret) == "") {
		items = append(items, map[string]any{"code": "oidc_incomplete", "message": "OIDC is enabled but not fully configured."})
	}
	if settings.OTPEnabled && (!settings.SMTPEnabled || strings.TrimSpace(settings.SMTPHost) == "") {
		items = append(items, map[string]any{"code": "otp_without_smtp", "message": "OTP is enabled but SMTP is not fully configured."})
	}
	if settings.RequireMFAForAdmins {
		items = append(items, map[string]any{"code": "admin_mfa_required", "message": "Admin MFA enforcement is enabled."})
	}
	if settings.SMTPEnabled && (strings.TrimSpace(settings.SMTPHost) == "" || settings.SMTPPort <= 0 || strings.TrimSpace(settings.SMTPFromEmail) == "") {
		items = append(items, map[string]any{"code": "smtp_incomplete", "message": "SMTP is enabled but incomplete."})
	}
	return items
}

func countAdminsWithoutMFA(users []domain.User) int {
	total := 0
	for _, user := range users {
		if user.Role == domain.RoleAdmin && user.Active && !user.MFAEnabled {
			total++
		}
	}
	return total
}
