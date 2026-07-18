// internal/config/servers_test.go
package config

import "testing"

func TestValidateServer(t *testing.T) {
	valid := ServerConfig{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "romain", SSHPort: 22, WOLPort: 9}
	if err := ValidateServer(valid); err != nil {
		t.Fatalf("expected valid server to pass, got %v", err)
	}

	cases := []struct {
		name string
		s    ServerConfig
	}{
		{"empty name", ServerConfig{Host: "h", MAC: "AA:BB:CC:DD:EE:FF", SSHPort: 22, WOLPort: 9}},
		{"empty host", ServerConfig{Name: "n", MAC: "AA:BB:CC:DD:EE:FF", SSHPort: 22, WOLPort: 9}},
		{"bad mac", ServerConfig{Name: "n", Host: "h", MAC: "not-a-mac", SSHPort: 22, WOLPort: 9}},
		{"bad ssh port", ServerConfig{Name: "n", Host: "h", MAC: "AA:BB:CC:DD:EE:FF", SSHPort: 0, WOLPort: 9}},
		{"bad wol port", ServerConfig{Name: "n", Host: "h", MAC: "AA:BB:CC:DD:EE:FF", SSHPort: 22, WOLPort: 70000}},
	}
	for _, c := range cases {
		if err := ValidateServer(c.s); err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
	}
}

func TestAddServer(t *testing.T) {
	cfg := &Config{}
	entry := ServerConfig{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "romain", SSHPort: 22, WOLPort: 9}

	if err := AddServer(cfg, entry); err != nil {
		t.Fatalf("expected AddServer to succeed, got %v", err)
	}
	if len(cfg.Servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(cfg.Servers))
	}

	if err := AddServer(cfg, entry); err == nil {
		t.Fatal("expected error when adding duplicate server name, got nil")
	}
}

func TestFindServer(t *testing.T) {
	cfg := &Config{Servers: []ServerConfig{{Name: "home", Host: "example.com"}}}

	found, err := FindServer(cfg, "home")
	if err != nil {
		t.Fatalf("expected to find server, got error %v", err)
	}
	if found.Host != "example.com" {
		t.Errorf("expected host example.com, got %s", found.Host)
	}

	if _, err := FindServer(cfg, "missing"); err == nil {
		t.Fatal("expected error for missing server, got nil")
	}
}
