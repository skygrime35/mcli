// internal/health/battery_test.go
package health

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadIntFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "capacity")
	if err := os.WriteFile(path, []byte("73\n"), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	v, ok := readIntFile(path)
	if !ok || v != 73 {
		t.Errorf("readIntFile() = (%d, %v), want (73, true)", v, ok)
	}
}

func TestReadIntFile_Missing(t *testing.T) {
	if _, ok := readIntFile("/nonexistent/path/capacity"); ok {
		t.Error("expected ok=false for a missing file")
	}
}

func TestReadFloatFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "energy_full")
	if err := os.WriteFile(path, []byte("48000000\n"), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	v, ok := readFloatFile(path)
	if !ok || v != 48000000 {
		t.Errorf("readFloatFile() = (%v, %v), want (48000000, true)", v, ok)
	}
}
