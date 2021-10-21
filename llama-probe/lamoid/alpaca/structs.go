package alpaca

// Create struct for JSON we send to the server for registration
type PayLoad struct {
	Port      int    `json:"port"`
	Keepalive int    `json:"keepalive,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Tags      struct {
		Version        string `json:"version"`
		ProbeShortname string `json:"probe_shortname"`
		ProbeName      string `json:"probe_name"`
	} `json:"tags"`
	Group string `json:"group,omitempty"`
}

// LamoidEnv struct containing the running environment information
// for the grazzing llama probe.
type LamoidEnv struct {
	SourceIP       string `env:"LLAMA_SOURCE_IP"`
	Server         string `env:"LLAMA_SERVER"`
	Group          string `env:"LLAMA_GROUP"`
	Port           int    `env:"LLAMA_PORT"`
	KeepAlive      int    `env:"LLAMA_KEEPALIVE"`
	ProbeName      string `env:"PROBE_NAME"`
	ProbeShortName string `env:"PROBE_SHORTNAME"`
	ServerURL      string
	ReflectorPID   int
	CollectorPID   int
}
