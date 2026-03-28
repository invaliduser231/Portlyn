package http

import (
	stdhttp "net/http"
	"time"

	"portlyn/internal/auth"
)

func (s *Server) handleListMySessions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, stdhttp.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	items, err := s.auth.ListUserSessions(r.Context(), user.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, items)
}

func (s *Server) handleRevokeMySession(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, stdhttp.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	sessionID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	if err := s.auth.RevokeSession(r.Context(), user.ID, sessionID); err != nil {
		s.handleStoreError(w, err)
		return
	}
	token := s.auth.SessionTokenFromRequest(r)
	if claims, err := s.auth.ParseToken(token); err == nil && claims.SessionID == sessionID {
		s.auth.ClearSessionCookie(w)
		s.auth.ClearRefreshCookie(w)
	}
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "session_revoked", "session", &sessionID, map[string]any{"session_id": sessionID})
	writeJSON(w, stdhttp.StatusOK, map[string]any{"ok": true, "revoked_at": time.Now().UTC()})
}

func (s *Server) handleRevokeAllMySessions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, stdhttp.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}
	if err := s.auth.RevokeAllUserSessions(r.Context(), user.ID); err != nil {
		s.internalError(w, err)
		return
	}
	s.auth.ClearSessionCookie(w)
	s.auth.ClearRefreshCookie(w)
	_ = s.audit.LogRequest(r.Context(), r, &user.ID, "session_revoke_all", "session", nil, map[string]any{"user_id": user.ID})
	writeJSON(w, stdhttp.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleListUserSessions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	items, err := s.auth.ListUserSessions(r.Context(), userID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, items)
}

func (s *Server) handleRevokeUserSession(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	sessionID, ok := s.parseIDParam(w, r, "sessionId")
	if !ok {
		return
	}
	if err := s.auth.RevokeSession(r.Context(), userID, sessionID); err != nil {
		s.handleStoreError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "session_revoked", "session", &sessionID, map[string]any{"user_id": userID, "session_id": sessionID})
	writeJSON(w, stdhttp.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleRevokeAllUserSessions(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	if err := s.auth.RevokeAllUserSessions(r.Context(), userID); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "session_revoke_all", "session", nil, map[string]any{"user_id": userID})
	writeJSON(w, stdhttp.StatusOK, map[string]any{"ok": true})
}
