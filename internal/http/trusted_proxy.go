package http

import (
	"net"
	stdhttp "net/http"
	"net/netip"
	"strings"
)

func (s *Server) requestSecure(r *stdhttp.Request) bool {
	if r.TLS != nil {
		return true
	}
	if !s.requestFromTrustedProxy(r) {
		return false
	}
	return strings.EqualFold(firstForwardedValue(r.Header.Get("X-Forwarded-Proto")), "https")
}

func (s *Server) requestFromTrustedProxy(r *stdhttp.Request) bool {
	if r == nil || len(s.cfg.TrustedProxyCIDRs) == 0 {
		return false
	}
	addr, ok := remoteAddrFromRequest(r)
	if !ok {
		return false
	}
	for _, raw := range s.cfg.TrustedProxyCIDRs {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
		if err == nil && prefix.Contains(addr) {
			return true
		}
	}
	return false
}

func (s *Server) clientIPForRequest(r *stdhttp.Request) string {
	if s.requestFromTrustedProxy(r) {
		if forwarded := firstForwardedValue(r.Header.Get("X-Forwarded-For")); forwarded != "" {
			return forwarded
		}
		if realIP := strings.TrimSpace(r.Header.Get("X-Real-Ip")); realIP != "" {
			return realIP
		}
	}
	if host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func remoteAddrFromRequest(r *stdhttp.Request) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		host = strings.TrimSpace(r.RemoteAddr)
	}
	addr, err := netip.ParseAddr(host)
	return addr, err == nil
}

func firstForwardedValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}
