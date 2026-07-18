// internal/config/servers.go
package config

import (
	"fmt"
	"regexp"
)

var macRegexp = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)

func ValidateServer(s ServerConfig) error {
	if s.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	if s.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if !macRegexp.MatchString(s.MAC) {
		return fmt.Errorf("invalid MAC address: %q", s.MAC)
	}
	if s.SSHPort < 1 || s.SSHPort > 65535 {
		return fmt.Errorf("invalid ssh_port: %d", s.SSHPort)
	}
	if s.WOLPort < 1 || s.WOLPort > 65535 {
		return fmt.Errorf("invalid wol_port: %d", s.WOLPort)
	}
	return nil
}

func AddServer(cfg *Config, entry ServerConfig) error {
	if err := ValidateServer(entry); err != nil {
		return err
	}
	for _, s := range cfg.Servers {
		if s.Name == entry.Name {
			return fmt.Errorf("a server named %q already exists", entry.Name)
		}
	}
	cfg.Servers = append(cfg.Servers, entry)
	return nil
}

func FindServer(cfg *Config, name string) (ServerConfig, error) {
	for _, s := range cfg.Servers {
		if s.Name == name {
			return s, nil
		}
	}
	return ServerConfig{}, fmt.Errorf("no server named %q in config", name)
}
