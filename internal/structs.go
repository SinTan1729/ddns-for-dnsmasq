package internal

type ipInfo struct {
	IP   string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}

type httpError struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason,omitempty"`
}

type HostEntry struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"api-key"`
}

type AppData struct {
	IPHeader string      `yaml:"ip-header"`
	Port     int         `yaml:"port"`
	Hosts    []HostEntry `yaml:"hosts"`
}

type updatePayload struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}
