// internal/health/system_test.go
package health

import "testing"

func TestParseOSRelease(t *testing.T) {
	content := "NAME=\"Ubuntu\"\nVERSION=\"22.04\"\nID=ubuntu\nPRETTY_NAME=\"Ubuntu 22.04.3 LTS\"\n"
	got := parseOSRelease(content)
	want := "Ubuntu 22.04.3 LTS"
	if got != want {
		t.Errorf("parseOSRelease() = %q, want %q", got, want)
	}
}

func TestParseOSRelease_Missing(t *testing.T) {
	if got := parseOSRelease("NAME=\"Ubuntu\"\n"); got != "Unknown" {
		t.Errorf("expected \"Unknown\" when PRETTY_NAME is absent, got %q", got)
	}
}
