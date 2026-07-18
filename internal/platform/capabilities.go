package platform

import "os/exec"

type Capabilities struct {
	Docker     bool
	Systemd    bool
	Nmcli      bool
	Apt        bool
	Smartctl   bool
	Upower     bool
	Journalctl bool
	Sensors    bool
	Ufw        bool
	Iptables   bool
}

type lookupFunc func(string) (string, error)

func Detect() Capabilities {
	return detect(exec.LookPath)
}

func detect(lookPath lookupFunc) Capabilities {
	has := func(bin string) bool {
		_, err := lookPath(bin)
		return err == nil
	}
	return Capabilities{
		Docker:     has("docker"),
		Systemd:    has("systemctl"),
		Nmcli:      has("nmcli"),
		Apt:        has("apt-get"),
		Smartctl:   has("smartctl"),
		Upower:     has("upower"),
		Journalctl: has("journalctl"),
		Sensors:    has("sensors"),
		Ufw:        has("ufw"),
		Iptables:   has("iptables"),
	}
}
