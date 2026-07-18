// internal/server/wol_test.go
package server

import (
	"bytes"
	"testing"
)

func TestBuildMagicPacket(t *testing.T) {
	packet, err := BuildMagicPacket("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(packet) != 102 {
		t.Fatalf("expected 102 bytes, got %d", len(packet))
	}

	header := packet[:6]
	if !bytes.Equal(header, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}) {
		t.Errorf("expected header of 6x 0xFF, got %x", header)
	}

	mac := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	for i := 0; i < 16; i++ {
		chunk := packet[6+i*6 : 6+i*6+6]
		if !bytes.Equal(chunk, mac) {
			t.Errorf("repetition %d: expected %x, got %x", i, mac, chunk)
		}
	}
}

func TestBuildMagicPacket_InvalidMAC(t *testing.T) {
	if _, err := BuildMagicPacket("not-a-mac"); err == nil {
		t.Fatal("expected error for invalid MAC, got nil")
	}
}
