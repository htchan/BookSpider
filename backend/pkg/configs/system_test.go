package configs

import (
	"testing"
	"os"
)

func Test_SystemConfig(t *testing.T) {
	systemConfigDirectory := os.Getenv("ASSETS_LOCATION") + "/configs"
	t.Run("func LoadSystemConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadSystemConfigs(systemConfigDirectory)
			if result == nil || len(result.AvailableSiteConfigs) != 3 {
				t.Fatalf("result: %v", result)
			}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadSystemConfigs(systemConfigDirectory + "abc")
			if result != nil {
				t.Fatalf("result: %v", result)
			}
		})
	})
}