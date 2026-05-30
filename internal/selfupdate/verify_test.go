package selfupdate

import (
	"strings"
	"testing"
)

func TestVerifySHA256Success(t *testing.T) {
	payload := []byte("hello world\n")
	checksums := "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447  hello.txt\n"
	if err := VerifySHA256(payload, checksums, "hello.txt"); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

func TestVerifySHA256MissingEntry(t *testing.T) {
	err := VerifySHA256([]byte("x"), "deadbeef  other.txt\n", "missing.txt")
	if err == nil || !strings.Contains(err.Error(), "no checksum entry") {
		t.Fatalf("expected no-entry error, got %v", err)
	}
}

func TestVerifySHA256Mismatch(t *testing.T) {
	checksums := "0000000000000000000000000000000000000000000000000000000000000000  hello.txt\n"
	err := VerifySHA256([]byte("hello world\n"), checksums, "hello.txt")
	if err == nil || !strings.Contains(err.Error(), "sha256 mismatch") {
		t.Fatalf("expected mismatch, got %v", err)
	}
}

func TestVerifySHA256IgnoresAsteriskPrefix(t *testing.T) {
	payload := []byte("hello world\n")
	checksums := "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447  *hello.txt\n"
	if err := VerifySHA256(payload, checksums, "hello.txt"); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

func TestVerifyCosignRejectsBadSignature(t *testing.T) {
	err := VerifyCosign([]byte("data"), "not-base64!!", "no-pem", CosignIdentity{})
	if err == nil {
		t.Fatal("expected error on bad signature input")
	}
}

func TestVerifyCosignRejectsBadCert(t *testing.T) {
	err := VerifyCosign([]byte("data"), "aGVsbG8=", "not-a-pem", CosignIdentity{})
	if err == nil || !strings.Contains(err.Error(), "no PEM block") {
		t.Fatalf("expected PEM error, got %v", err)
	}
}
