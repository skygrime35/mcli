// internal/system/analyze.go
package system

import (
	"fmt"

	"github.com/skygrime35/mcli/internal/platform"
)

// Analyze inspects the system (read-only: uname, dpkg --list - never
// mutates anything) and builds the Plan of actions that a subsequent
// Execute call would perform. Safe to call at any time, including
// repeatedly, with no side effects.
func Analyze(opts AnalyzeOptions, caps platform.Capabilities) (Plan, error) {
	if !caps.Apt {
		return Plan{}, fmt.Errorf("apt-get not found - System Update is only supported on Debian/Ubuntu-based systems")
	}

	var plan Plan
	if opts.Update {
		plan.Actions = append(plan.Actions,
			Action{Name: "Update package lists", Tier: TierSafe, Command: []string{"apt-get", "update"}},
			Action{Name: "Upgrade packages", Tier: TierSafe, Command: []string{"apt-get", "upgrade", "-y"}},
		)
	}

	if opts.Clean {
		plan.Actions = append(plan.Actions,
			Action{Name: "Remove unused packages", Tier: TierSafe, Command: []string{"apt-get", "autoremove", "-y"}},
			Action{Name: "Clean APT cache", Tier: TierSafe, Command: []string{"apt-get", "clean"}},
			Action{Name: "Deep clean obsolete packages", Tier: TierSafe, Command: []string{"apt-get", "autoremove", "--purge", "-y"}},
		)

		kernelAction, err := analyzeKernels()
		if err != nil {
			return plan, err
		}
		if kernelAction != nil {
			plan.Actions = append(plan.Actions, *kernelAction)
		}

		configAction, err := analyzeOrphanedConfigs()
		if err != nil {
			return plan, err
		}
		if configAction != nil {
			plan.Actions = append(plan.Actions, *configAction)
		}
	}

	return plan, nil
}
