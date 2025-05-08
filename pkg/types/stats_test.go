package types

import (
	_ "embed"
	"testing"
)

func TestConvertClassToSensorName(t *testing.T) {
	cases := []struct {
		class    string
		expected string
	}{
		{"0x030000", "amdgpu-pci-0300"},
		{"0x040000", "amdgpu-pci-0400"},
		{"0x050000", "amdgpu-pci-0500"},
		{"0x060000", "amdgpu-pci-0600"},
		{"0x070000", "amdgpu-pci-0700"},
		{"0x080000", "amdgpu-pci-0800"},
		{"0x090000", "amdgpu-pci-0900"},
	}
	for _, c := range cases {
		result := convertClassToSensorName(c.class)
		if result != c.expected {
			t.Errorf("convertClassToSensorName(%s) = %s; expected %s", c.class, result, c.expected)
		}
	}
}

//go:embed testdata/sensors.json
var sensorsJson []byte

func TestParseSensorsFile(t *testing.T) {
	if len(sensorsJson) == 0 {
		t.Fatal("Failed to load test data: sensors.json is empty")
	}

	sensors, err := ParseSensorsFile(sensorsJson)
	if err != nil {
		t.Fatalf("ParseSensorsFile() failed: %v", err)
	}

	if sensors == nil {
		t.Fatal("ParseSensorsFile() returned nil sensors")
	}

	if len(sensors.Sensors) != 1 {
		t.Errorf("ParseSensorsFile() got %v sensors, want 1", len(sensors.Sensors))
	}

	amdGpu, exists := sensors.Sensors["amdgpu-pci-0300"]
	if !exists {
		t.Fatal("Expected to find amdgpu-pci-0300 in sensors")
	}

	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{"Adapter name", float64(len(amdGpu.Adapter)), float64(len("PCI adapter"))},
		{"vddgfx voltage", amdGpu.Vddgfx.In0_input, 0.702},
		{"fan speed", amdGpu.Fan1.Fan1_input, 0.000},
		{"edge temperature", amdGpu.Edge.Temp1_input, 33.000},
		{"junction temperature", amdGpu.Junction.Temp2_input, 38.000},
		{"memory temperature", amdGpu.Mem.Temp3_input, 42.000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
