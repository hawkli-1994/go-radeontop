package monitor

import (
	"github.com/hawkli-1994/go-radeontop/internal/hardware"
	"github.com/hawkli-1994/go-radeontop/pkg/types"
)

// Monitor represents an AMD GPU monitoring instance
type Monitor struct {
	device *hardware.Device
}

// New creates a new Monitor instance
func New() (*Monitor, error) {
	dev, err := hardware.NewDevice()
	if err != nil {
		return nil, err
	}

	return &Monitor{
		device: dev,
	}, nil
}

// GetStats returns current GPU statistics
func (m *Monitor) GetStats() (*types.GPUStats, error) {
	return m.device.GetStats()
}

// GetDeviceInfo returns static information about the GPU
func (m *Monitor) GetDeviceInfo() (*types.DeviceInfo, error) {
	return m.device.GetDeviceInfo()
}

// Close releases any resources used by the monitor
func (m *Monitor) Close() error {
	return m.device.Close()
}
