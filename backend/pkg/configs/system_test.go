package configs

import (
	"testing"
	"os"
)

func Test_SystemConfig(t *testing.T) {
	systemConfigLocation := os.Getenv("ASSETS_LOCATION") + "/configs/system_config.yaml"
	t.Run("func LoadSystemConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadSystemConfigs(systemConfigLocation)
			if result == nil || len(result.AvailableSites) != 3 {
				t.Fatalf("result: %v", result)
			}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadSystemConfigs(systemConfigLocation + "abc")
			if result != nil {
				t.Fatalf("result: %v", result)
			}
		})
	})
}