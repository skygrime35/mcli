// internal/hotspot/status_test.go
package hotspot

import "testing"

func TestParseActiveConnectionNames(t *testing.T) {
	output := "Le Super Telephone\ndocker0\nlo\n"
	got := parseActiveConnectionNames(output)
	want := []string{"Le Super Telephone", "docker0", "lo"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestContainsName(t *testing.T) {
	names := []string{"MyHotspot", "docker0"}
	if !containsName(names, "MyHotspot") {
		t.Error("expected containsName to find an exact match")
	}
	if containsName(names, "MyHotspotExtra") {
		t.Error("expected containsName NOT to match a superstring (exact match only, unlike the old grep -w reference)")
	}
}
