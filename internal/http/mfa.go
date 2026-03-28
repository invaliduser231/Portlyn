package http

import (
	"net/http"

	"portlyn/internal/auth"
)

func (s *Server) handleGetMyMFAStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	status, err := s.auth.MFAStatusForUser(r.Context(), user.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleBeginMyMFASetup(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	setup, err := s.auth.BeginTOTPSetup(r.Context(), user.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "mfa_setup_started", "user", &user.ID, map[string]any{"user_id": user.ID})
	writeJSON(w, http.StatusOK, setup)
}

func (s *Server) handleEnableMyMFA(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	var req mfaCodeRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	status, err := s.auth.EnableTOTP(r.Context(), user.ID, req.Code)
	if err != nil {
		if err == auth.ErrMFACodeInvalid {
			writeError(w, http.StatusUnauthorized, "invalid_mfa_code", "invalid authenticator code")
			return
		}
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "mfa_enabled", "user", &user.ID, map[string]any{"user_id": user.ID})
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleDisableMyMFA(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	var req mfaCodeRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	status, err := s.auth.DisableMFA(r.Context(), user.ID, req.Code)
	if err != nil {
		if err == auth.ErrMFACodeInvalid {
			writeError(w, http.StatusUnauthorized, "invalid_mfa_code", "invalid authenticator or recovery code")
			return
		}
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "mfa_disabled", "user", &user.ID, map[string]any{"user_id": user.ID})
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleRegenerateMyRecoveryCodes(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	var req mfaCodeRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	codes, err := s.auth.RegenerateRecoveryCodes(r.Context(), user.ID, req.Code)
	if err != nil {
		if err == auth.ErrMFACodeInvalid {
			writeError(w, http.StatusUnauthorized, "invalid_mfa_code", "invalid authenticator or recovery code")
			return
		}
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "mfa_recovery_regenerated", "user", &user.ID, map[string]any{"user_id": user.ID})
	writeJSON(w, http.StatusOK, map[string]any{"recovery_codes": codes})
}

func (s *Server) handleResetUserMFA(w http.ResponseWriter, r *http.Request) {
	user, ok := s.loadUser(w, r)
	if !ok {
		return
	}
	if err := s.auth.ResetUserMFA(r.Context(), user.ID); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "mfa_reset", "user", &user.ID, map[string]any{"user_id": user.ID, "email": user.Email})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
