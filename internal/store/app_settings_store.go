package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"portlyn/internal/config"
	"portlyn/internal/domain"
)

const appSettingsSingletonID uint = 1

type AppSettingsStore struct {
	db *gorm.DB
}

func NewAppSettingsStore(db *gorm.DB) *AppSettingsStore {
	return &AppSettingsStore{db: db}
}

func (s *AppSettingsStore) Get(ctx context.Context) (*domain.AppSettings, error) {
	var item domain.AppSettings
	err := s.db.WithContext(ctx).First(&item, appSettingsSingletonID).Error
	if err == nil {
		return &item, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return nil, err
}

func (s *AppSettingsStore) Upsert(ctx context.Context, item *domain.AppSettings) error {
	item.ID = appSettingsSingletonID
	return s.db.WithContext(ctx).Save(item).Error
}

func (s *AppSettingsStore) SeedDefaults(ctx context.Context, cfg config.Config) error {
	_, err := s.Get(ctx)
	if err == nil {
		return nil
	}
	if err != ErrNotFound {
		return err
	}

	item := &domain.AppSettings{
		ID:                        appSettingsSingletonID,
		FrontendBaseURL:           cfg.FrontendBaseURL,
		AuthBrandName:             "Portlyn",
		AuthBackgroundColor:       "#0a0d14",
		AuthBackgroundAccent:      "#162033",
		AuthPanelColor:            "#111826",
		AuthButtonColor:           "#2f6fed",
		AuthTextColor:             "#f8fafc",
		AuthMutedTextColor:        "#94a3b8",
		AuthLoginTitle:            "Login",
		AuthRouteLoginTitle:       "Login",
		AuthForbiddenTitle:        "Access denied",
		AuthLoginPasswordLabel:    "Login",
		AuthLoginOIDCLabel:        "Continue with SSO",
		AuthLoginOTPRequestLabel:  "Request code",
		AuthLoginOTPVerifyLabel:   "Verify code",
		AuthRouteContinueLabel:    "Continue",
		AuthRouteOIDCLabel:        "Continue with SSO",
		AuthRoutePINLabel:         "Unlock",
		AuthRouteEmailSendLabel:   "Send code",
		AuthRouteEmailVerifyLabel: "Verify code",
		AuthForbiddenRetryLabel:   "Try again",
		OIDCEnabled:               cfg.OIDC.Enabled,
		OIDCIssuerURL:             cfg.OIDC.IssuerURL,
		OIDCClientID:              cfg.OIDC.ClientID,
		OIDCClientSecret:          cfg.OIDC.ClientSecret,
		OIDCRedirectURL:           cfg.OIDC.RedirectURL,
		OIDCAllowedEmailDomains:   domain.JSONStringSlice(cfg.OIDC.AllowedEmailDomains),
		OIDCAdminRoleClaimPath:    cfg.OIDC.AdminRoleClaimPath,
		OIDCAdminRoleValue:        cfg.OIDC.AdminRoleValue,
		OIDCProviderLabel:         cfg.OIDC.DefaultProviderLabel,
		OIDCAllowEmailLinking:     cfg.OIDC.AllowEmailLinking,
		OIDCRequireVerifiedEmail:  cfg.OIDC.RequireVerifiedEmail,
		OTPEnabled:                cfg.OTP.Enabled,
		OTPTokenTTLSeconds:        int(cfg.OTP.TokenTTL.Seconds()),
		OTPRequestLimit:           cfg.OTP.RequestLimit,
		OTPRequestWindowSeconds:   int(cfg.OTP.RequestWindow.Seconds()),
		RequireMFAForAdmins:       cfg.RequireMFAForAdmins,
		SMTPPort:                  587,
		SMTPEncryption:            "starttls",
	}
	return s.Upsert(ctx, item)
}
