package http

import (
	stdhttp "net/http"
	"time"
)

func (s *Server) handleRetryCertificate(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	item, ok := s.loadCertificate(w, r)
	if !ok {
		return
	}
	item.Status = "pending"
	item.LastError = ""
	if _, err := s.acme.SyncCertificate(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "retry", "certificate", &item.ID, map[string]any{"certificate_id": item.ID})
	writeJSON(w, stdhttp.StatusOK, item)
}

func (s *Server) handleRenewCertificate(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	item, ok := s.loadCertificate(w, r)
	if !ok {
		return
	}
	item.Status = "renewing"
	item.LastCheckedAt = ptrTime(time.Now().UTC())
	if _, err := s.acme.SyncCertificate(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "renew", "certificate", &item.ID, map[string]any{"certificate_id": item.ID})
	writeJSON(w, stdhttp.StatusOK, item)
}

func (s *Server) handleSyncCertificateStatus(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	item, ok := s.loadCertificate(w, r)
	if !ok {
		return
	}
	if _, err := s.acme.SyncCertificate(r.Context(), item); err != nil {
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "sync_status", "certificate", &item.ID, map[string]any{"certificate_id": item.ID})
	writeJSON(w, stdhttp.StatusOK, item)
}

func ptrTime(value time.Time) *time.Time {
	return &value
}
