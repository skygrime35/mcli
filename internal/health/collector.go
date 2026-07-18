package health

import (
	"time"

	"github.com/skygrime35/mcli/internal/platform"
)

func Collect(caps platform.Capabilities) Snapshot {
	return Snapshot{
		Timestamp: time.Now(),
		System:    CollectSystem(),
		CPU:       CollectCPU(caps),
		Memory:    CollectMemory(),
		Disks:     CollectDisks(),
		SMART:     CollectSMART(caps),
		Network:   CollectNetwork(),
		Services:  CollectServices(caps),
		Battery:   CollectBattery(caps),
		Errors:    CollectErrors(caps),
		Security:  CollectSecurity(caps),
	}
}
