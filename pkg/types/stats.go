package types

import (
	// "os"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"log/slog"
)

var drmPath = "/sys/class/drm/"

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

	GpuTempEdge     float64
	GpuTempMem      float64
	GpuTempJunction float64

	// Fan speed percentage (0-100)
	FanSpeed uint8

	// Power consumption in watts
	PowerConsumption float64

	// VRAM usage in bytes
	VRAMUsed uint64

	// Total VRAM available in bytes
	VRAMTotal uint64
}

type DeviceInfoList struct {
	Items []DeviceInfo
}

// DeviceInfo contains static information about the GPU
type DeviceInfo struct {
	// Device ID
	DeviceID string

	// Device name (e.g. "card0")
	Name string

	// Driver version
	DriverVersion string

	// PCI bus information
	PCIBusID string

	// Sensor name
	SensorName string
	Stats      *GPUStats
}

// sensors is the data structure for the sensors.json file
type sensorsGpu struct {
	Adapter string `json:"Adapter"`
	Vddgfx  struct {
		In0_input float64 `json:"in0_input"`
	} `json:"vddgfx"`
	Fan1 struct {
		Fan1_input float64 `json:"fan1_input"`
	} `json:"fan1"`
	Edge struct {
		Temp1_input float64 `json:"temp1_input"`
	} `json:"edge"`
	Junction struct {
		Temp2_input float64 `json:"temp2_input"`
	} `json:"junction"`
	Mem struct {
		Temp3_input float64 `json:"temp3_input"`
	} `json:"mem"`
	PPT struct {
		Power1_average float64 `json:"power1_average"`
	} `json:"PPT"`
}

type sensorsOutput struct {
	Sensors map[string]sensorsGpu
}

// ParseSensorsFile parses the sensors.json file and returns only AMD GPU entries
func ParseSensorsFile(data []byte) (*sensorsOutput, error) {
	var rawOutput map[string]sensorsGpu
	if err := json.Unmarshal(data, &rawOutput); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensors data: %w", err)
	}

	// Create new map with only AMD GPU entries
	filteredOutput := make(map[string]sensorsGpu)
	for key, value := range rawOutput {
		if strings.HasPrefix(key, "amdgpu-") {
			filteredOutput[key] = value
		}
	}

	return &sensorsOutput{
		Sensors: filteredOutput,
	}, nil
}

func NewDeviceInfoList(logger *slog.Logger) (*DeviceInfoList, error) {

	sensorsCmd := exec.Command("sensors", "-j")
	sensorsOutput, err := sensorsCmd.Output()
	if err != nil {
		logger.Error("failed to get sensors output", "error", err)
		return nil, err
	}

	sensors, err := ParseSensorsFile(sensorsOutput)
	if err != nil {
		logger.Error("failed to parse sensors output", "error", err)
		return nil, err
	}

	logger.Info("sensors", "sensors", sensors)

	cmd := exec.Command("ls", drmPath)
	output, err := cmd.Output()
	if err != nil {
		logger.Error("failed to get device info list", "error", err)
		return nil, err
	}

	// Create regexp to match only cardX directories
	re := regexp.MustCompile(`^card[0-9]+$`)

	var devices []DeviceInfo
	// Split output into lines and process each line
	for _, line := range strings.Split(string(output), "\n") {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}

		// Only add if it matches cardX pattern
		if re.MatchString(name) {
			devices = append(devices, *NewDeviceInfo(name, logger, sensors))
		}
	}

	return &DeviceInfoList{
		Items: devices,
	}, nil
}

func NewDeviceInfo(name string, logger *slog.Logger, sensors *sensorsOutput) *DeviceInfo {

	path := filepath.Join(drmPath, name, "device")

	devicePath := filepath.Join(path, "device")
	deviceIDCmd := exec.Command("cat", devicePath)
	deviceID, err := deviceIDCmd.Output()
	if err != nil {
		logger.Error("failed to get device ID", "error", err)
		return nil
	}

	classCmd := exec.Command("cat", filepath.Join(path, "class"))
	class, err := classCmd.Output()
	if err != nil {
		logger.Error("failed to get device name", "error", err)
		return nil
	}
	// class like 0x030000, convert to amdgpu-pci-0300
	SensorName := convertClassToSensorName(string(class))
	sensor, exists := sensors.Sensors[SensorName]
	if !exists {
		logger.Error("failed to get sensor name", "error", err)
		return nil
	}

	gpuUsageCmd := exec.Command("cat", filepath.Join(path, "gpu_busy_percent"))
	gpuUsage, err := gpuUsageCmd.Output()
	if err != nil {
		logger.Error("failed to get gpu usage", "error", err)
		return nil
	}
	gpuUsageInt, err := strconv.Atoi(strings.TrimSpace(string(gpuUsage)))
	if err != nil {
		logger.Error("failed to convert gpu usage to int", "error", err)
		return nil
	}

	memoryUsageCmd := exec.Command("cat", filepath.Join(path, "mem_info_vram_used"))
	memoryUsage, err := memoryUsageCmd.Output()
	if err != nil {
		logger.Error("failed to get memory usage", "error", err)
		return nil
	}
	memoryUsageInt, err := strconv.Atoi(strings.TrimSpace(string(memoryUsage)))
	if err != nil {
		logger.Error("failed to convert memory usage to int", "error", err)
		return nil
	}
	// mem_info_vram_total
	memoryTotalCmd := exec.Command("cat", filepath.Join(path, "mem_info_vram_total"))
	memoryTotal, err := memoryTotalCmd.Output()
	if err != nil {
		logger.Error("failed to get memory total", "error", err)
		return nil
	}
	memoryTotalInt, err := strconv.Atoi(strings.TrimSpace(string(memoryTotal)))
	if err != nil {
		logger.Error("failed to convert memory total to int", "error", err)
		return nil
	}

	return &DeviceInfo{
		Name:       name,
		DeviceID:   string(deviceID),
		PCIBusID:   string(class),
		SensorName: SensorName,
		Stats: &GPUStats{
			GPUUsage:         float64(gpuUsageInt),
			MemoryUsage:      float64(memoryUsageInt),
			VRAMTotal:        uint64(memoryTotalInt),
			VRAMUsed:         uint64(memoryUsageInt),
			ClockSpeed:       0,
			MemoryClockSpeed: 0,
			GpuTempEdge:      sensor.Edge.Temp1_input,
			GpuTempMem:       sensor.Mem.Temp3_input,
			GpuTempJunction:  sensor.Junction.Temp2_input,
		},
	}
}

func convertClassToSensorName(class string) string {
	// class like 0x030000, convert to amdgpu-pci-0300
	classStr := strings.TrimSpace(class)
	// Remove "0x" prefix and last two zeros
	classStr = strings.TrimSuffix(strings.TrimPrefix(classStr, "0x"), "00")
	return "amdgpu-pci-" + classStr
}
