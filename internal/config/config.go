// internal/config/config.go
package config

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
	Hotspot HotspotConfig  `yaml:"hotspot"`
	Share   ShareConfig    `yaml:"share"`
}

type ServerConfig struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	MAC     string `yaml:"mac"`
	SSHUser string `yaml:"ssh_user"`
	SSHPort int    `yaml:"ssh_port"`
	WOLPort int    `yaml:"wol_port"`
}

type HotspotConfig struct {
	SSID     string `yaml:"ssid"`
	Password string `yaml:"password"`
}

type ShareConfig struct {
	Dir      string `yaml:"dir"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password,omitempty"`
}

func Default() Config {
	return Config{
		Servers: []ServerConfig{},
		Hotspot: HotspotConfig{},
		Share:   ShareConfig{Dir: "share", Port: 8000},
	}
}
