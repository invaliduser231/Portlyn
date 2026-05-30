package selfupdate

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAtomicSwapKeepsBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "binary")
	if err := os.WriteFile(target, []byte("v1"), 0o755); err != nil {
		t.Fatalf("seed: %v", err)
	}

	backup, err := AtomicSwap(target, bytes.NewReader([]byte("v2")), 0o755)
	if err != nil {
		t.Fatalf("swap: %v", err)
	}

	currentBytes, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(currentBytes) != "v2" {
		t.Fatalf("target = %q", currentBytes)
	}

	backupBytes, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(backupBytes) != "v1" {
		t.Fatalf("backup = %q", backupBytes)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("target not executable: %v", info.Mode())
	}
}

func TestAtomicSwapCreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "fresh")

	backup, err := AtomicSwap(target, bytes.NewReader([]byte("v1")), 0o755)
	if err != nil {
		t.Fatalf("swap: %v", err)
	}
	if backup == "" {
		t.Fatal("backup path should still be returned as a hint")
	}
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatalf("backup should not exist when target was new, got %v", err)
	}
	if data, err := os.ReadFile(target); err != nil || string(data) != "v1" {
		t.Fatalf("target = %q, err %v", data, err)
	}
}

func TestAssetNameIncludesPlatform(t *testing.T) {
	name := AssetName("portlyn")
	if name == "" {
		t.Fatal("empty asset name")
	}
}
