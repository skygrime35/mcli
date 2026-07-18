// internal/docker/client.go
package docker

import "os/exec"

// IsAvailable reports whether the docker binary is present on PATH.
func IsAvailable() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// CheckDaemon reports whether the Docker daemon is actually reachable
// (docker info succeeds) - unlike the old Python reference, this IS
// checked before every destructive operation in this package (see
// Global Constraints).
func CheckDaemon() bool {
	return exec.Command("docker", "info").Run() == nil
}
