package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/hawkli-1994/go-radeontop/pkg/monitor"
)

func main() {
	// Create a new monitor instance
	mon, err := monitor.New(slog.Default())
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	// Get device information
	info, err := mon.GetDeviceInfoList()
	if err != nil {
		log.Fatalf("Failed to get device info: %v", err)
	}
	fmt.Printf("Monitoring GPU: %s (Driver: %s)\n", info.Items[0].Name, info.Items[0].DriverVersion)

	// Monitor GPU stats every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	fmt.Println("\nPress Ctrl+C to stop monitoring...")
	for range ticker.C {
		stats, err := mon.GetDeviceInfoList()
		if err != nil {
			log.Printf("Error getting stats: %v", err)
			continue
		}

		fmt.Printf("\033[2K\r") // Clear line
		fmt.Printf("GPU: %.1f%% | MemoryUsed: %.1f%% | MEMTotal: %.1f%% | TempEdge: %d°C | TempMem: %d°C | TempJunction : %d°C |",
			stats.Items[0].Stats.GPUUsage,
			stats.Items[0].Stats.MemoryUsage,
			stats.Items[0].Stats.VRAMTotal,
			stats.Items[0].Stats.GpuTempEdge,
			stats.Items[0].Stats.GpuTempMem,
			stats.Items[0].Stats.GpuTempJunction)
	}
}
