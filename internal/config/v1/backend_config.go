package config

type BackendConfig struct {
	APIRoutePrefix  string `yaml:"api_route_prefix" json:"api_route_prefix"`
	LiteRoutePrefix string `yaml:"lite_route_prefix" json:"lite_route_prefix"`

	Database DatabaseConfig `yaml:"database" json:"database"`
	Trace    TraceConfig    `yaml:"trace" json:"trace"`

	StoragePath    string   `yaml:"storage_path" json:"storage_path"`
	AvailableSites []string `yaml:"available_sites" json:"available_sites"`
}
