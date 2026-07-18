// internal/cli/server.go
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/server"
	"github.com/spf13/cobra"
)

func newServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage remote servers (Wake-on-LAN, SSH)",
	}
	cmd.AddCommand(newServerListCmd())
	cmd.AddCommand(newServerAddCmd())
	cmd.AddCommand(newServerStatusCmd())
	cmd.AddCommand(newServerWOLCmd())
	cmd.AddCommand(newServerSSHCmd())
	cmd.AddCommand(newServerSSHKeyCmd())
	return cmd
}

func newServerListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if len(cfg.Servers) == 0 {
				fmt.Println("No servers configured. Use `mcli server add` to add one.")
				return nil
			}
			for _, s := range cfg.Servers {
				fmt.Printf("%s\t%s@%s:%d (MAC %s, WOL port %d)\n", s.Name, s.SSHUser, s.Host, s.SSHPort, s.MAC, s.WOLPort)
			}
			return nil
		},
	}
}

func newServerAddCmd() *cobra.Command {
	var name, host, mac, sshUser string
	var sshPort, wolPort int
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a server to the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			entry := config.ServerConfig{
				Name:    name,
				Host:    host,
				MAC:     mac,
				SSHUser: sshUser,
				SSHPort: sshPort,
				WOLPort: wolPort,
			}
			if err := config.AddServer(cfg, entry); err != nil {
				return err
			}
			if err := config.Save(cfg); err != nil {
				return err
			}
			fmt.Printf("Server %q added.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "server name (required)")
	cmd.Flags().StringVar(&host, "host", "", "server host/IP (required)")
	cmd.Flags().StringVar(&mac, "mac", "", "server MAC address (required)")
	cmd.Flags().StringVar(&sshUser, "ssh-user", "", "SSH username (required)")
	cmd.Flags().IntVar(&sshPort, "ssh-port", 22, "SSH port")
	cmd.Flags().IntVar(&wolPort, "wol-port", 9, "Wake-on-LAN port")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("mac")
	cmd.MarkFlagRequired("ssh-user")
	return cmd
}

func newServerStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <name>",
		Short: "Check if a server is online (SSH port reachable)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			s, err := config.FindServer(cfg, args[0])
			if err != nil {
				return err
			}
			status := server.CheckStatus(context.Background(), s.Host, s.SSHPort, 5*time.Second)
			if status.Online {
				fmt.Printf("%s is online.\n", s.Name)
			} else {
				fmt.Printf("%s is offline.\n", s.Name)
			}
			return nil
		},
	}
}

func newServerWOLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wol <name>",
		Short: "Wake a server via Wake-on-LAN and wait for it to come online",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			s, err := config.FindServer(cfg, args[0])
			if err != nil {
				return err
			}
			for msg := range server.WakeOnLAN(context.Background(), s) {
				if msg.Err != nil {
					return msg.Err
				}
				fmt.Println(msg.Text)
			}
			return nil
		},
	}
}

func newServerSSHCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ssh <name>",
		Short: "Open an interactive SSH session to a server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			s, err := config.FindServer(cfg, args[0])
			if err != nil {
				return err
			}
			return server.Connect(s)
		},
	}
}

func newServerSSHKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ssh-key <name>",
		Short: "Generate (if needed) and copy an SSH key to a server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			s, err := config.FindServer(cfg, args[0])
			if err != nil {
				return err
			}
			for msg := range server.SetupSSHKey(s) {
				if msg.Err != nil {
					return msg.Err
				}
				fmt.Println(msg.Text)
			}
			return nil
		},
	}
}
