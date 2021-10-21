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

// YAML Config Strongly Typed
type LLamaConfig struct {
	Summarization struct {
		Interval int `yaml:"interval"`
		Handlers int `yaml:"handlers"`
	} `yaml:"summarization"`
	API struct {
		Bind string `yaml:"bind"`
	} `yaml:"api"`
	Ports struct {
		Default struct {
			IP      string `yaml:"ip"`
			Port    int    `yaml:"port"`
			Tos     int    `yaml:"tos"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"default"`
	} `yaml:"ports"`
	PortGroups struct {
		Default []struct {
			Port  string `yaml:"port"`
			Count int    `yaml:"count"`
		} `yaml:"default"`
	} `yaml:"port_groups"`
	RateLimits struct {
		Default struct {
			Cps float64 `yaml:"cps"`
		} `yaml:"default"`
	} `yaml:"rate_limits"`
	Tests []struct {
		Targets   string `yaml:"targets"`
		PortGroup string `yaml:"port_group"`
		RateLimit string `yaml:"rate_limit"`
	} `yaml:"tests"`
	Targets struct {
		Default []struct {
			IP   string `yaml:"ip"`
			Port int    `yaml:"port"`
			Tags struct {
				Version        string `yaml:"version"`
				ProbeShortname string `yaml:"probe_shortname"`
				ProbeName      string `yaml:"probe_name"`
				DstName        string `yaml:"dst_name"`
				DstShortname   string `yaml:"dst_shortname"`
				SrcName        string `yaml:"src_name"`
				SrcShortname   string `yaml:"src_shortname"`
				Group          string `yaml:"group"`
			} `yaml:"tags"`
		} `yaml:"default"`
	} `yaml:"targets"`
}
