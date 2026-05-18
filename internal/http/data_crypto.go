package http

import (
	"portlyn/internal/secureconfig"
)

func (s *Server) dataEncryptJSON(value map[string]string) (string, error) {
	return secureconfig.EncryptJSON([]byte(s.cfg.DataEncryptionSecret), value)
}

func (s *Server) dataDecryptJSON(value string) (map[string]string, error) {
	secrets := s.dataSecrets()
	return secureconfig.DecryptJSONWithSecrets(secrets, value)
}

func (s *Server) dataDecryptJSONWithActiveKey(value string) (map[string]string, error) {
	return secureconfig.DecryptJSON([]byte(s.cfg.DataEncryptionSecret), value)
}

func (s *Server) dataSecrets() [][]byte {
	values := s.cfg.DataEncryptionSecrets()
	out := make([][]byte, 0, len(values))
	for _, value := range values {
		out = append(out, []byte(value))
	}
	return out
}
