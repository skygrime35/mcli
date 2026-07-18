package platform

import (
	"errors"
	"testing"
)

func TestDetect_UsesLookupFunc(t *testing.T) {
	fakeLookPath := func(bin string) (string, error) {
		available := map[string]bool{
			"docker":  true,
			"nmcli":   true,
			"apt-get": false,
		}
		if available[bin] {
			return "/usr/bin/" + bin, nil
		}
		return "", errors.New("not found")
	}

	caps := detect(fakeLookPath)

	if !caps.Docker {
		t.Error("expected Docker to be true")
	}
	if !caps.Nmcli {
		t.Error("expected Nmcli to be true")
	}
	if caps.Apt {
		t.Error("expected Apt to be false")
	}
	if caps.Systemd {
		t.Error("expected Systemd to be false (not in fake lookup map)")
	}
}

func TestDetect_IncludesSensorsUfwIptables(t *testing.T) {
	fakeLookPath := func(bin string) (string, error) {
		available := map[string]bool{
			"sensors":  true,
			"ufw":      false,
			"iptables": true,
		}
		if available[bin] {
			return "/usr/bin/" + bin, nil
		}
		return "", errors.New("not found")
	}

	caps := detect(fakeLookPath)

	if !caps.Sensors {
		t.Error("expected Sensors to be true")
	}
	if caps.Ufw {
		t.Error("expected Ufw to be false")
	}
	if !caps.Iptables {
		t.Error("expected Iptables to be true")
	}
}
