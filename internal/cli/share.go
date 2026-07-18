// internal/cli/share.go
package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skygrime35/mcli/internal/share"
	"github.com/spf13/cobra"
)

func newShareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Share a local directory over HTTP",
	}
	cmd.AddCommand(newShareStartCmd())
	return cmd
}

func newShareStartCmd() *cobra.Command {
	var dirFlag string
	var portFlag int
	var passwordFlag string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the file-share HTTP server (Ctrl+C to stop)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			dir := cfg.Share.Dir
			if cmd.Flags().Changed("dir") {
				dir = dirFlag
			}
			port := cfg.Share.Port
			if cmd.Flags().Changed("port") {
				port = portFlag
			}
			password := cfg.Share.Password
			if cmd.Flags().Changed("password") {
				password = passwordFlag
			}

			srv, err := share.Start(dir, port, password)
			if err != nil {
				return err
			}

			fmt.Printf("Sharing '%s' at: http://%s:%d\n", dir, share.LocalIP(), port)
			if password != "" {
				fmt.Println("Password required (username is ignored).")
			} else {
				fmt.Println("No password set - open access.")
			}
			fmt.Println("Press Ctrl+C to stop sharing.")

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			<-sigCh

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			fmt.Println("\nStopping...")
			return srv.Stop(ctx)
		},
	}

	cmd.Flags().StringVar(&dirFlag, "dir", "", "directory to share (defaults to config)")
	cmd.Flags().IntVar(&portFlag, "port", 0, "port to listen on (defaults to config)")
	cmd.Flags().StringVar(&passwordFlag, "password", "", "password required to download (defaults to config; empty = open access)")

	return cmd
}
