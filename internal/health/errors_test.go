package health

import "testing"

func TestParseCount(t *testing.T) {
	if got := parseCount("line one\nline two\nline three\n"); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestParseCount_Empty(t *testing.T) {
	if got := parseCount(""); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestParseCount_TrailingBlankLine(t *testing.T) {
	// journalctl/dmesg output with a trailing newline shouldn't count an
	// extra empty line.
	if got := parseCount("line one\nline two\n"); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}
