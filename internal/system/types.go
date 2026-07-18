// internal/system/types.go
package system

type Tier string

const (
	TierSafe    Tier = "safe"
	TierWarning Tier = "warning"
	TierUnsafe  Tier = "unsafe"
)

// Action is one apt/dpkg operation, classified by risk tier. Command is
// the binary and its arguments as discrete slice elements - never a shell
// string (see the plan's Global Constraints on avoiding eval-style
// execution).
type Action struct {
	Name    string
	Tier    Tier
	Command []string
	Reason  string
}

type Plan struct {
	Actions []Action
}

// AnalyzeOptions selects which categories of actions Analyze considers.
type AnalyzeOptions struct {
	Update bool
	Clean  bool
}

// ExecuteOptions controls how Warning/Unsafe actions are handled during
// Execute - Safe actions always run regardless of these options.
type ExecuteOptions struct {
	DoWarnings   bool
	DoUnsafe     bool
	SkipWarnings bool
	SkipUnsafe   bool
}
