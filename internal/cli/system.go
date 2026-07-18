// internal/cli/system.go
package cli

import (
	"fmt"

	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/system"
	"github.com/spf13/cobra"
)

func newSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System updates and cleanup (apt-based)",
	}
	cmd.AddCommand(newSystemUpdateCmd())
	return cmd
}

func newSystemUpdateCmd() *cobra.Command {
	var (
		yes          bool
		updateFlag   bool
		cleanFlag    bool
		doWarnings   bool
		doUnsafe     bool
		skipWarnings bool
		skipUnsafe   bool
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Analyze (always) and optionally apply system updates and cleanup",
		RunE: func(cmd *cobra.Command, args []string) error {
			caps := platform.Detect()
			plan, err := system.Analyze(system.AnalyzeOptions{Update: updateFlag, Clean: cleanFlag}, caps)
			if err != nil {
				return err
			}

			fmt.Println("Planned actions:")
			for _, a := range plan.Actions {
				fmt.Printf("  [%s] %s\n", a.Tier, a.Name)
				if a.Reason != "" {
					fmt.Printf("        %s\n", a.Reason)
				}
			}

			if !yes {
				fmt.Println("\nDry run only - re-run with --yes to apply.")
				return nil
			}

			opts := system.ExecuteOptions{
				DoWarnings:   doWarnings,
				DoUnsafe:     doUnsafe,
				SkipWarnings: skipWarnings,
				SkipUnsafe:   skipUnsafe,
			}
			fmt.Println()
			return system.Execute(plan, opts, caps)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "actually apply the plan (default is dry-run/analysis only)")
	cmd.Flags().BoolVar(&updateFlag, "update", true, "include update package lists/upgrade packages")
	cmd.Flags().BoolVar(&cleanFlag, "clean", true, "include autoremove/clean/kernel+config cleanup")
	cmd.Flags().BoolVar(&doWarnings, "do-warnings", false, "auto-approve Warning-tier actions")
	cmd.Flags().BoolVar(&doUnsafe, "do-unsafe", false, "auto-approve Unsafe-tier actions")
	cmd.Flags().BoolVar(&skipWarnings, "skip-warnings", false, "skip Warning-tier actions entirely")
	cmd.Flags().BoolVar(&skipUnsafe, "skip-unsafe", false, "skip Unsafe-tier actions entirely")
	return cmd
}
