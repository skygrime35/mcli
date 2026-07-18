// internal/docker/types.go
package docker

// Container is one row from `docker ps -a`.
type Container struct {
	ID     string
	Name   string
	Status string
	Image  string
}

// ProgressMsg streams a line of progress or a terminal error from a
// long-running Docker operation (FullPurge/ClearAll/SelectPurge), read the
// same way as internal/server.ProgressMsg but deliberately not shared
// across packages - see the plan's Global Constraints.
type ProgressMsg struct {
	Text string
	Err  error
}
