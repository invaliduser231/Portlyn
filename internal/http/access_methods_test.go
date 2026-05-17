package http

import (
	"testing"

	"portlyn/internal/domain"
)

func TestRestrictedPolicyAllowsUserByRoleAndGroup(t *testing.T) {
	user := &domain.User{ID: 1, Role: domain.RoleViewer, Active: true}

	if !restrictedPolicyAllowsUser(user, nil, domain.AccessPolicy{
		AccessMode:   domain.AccessModeRestricted,
		AllowedRoles: domain.JSONStringSlice{domain.RoleViewer},
	}) {
		t.Fatal("expected viewer role to satisfy restricted policy")
	}

	if !restrictedPolicyAllowsUser(user, []uint{42}, domain.AccessPolicy{
		AccessMode:    domain.AccessModeRestricted,
		AllowedGroups: domain.JSONUintSlice{42},
	}) {
		t.Fatal("expected matching group to satisfy restricted policy")
	}

	if restrictedPolicyAllowsUser(user, []uint{7}, domain.AccessPolicy{
		AccessMode:    domain.AccessModeRestricted,
		AllowedGroups: domain.JSONUintSlice{42},
	}) {
		t.Fatal("expected non-member group to be denied")
	}
}

func TestViewerCanAccessServiceHonorsOIDCOnlyRequirement(t *testing.T) {
	service := domain.Service{
		AccessMode:   domain.AccessModeAuthenticated,
		AccessMethod: domain.AccessMethodOIDCOnly,
	}

	localViewer := &domain.User{Role: domain.RoleViewer, Active: true, AuthProvider: domain.AuthProviderLocal}
	if viewerCanAccessService(localViewer, nil, service) {
		t.Fatal("expected local-auth viewer to be denied for oidc_only service")
	}

	oidcViewer := &domain.User{Role: domain.RoleViewer, Active: true, AuthProvider: domain.AuthProviderOIDC}
	if !viewerCanAccessService(oidcViewer, nil, service) {
		t.Fatal("expected oidc-authenticated viewer to be allowed")
	}
}
