// internal/network/speedtest.go
package network

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

const fallbackTestURL = "http://speedtest.tele2.net/1MB.zip"
const fallbackTestSizeBytes = 1 << 20 // 1MB.zip is exactly 1 MiB

// RunSpeedTest measures download/upload throughput and latency against
// the nearest speedtest.net server. If that fails for any reason (no
// server reachable, network error), it falls back to timing a plain HTTP
// download of a small public test file - matching the old Python
// reference's two-tier behavior (speedtest library, then a manual HTTP
// download fallback).
func RunSpeedTest() (SpeedTestResult, error) {
	result, err := runLibrarySpeedTest()
	if err == nil {
		return result, nil
	}
	return fallbackSpeedTest()
}

func runLibrarySpeedTest() (SpeedTestResult, error) {
	client := speedtest.New()
	serverList, err := client.FetchServers()
	if err != nil {
		return SpeedTestResult{}, err
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil || len(targets) == 0 {
		return SpeedTestResult{}, fmt.Errorf("no speedtest server found")
	}

	server := targets[0]
	if err := server.PingTest(nil); err != nil {
		return SpeedTestResult{}, err
	}
	if err := server.DownloadTest(); err != nil {
		return SpeedTestResult{}, err
	}
	if err := server.UploadTest(); err != nil {
		return SpeedTestResult{}, err
	}

	return SpeedTestResult{
		DownloadMbps: server.DLSpeed.Mbps(),
		UploadMbps:   server.ULSpeed.Mbps(),
		PingMs:       float64(server.Latency.Microseconds()) / 1000.0,
		Server:       server.Name,
		Country:      server.Country,
	}, nil
}

func fallbackSpeedTest() (SpeedTestResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	start := time.Now()
	resp, err := client.Get(fallbackTestURL)
	if err != nil {
		return SpeedTestResult{}, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return SpeedTestResult{}, err
	}
	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		return SpeedTestResult{}, fmt.Errorf("download completed too fast to measure")
	}

	bitsPerSecond := float64(n) * 8 / elapsed
	return SpeedTestResult{
		DownloadMbps: bitsPerSecond / 1_000_000,
		UploadMbps:   0,
		PingMs:       0,
		Server:       "Fallback (HTTP Download)",
		Country:      "Unknown",
	}, nil
}
