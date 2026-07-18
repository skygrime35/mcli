// internal/server/wol.go
package server

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func BuildMagicPacket(mac string) ([]byte, error) {
	macBytes, err := parseMAC(mac)
	if err != nil {
		return nil, err
	}
	packet := make([]byte, 0, 102)
	for i := 0; i < 6; i++ {
		packet = append(packet, 0xFF)
	}
	for i := 0; i < 16; i++ {
		packet = append(packet, macBytes...)
	}
	return packet, nil
}

func parseMAC(mac string) ([]byte, error) {
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid MAC address: %q", mac)
	}
	out := make([]byte, 6)
	for i, p := range parts {
		b, err := hex.DecodeString(p)
		if err != nil || len(b) != 1 {
			return nil, fmt.Errorf("invalid MAC address: %q", mac)
		}
		out[i] = b[0]
	}
	return out, nil
}

func SendMagicPacket(host string, port int, mac string) error {
	packet, err := BuildMagicPacket(mac)
	if err != nil {
		return err
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("dialing %s: %w", addr, err)
	}
	defer conn.Close()
	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("sending magic packet: %w", err)
	}
	return nil
}
