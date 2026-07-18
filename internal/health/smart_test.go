// internal/health/smart_test.go
package health

import "testing"

const sampleSmartctlJSONPassed = `{
  "smart_status": {"passed": true},
  "device": {"name": "/dev/sda"}
}`

const sampleSmartctlJSONFailed = `{
  "smart_status": {"passed": false},
  "device": {"name": "/dev/sda"}
}`

func TestParseSmartctlJSON_Passed(t *testing.T) {
	healthy, ok := parseSmartctlJSON([]byte(sampleSmartctlJSONPassed))
	if !ok {
		t.Fatal("expected ok=true for valid JSON")
	}
	if !healthy {
		t.Error("expected healthy=true when smart_status.passed is true")
	}
}

func TestParseSmartctlJSON_Failed(t *testing.T) {
	healthy, ok := parseSmartctlJSON([]byte(sampleSmartctlJSONFailed))
	if !ok {
		t.Fatal("expected ok=true for valid JSON")
	}
	if healthy {
		t.Error("expected healthy=false when smart_status.passed is false")
	}
}

func TestParseSmartctlJSON_Invalid(t *testing.T) {
	if _, ok := parseSmartctlJSON([]byte("not json")); ok {
		t.Error("expected ok=false for invalid JSON")
	}
}

func TestParseSmartctlText_Passed(t *testing.T) {
	output := "=== START OF READ SMART DATA SECTION ===\nSMART overall-health self-assessment test result: PASSED\n"
	healthy, ok := parseSmartctlText(output)
	if !ok || !healthy {
		t.Errorf("expected ok=true healthy=true, got ok=%v healthy=%v", ok, healthy)
	}
}

func TestParseSmartctlText_Failed(t *testing.T) {
	output := "SMART overall-health self-assessment test result: FAILED!\n"
	healthy, ok := parseSmartctlText(output)
	if !ok || healthy {
		t.Errorf("expected ok=true healthy=false, got ok=%v healthy=%v", ok, healthy)
	}
}

func TestParseSmartctlText_NoMatch(t *testing.T) {
	if _, ok := parseSmartctlText("smartctl 7.3\nsome unrelated output\n"); ok {
		t.Error("expected ok=false when the health line isn't present")
	}
}

func TestIsNVMeNamespace(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/dev/nvme0n1", true},
		{"/dev/nvme1n1", true},
		{"/dev/nvme0n2", true},
		{"/dev/nvme10n1", true},
		{"/dev/nvme0n1p1", false},
		{"/dev/nvme0n1p12", false},
		{"/dev/nvme1n1p3", false},
		{"/dev/sda", false},
		{"/dev/nvme0", false},
	}
	for _, c := range cases {
		if got := isNVMeNamespace(c.path); got != c.want {
			t.Errorf("isNVMeNamespace(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}
