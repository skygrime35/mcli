// internal/tui/hotspotui/menu.go
package hotspotui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/hotspot"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

func credentials() (ssid, password string) {
	cfg, err := config.Load()
	ssid, password = "Hotspot", "123456789"
	if err == nil {
		if cfg.Hotspot.SSID != "" {
			ssid = cfg.Hotspot.SSID
		}
		if cfg.Hotspot.Password != "" {
			password = cfg.Hotspot.Password
		}
	}
	return ssid, password
}

// NewMenuScreen is the Hotspot Manager's entry screen. It adapts to the
// hotspot's current state, matching the old app.py's behavior: shows
// "Activate" when inactive, "Deactivate" when active - never both at
// once - plus "Show Statistics" always. The screen's title stays
// constant ("Hotspot Manager"), matching the convention already
// established by dockerui/systemui's menu screens.
func NewMenuScreen() tui.Screen {
	if !platform.Detect().Nmcli {
		return menu.New("Hotspot Manager", []menu.Item{
			{Title: "nmcli not found", Disabled: true},
		})
	}

	ssid, password := credentials()
	active := hotspot.IsActive(ssid)

	var items []menu.Item
	if active {
		items = append(items, menu.Item{
			Title:       "Deactivate Hotspot",
			Description: fmt.Sprintf("Bring down '%s'", ssid),
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newProgressScreen("Deactivate Hotspot", hotspot.Deactivate(ssid)))
			},
		})
	} else {
		items = append(items, menu.Item{
			Title:       "Activate Hotspot",
			Description: fmt.Sprintf("Bring up '%s'", ssid),
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newProgressScreen("Activate Hotspot", hotspot.Activate(ssid, password)))
			},
		})
	}
	items = append(items, menu.Item{
		Title:       "Show Hotspot Statistics",
		Description: "Active interface, connected clients, data usage",
		OnSelect: func() tea.Cmd {
			return tui.PushScreen(newStatsScreen(ssid))
		},
	})

	return menu.New("Hotspot Manager", items)
}
