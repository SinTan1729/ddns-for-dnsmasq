package internal

type ipInfo struct {
	IP   string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}

type httpError struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason,omitempty"`
}

type HostConfig struct {
	Host   string `yaml:"host"`
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
