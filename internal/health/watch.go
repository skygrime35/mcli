package health

import (
	"context"
	"time"

	"github.com/skygrime35/mcli/internal/platform"
)

// Watch emits a fresh Snapshot every interval until ctx is canceled, then
// closes the channel. The caller drives redraw/print behavior per
// snapshot - this function only owns the collection cadence.
func Watch(ctx context.Context, caps platform.Capabilities, interval time.Duration) <-chan Snapshot {
	ch := make(chan Snapshot)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		select {
		case ch <- Collect(caps):
		case <-ctx.Done():
			return
		}

		for {
			select {
			case <-ticker.C:
				select {
				case ch <- Collect(caps):
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
