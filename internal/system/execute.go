// internal/system/execute.go
package system

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/skygrime35/mcli/internal/platform"
)

// ShouldRun is the pure decision logic for whether an action executes,
// given the caller's approval flags. Safe actions always run. Warning
// and Unsafe actions require their corresponding "do" flag, and a "skip"
// flag always takes precedence over a "do" flag for the same tier.
func ShouldRun(action Action, opts ExecuteOptions) (run bool, note string) {
	switch action.Tier {
	case TierSafe:
		return true, ""
	case TierWarning:
		if opts.SkipWarnings {
			return false, " (skipped: warning)"
		}
		if !opts.DoWarnings {
			return false, fmt.Sprintf(" (skipped: warning not approved - %s)", action.Reason)
		}
		return true, ""
	case TierUnsafe:
		if opts.SkipUnsafe {
			return false, " (skipped: unsafe)"
		}
		if !opts.DoUnsafe {
			return false, fmt.Sprintf(" (skipped: unsafe not approved - %s)", action.Reason)
		}
		return true, ""
	default:
		return false, " (skipped: unknown tier)"
	}
}

// Execute runs every action in plan that ShouldRun approves, in order,
// with the action's command prefixed by "sudo" and stdio connected
// directly to the real terminal (interactive sudo password prompts and
// apt's own live output both need a real terminal - this deliberately
// does NOT use the channel-based ProgressMsg pattern from the
// server/docker plans, see this plan's architecture notes). Stops and
// returns an error on the first action that actually fails to run.
func Execute(plan Plan, opts ExecuteOptions, caps platform.Capabilities) error {
	if !caps.Apt {
		return fmt.Errorf("apt-get not found - System Update is only supported on Debian/Ubuntu-based systems")
	}
	for _, action := range plan.Actions {
		run, note := ShouldRun(action, opts)
		fmt.Printf("-> %s%s\n", action.Name, note)
		if !run {
			continue
		}
		args := append([]string{action.Command[0]}, action.Command[1:]...)
		cmd := exec.Command("sudo", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", action.Name, err)
		}
		fmt.Println("   done.")
	}
	return nil
}
