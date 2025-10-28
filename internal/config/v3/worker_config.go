package config

type WorkerConfig struct {
	WorkerCount    int      `yaml:"worker_count" json:"worker_count"`
	AvailableSites []string `yaml:"available_sites" json:"available_sites"`

	DatabaseConfig DatabaseConfig          `yaml:"database" json:"database"`
	TraceConfig    TraceConfig             `yaml:"trace" json:"trace"`
	Client         map[string]ClientConfig `yaml:"client" json:"client"`

	StoragePath string `yaml:"storage_path" json:"storage_path"`
}
