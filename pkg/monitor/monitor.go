package monitor

import (
	"log/slog"

	"github.com/hawkli-1994/go-radeontop/internal/hardware"
	"github.com/hawkli-1994/go-radeontop/pkg/types"
)

// Monitor represents an AMD GPU monitoring instance
type Monitor struct {
	device *hardware.Device
	logger *slog.Logger
}

// New creates a new Monitor instance
func New(logger *slog.Logger) (*Monitor, error) {
	dev, err := hardware.NewDevice()
	if err != nil {
		return nil, err
	}

	return &Monitor{
		device: dev,
	}, nil
}

func (m *Monitor) GetDeviceInfoList() (*types.DeviceInfoList, error) {
	return hardware.GetDeviceInfoList(m.logger)
}
