package main

import (
    "fmt"
    "time"
    "log"
    "github.com/hawkli-1994/go-radeontop/pkg/monitor"
)

func main() {
    // Create a new monitor instance
    mon, err := monitor.New()
    if err != nil {
        log.Fatalf("Failed to create monitor: %v", err)
    }
    defer mon.Close()

    // Get device information
    info, err := mon.GetDeviceInfo()
    if err != nil {
        log.Fatalf("Failed to get device info: %v", err)
    }
    fmt.Printf("Monitoring GPU: %s (Driver: %s)\n", info.Name, info.DriverVersion)

    // Monitor GPU stats every second
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    fmt.Println("\nPress Ctrl+C to stop monitoring...")
    for range ticker.C {
        stats, err := mon.GetStats()
        if err != nil {
            log.Printf("Error getting stats: %v", err)
            continue
        }

        fmt.Printf("\033[2K\r") // Clear line
        fmt.Printf("GPU: %.1f%% | Memory: %.1f%% | Temp: %d°C | Power: %.1fW",
            stats.GPUUsage,
            stats.MemoryUsage,
            stats.Temperature,
            stats.PowerConsumption)
    }
} 