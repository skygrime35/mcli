// internal/config/load_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFrom_CreatesDefaultWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcli", "config.yaml")

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cfg.Servers) != 0 {
		t.Errorf("expected empty servers list, got %d", len(cfg.Servers))
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file to be created at %s, got %v", path, err)
	}
}

func TestSaveTo_LoadFrom_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &Config{
		Servers: []ServerConfig{
			{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "romain", SSHPort: 22, WOLPort: 9},
		},
		Hotspot: HotspotConfig{SSID: "MyHotspot", Password: "secret"},
		Share:   ShareConfig{Dir: "share", Port: 8000},
	}

	if err := SaveTo(path, cfg); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if len(loaded.Servers) != 1 || loaded.Servers[0].Name != "home" || loaded.Servers[0].MAC != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("servers not preserved: %+v", loaded.Servers)
	}
	if loaded.Hotspot.SSID != "MyHotspot" {
		t.Errorf("hotspot ssid not preserved: %+v", loaded.Hotspot)
	}
}
