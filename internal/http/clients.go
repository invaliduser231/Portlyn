package http

import (
	"errors"
	stdhttp "net/http"
	"net/netip"
	"strings"

	"portlyn/internal/domain"
	"portlyn/internal/tunnel"
)

type createClientRequest struct {
	Name           string `json:"name" validate:"required,max=255"`
	Description    string `json:"description,omitempty" validate:"omitempty,max=1024"`
	AllowedNodeIDs []uint `json:"allowed_node_ids,omitempty"`
}

type clientConfigResponse struct {
	Client     *domain.Client `json:"client"`
	ConfigText string         `json:"config_text"`
	Address    string         `json:"address"`
	AllowedIPs []string       `json:"allowed_ips"`
}

func (s *Server) handleListClients(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if s.clients == nil {
		writeJSON(w, stdhttp.StatusOK, []domain.Client{})
		return
	}
	items, err := s.clients.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, items)
}

func (s *Server) handleCreateClient(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if s.tunnel == nil {
		writeError(w, stdhttp.StatusServiceUnavailable, "tunnel_unavailable", "tunnel manager is not configured")
		return
	}
	var req createClientRequest
	if !s.decodeAndValidate(w, r, &req) {
		return
	}
	bundle, client, err := s.tunnel.ProvisionClient(r.Context(), req.Name, req.Description, req.AllowedNodeIDs)
	if err != nil {
		if errors.Is(err, tunnel.ErrInvalidServerSettings) {
			writeError(w, stdhttp.StatusPreconditionFailed, "tunnel_not_configured", "enable and configure the tunnel server first")
			return
		}
		if errors.Is(err, tunnel.ErrPoolExhausted) {
			writeError(w, stdhttp.StatusConflict, "pool_exhausted", "tunnel ip pool is exhausted")
			return
		}
		s.internalError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "create", "client", &client.ID, map[string]any{
		"name":      client.Name,
		"tunnel_ip": client.WGTunnelIP,
	})
	writeJSON(w, stdhttp.StatusCreated, clientConfigResponse{
		Client:     client,
		ConfigText: tunnel.RenderClientConfig(*bundle),
		Address:    bundle.Address,
		AllowedIPs: bundle.AllowedIPs,
	})
}

func (s *Server) handleRotateClient(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if s.tunnel == nil {
		writeError(w, stdhttp.StatusServiceUnavailable, "tunnel_unavailable", "tunnel manager is not configured")
		return
	}
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	bundle, client, err := s.tunnel.RotateClient(r.Context(), id)
	if err != nil {
		s.handleStoreError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "rotate", "client", &client.ID, nil)
	writeJSON(w, stdhttp.StatusOK, clientConfigResponse{
		Client:     client,
		ConfigText: tunnel.RenderClientConfig(*bundle),
		Address:    bundle.Address,
		AllowedIPs: bundle.AllowedIPs,
	})
}

func (s *Server) handleDeleteClient(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if s.tunnel == nil {
		writeError(w, stdhttp.StatusServiceUnavailable, "tunnel_unavailable", "tunnel manager is not configured")
		return
	}
	id, ok := s.parseIDParam(w, r, "id")
	if !ok {
		return
	}
	if err := s.tunnel.RevokeClient(r.Context(), id); err != nil {
		s.handleStoreError(w, err)
		return
	}
	_ = s.audit.LogRequest(r.Context(), r, s.currentUserID(r), "delete", "client", &id, map[string]any{"id": id})
	w.WriteHeader(stdhttp.StatusNoContent)
}

func normalizeSubnetCSV(value string) (string, error) {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(trimmed)
		if err != nil {
			return "", err
		}
		out = append(out, prefix.Masked().String())
	}
	return strings.Join(out, ","), nil
}
