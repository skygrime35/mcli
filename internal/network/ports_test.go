package network

import "testing"

const sampleProcNetTCP = `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1538 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 21234 1 0000000000000000 100 0 0 10 0
   1: 00000000:01BB 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 21235 1 0000000000000000 100 0 0 10 0
   2: 0100007F:C350 0100007F:1F90 01 00000000:00000000 00:00000000 00000000  1000        0 21236 1 0000000000000000 20 4 30 10 -1
`

const sampleProcNetUDP = `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0044 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 21237 2 0000000000000000 0
`

func TestParseProcNet_TCPListeningOnly(t *testing.T) {
	got := parseProcNet(sampleProcNetTCP, "tcp", true)
	want := []ListeningPort{
		{Protocol: "tcp", Port: 5432}, // 0x1538
		{Protocol: "tcp", Port: 443},  // 0x01BB
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v (only st=0A rows must be included)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestParseProcNet_UDPAllRows(t *testing.T) {
	got := parseProcNet(sampleProcNetUDP, "udp", false)
	want := []ListeningPort{{Protocol: "udp", Port: 68}} // 0x0044
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestParseProcNet_EmptyInput(t *testing.T) {
	header := "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"
	if got := parseProcNet(header, "tcp", true); len(got) != 0 {
		t.Errorf("expected no entries for header-only input, got %v", got)
	}
}
