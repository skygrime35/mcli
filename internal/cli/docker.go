// internal/cli/docker.go
package cli

import (
	"fmt"
	"strings"

	"github.com/skygrime35/mcli/internal/docker"
	"github.com/spf13/cobra"
)

func newDockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Manage Docker containers, images, and volumes",
	}
	cmd.AddCommand(newDockerListCmd())
	cmd.AddCommand(newDockerFullPurgeCmd())
	cmd.AddCommand(newDockerClearAllCmd())
	cmd.AddCommand(newDockerSelectPurgeCmd())
	return cmd
}

func newDockerListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !docker.IsAvailable() {
				return fmt.Errorf("docker command not found")
			}
			containers, err := docker.ListContainers()
			if err != nil {
				return err
			}
			if len(containers) == 0 {
				fmt.Println("No containers found.")
				return nil
			}
			for _, c := range containers {
				id := c.ID
				if len(id) > 12 {
					id = id[:12]
				}
				fmt.Printf("%s\t%s\t%s\t%s\n", id, c.Name, c.Status, c.Image)
			}
			return nil
		},
	}
}

func runDockerProgress(ch <-chan docker.ProgressMsg) error {
	var firstErr error
	for msg := range ch {
		if msg.Err != nil {
			fmt.Println("Error:", msg.Err)
			if firstErr == nil {
				firstErr = msg.Err
			}
			continue
		}
		fmt.Println(msg.Text)
	}
	return firstErr
}

func newDockerFullPurgeCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "full-purge",
		Short: "Stop and remove ALL containers, images, volumes, networks, and build cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !docker.IsAvailable() {
				return fmt.Errorf("docker command not found")
			}
			if !yes {
				return fmt.Errorf("this is destructive - re-run with --yes to confirm")
			}
			return runDockerProgress(docker.FullPurge())
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm the destructive purge")
	return cmd
}

func newDockerClearAllCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "clear-all",
		Short: "Stop and remove all containers (images/volumes untouched)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !docker.IsAvailable() {
				return fmt.Errorf("docker command not found")
			}
			if !yes {
				return fmt.Errorf("this is destructive - re-run with --yes to confirm")
			}
			return runDockerProgress(docker.ClearAll())
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm the destructive purge")
	return cmd
}

func newDockerSelectPurgeCmd() *cobra.Command {
	var idsFlag string
	var yes bool
	cmd := &cobra.Command{
		Use:   "select-purge",
		Short: "Stop and remove specific containers by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !docker.IsAvailable() {
				return fmt.Errorf("docker command not found")
			}
			if idsFlag == "" {
				return fmt.Errorf("--ids is required (comma-separated container IDs)")
			}
			if !yes {
				return fmt.Errorf("this is destructive - re-run with --yes to confirm")
			}
			ids := strings.Split(idsFlag, ",")
			return runDockerProgress(docker.SelectPurge(ids))
		},
	}
	cmd.Flags().StringVar(&idsFlag, "ids", "", "comma-separated container IDs to purge (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm the destructive purge")
	return cmd
}
