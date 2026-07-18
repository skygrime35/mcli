// internal/docker/clear.go
package docker

import "fmt"

// ClearAll stops and removes every container, leaving images/volumes
// untouched - matching cleaner.py's clear_all_containers().
func ClearAll() <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		if !CheckDaemon() {
			ch <- ProgressMsg{Err: fmt.Errorf("docker daemon is not accessible - is Docker running?")}
			return
		}
		if stopAllContainers(ch) != nil {
			return
		}
		if removeAllContainers(ch) != nil {
			return
		}
		ch <- ProgressMsg{Text: "All containers cleared."}
	}()
	return ch
}
