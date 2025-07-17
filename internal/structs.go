package internal

type httpError struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason,omitempty"`
}

type HostConfig struct {
	APIKey string `yaml:"api-key"`
}

type Config struct {
	IPHeader string                `yaml:"ip-header"`
	Port     int                   `yaml:"port"`
	Hosts    map[string]HostConfig `yaml:"hosts"`
}

type hostEntry struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}
