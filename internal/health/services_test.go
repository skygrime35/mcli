package health

import "testing"

func TestParseFailedUnits(t *testing.T) {
	output := "docker.service    loaded failed failed Docker Application Container Engine\nfoo.service       loaded failed failed Some other service\n"
	got := parseFailedUnits(output)
	want := []string{"docker.service", "foo.service"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseFailedUnits_Empty(t *testing.T) {
	if got := parseFailedUnits(""); len(got) != 0 {
		t.Errorf("expected no units for empty output, got %v", got)
	}
}

func TestParseUnitFileNames(t *testing.T) {
	output := "" +
		"proc-sys-fs-binfmt_misc.automount enabled        enabled\n" +
		"ssh.service                        enabled        enabled\n" +
		"cron.service                       enabled        enabled\n" +
		"systemd-resolved.service           enabled        enabled\n" +
		"bluetooth.service                  disabled       disabled\n" +
		"\n" +
		"5 unit files listed.\n"

	got := parseUnitFileNames(output)

	for _, want := range []string{"ssh", "cron", "systemd-resolved", "bluetooth"} {
		if !got[want] {
			t.Errorf("expected %q to be present in parsed unit set, got %v", want, got)
		}
	}
	if got["sshd"] {
		t.Error("did not expect \"sshd\" to be present - it wasn't in the fixture output")
	}
	if got["NetworkManager"] {
		t.Error("did not expect \"NetworkManager\" to be present - it wasn't in the fixture output")
	}
}
