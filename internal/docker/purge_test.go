// internal/docker/purge_test.go
package docker

import "testing"

func TestParseCustomNetworks(t *testing.T) {
	output := `abc123 bridge
def456 host
ghi789 none
jkl012 my_custom_network
mno345 another_custom_net
`
	got := parseCustomNetworks(output)
	want := []string{"jkl012", "mno345"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseCustomNetworks_OnlyDefaults(t *testing.T) {
	output := "abc123 bridge\ndef456 host\nghi789 none\n"
	if got := parseCustomNetworks(output); len(got) != 0 {
		t.Errorf("expected no custom networks, got %v", got)
	}
}

func TestDedup(t *testing.T) {
	got := dedup([]string{"a", "b", "a", "c", "b", "a"})
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDedup_Empty(t *testing.T) {
	if got := dedup(nil); len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
