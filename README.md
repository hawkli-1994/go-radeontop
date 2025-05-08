# go-radeontop

A Go library for monitoring AMD GPU statistics and performance metrics, specifically designed for Linux systems where traditional monitoring tools like `rocm-smi` are not available.

## Features

- Machine-friendly API for programmatic access to GPU metrics
- No dependency on ROCm tools
- Real-time monitoring of AMD GPU metrics including:
  - GPU Usage
  - Memory Usage
  - VRAM Usage
  - Temperature (Edge, Memory, Junction)
- Support for multiple GPU devices
- Lightweight and efficient implementation
- Uses native Linux sysfs and `sensors` for data collection

## System Requirements

- Linux operating system
- AMD GPU
- `lm-sensors` package installed
- Appropriate system permissions to access GPU information

## Installation

```bash
go get github.com/hawkli-1994/go-radeontop
```

## Usage

Basic monitoring example:

```go
package main

import (
    "fmt"
    "log/slog"
    "time"
    "github.com/hawkli-1994/go-radeontop/pkg/monitor"
)

func main() {
    // Create a new monitor instance with default logger
    mon, err := monitor.New(slog.Default())
    if err != nil {
        log.Fatalf("Failed to create monitor: %v", err)
    }

    // Get device information
    info, err := mon.GetDeviceInfoList()
    if err != nil {
        log.Fatalf("Failed to get device info: %v", err)
    }

    // Print basic device information
    fmt.Printf("Monitoring GPU: %s (Driver: %s)\n", 
        info.Items[0].Name, 
        info.Items[0].DriverVersion)

    // Monitor GPU stats every second
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for range ticker.C {
        stats, err := mon.GetDeviceInfoList()
        if err != nil {
            log.Printf("Error getting stats: %v", err)
            continue
        }

        gpuStats := stats.Items[0].Stats
        fmt.Printf("GPU: %.1f%% | MemoryUsage: %.1f%% | VRAM: %d/%d | "+
            "Temp: Edge:%.1f°C Mem:%.1f°C Junction:%.1f°C\n",
            gpuStats.GPUUsage,
            gpuStats.MemoryUsage,
            gpuStats.VRAMUsed,
            gpuStats.VRAMTotal,
            gpuStats.GpuTempEdge,
            gpuStats.GpuTempMem,
            gpuStats.GpuTempJunction)
    }
}
```

## Available Metrics

- `GPUUsage`: GPU utilization percentage (0-100)
- `MemoryUsage`: Memory utilization percentage (0-100)
- `VRAMUsed`: VRAM usage in bytes
- `VRAMTotal`: Total VRAM available in bytes
- `GpuTempEdge`: Edge temperature in Celsius
- `GpuTempMem`: Memory temperature in Celsius
- `GpuTempJunction`: Junction temperature in Celsius

## License

This project is licensed under the MIT License - see the LICENSE file for details.
