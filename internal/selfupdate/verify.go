package selfupdate

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"regexp"
	"strings"
)

func VerifySHA256(payload []byte, checksumsTxt, assetName string) error {
	expected, err := lookupChecksum(checksumsTxt, assetName)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(payload)
	actual := fmt.Sprintf("%x", sum[:])
	if expected != actual {
		return fmt.Errorf("sha256 mismatch for %s: expected %s, got %s", assetName, expected, actual)
	}
	return nil
}

func lookupChecksum(checksumsTxt, assetName string) (string, error) {
	for _, line := range strings.Split(checksumsTxt, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[1] == assetName || strings.TrimPrefix(fields[1], "*") == assetName {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("no checksum entry for %s", assetName)
}

type CosignIdentity struct {
	IdentityRegex string
	OIDCIssuer    string
}

func VerifyCosign(blob []byte, signatureB64, certPEM string, identity CosignIdentity) error {
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(signatureB64))
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return fmt.Errorf("no PEM block in certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse certificate: %w", err)
	}
	if err := matchCertIdentity(cert, identity); err != nil {
		return err
	}
	pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("unsupported public key type %T (expected ecdsa)", cert.PublicKey)
	}
	digest := sha256.Sum256(blob)
	if !ecdsa.VerifyASN1(pub, digest[:], sig) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func matchCertIdentity(cert *x509.Certificate, identity CosignIdentity) error {
	if strings.TrimSpace(identity.IdentityRegex) == "" {
		return nil
	}
	re, err := regexp.Compile(identity.IdentityRegex)
	if err != nil {
		return fmt.Errorf("compile identity regex: %w", err)
	}
	candidates := certIdentityCandidates(cert)
	for _, candidate := range candidates {
		if re.MatchString(candidate) {
			if identity.OIDCIssuer == "" || certHasOIDCIssuer(cert, identity.OIDCIssuer) {
				return nil
			}
		}
	}
	return fmt.Errorf("certificate identity does not match %q (candidates: %v)", identity.IdentityRegex, candidates)
}

func certIdentityCandidates(cert *x509.Certificate) []string {
	var out []string
	for _, u := range cert.URIs {
		out = append(out, u.String())
	}
	out = append(out, cert.EmailAddresses...)
	out = append(out, cert.DNSNames...)
	return out
}

var oidcIssuerOID = []int{1, 3, 6, 1, 4, 1, 57264, 1, 1}

func certHasOIDCIssuer(cert *x509.Certificate, expected string) bool {
	expected = strings.TrimSpace(expected)
	for _, ext := range cert.Extensions {
		if !ext.Id.Equal(oidcIssuerOID) {
			continue
		}
		if strings.TrimSpace(string(ext.Value)) == expected {
			return true
		}
	}
	return false
}
