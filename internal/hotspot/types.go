// internal/hotspot/types.go
package hotspot

// ProgressMsg streams a line of progress or a terminal error from
// Activate/Deactivate - a small local type, matching the pattern already
// established (and deliberately duplicated, not shared) in the docker
// and server plans.
type ProgressMsg struct {
	Text string
	Err  error
}

// Stats reports the current state of a hotspot connection.
type Stats struct {
	Active    bool
	Interface string
	Clients   []string
	RXBytes   uint64
	TXBytes   uint64
}
