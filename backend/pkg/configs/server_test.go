package configs

import (
	"testing"
	"os"
)

func Test_ServerConfig(t *testing.T) {
	serverConfigLocation := os.Getenv("ASSETS_LOCATION") + "/configs/server_config.yaml"
	t.Run("func LoadServerConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadServerConfigs(serverConfigLocation)
			if result == nil || len(result.AvailableApi) != 7 {
				t.Fatalf("result: %v", result)
			}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadServerConfigs(serverConfigLocation + "abc")
			if result != nil {
				t.Fatalf("result: %v", result)
			}
		})
	})
}