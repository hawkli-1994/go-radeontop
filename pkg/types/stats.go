package types

// GPUStats represents the statistics of an AMD GPU
type GPUStats struct {
	// GPU utilization percentage (0-100)
	GPUUsage float64

	// Memory utilization percentage (0-100)
	MemoryUsage float64

	// Current GPU clock speed in MHz
	ClockSpeed uint32

	// Current memory clock speed in MHz
	MemoryClockSpeed uint32

	// Temperature in Celsius
	Temperature int32

	// Fan speed percentage (0-100)
	FanSpeed uint8

	// Power consumption in watts
	PowerConsumption float64

	// VRAM usage in bytes
	VRAMUsed uint64

	// Total VRAM available in bytes
	VRAMTotal uint64
}

// DeviceInfo contains static information about the GPU
type DeviceInfo struct {
	// Device ID
	DeviceID string

	// Device name
	Name string

	// Driver version
	DriverVersion string

	// PCI bus information
	PCIBusID string
}
