# go-radeontop

A Go library for monitoring AMD GPU statistics and performance metrics, specifically designed for RISC-V architecture systems where traditional monitoring tools like `rocm-smi` are not available.

## Background

On RISC-V architecture systems, monitoring AMD GPUs can be challenging due to:
- Limited support for traditional tools like `rocm-smi`
- The original `radeontop` tool's output is not machine-friendly
- Lack of proper GPU monitoring solutions for RISC-V platforms

This library aims to solve these issues by providing a clean, Go-friendly API for GPU monitoring on RISC-V systems.

## Features

- Specifically designed for RISC-V architecture systems
- Machine-friendly API for programmatic access to GPU metrics
- No dependency on ROCm tools
- Real-time monitoring of AMD GPU metrics
- Easy-to-use API for integration
- Support for multiple GPU devices
- Lightweight and efficient implementation

## Use Cases

This library is particularly useful for:
- RISC-V systems running AMD GPUs
- Applications requiring programmatic access to GPU metrics
- Monitoring solutions that need machine-parseable GPU data
- Systems where ROCm tools are not available or cannot be installed
- DevOps and monitoring tools on RISC-V platforms

## Installation

```bash
go get github.com/hawkli-1994/go-radeontop
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/hawkli-1994/go-radeontop/pkg/monitor"
)

func main() {
    monitor, err := monitor.New()
    if err != nil {
        panic(err)
    }
    
    stats, err := monitor.GetStats()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("GPU Usage: %.2f%%\n", stats.GPUUsage)
}
```

## Documentation

For detailed documentation, please check the [docs](./docs) directory.

## System Requirements

- RISC-V architecture system
- AMD GPU
- Linux operating system
- No ROCm installation required

## Limitations

- This tool is specifically designed for RISC-V systems and may not be necessary on platforms where ROCm tools are available
- Some advanced features available in ROCm might not be accessible
- Requires appropriate system permissions to access GPU information

## License

This project is licensed under the MIT License - see the LICENSE file for details.
