package model

type SSHHost struct {
	Host      string    `json:"Host"`
	Port      int       `json:"Port"`
	Username  string    `json:"Username"`
	Password  string    `json:"Password"`
	CmdFile   string    `json:"CmdFile"`
	CmdLine   string    `json:"CmdLine"`
	CmdList   []string  `json:"CmdList"`
	Key       string    `json:"Key"`
	LinuxMode bool      `json:"LinuxMode"`
	Result    RunResult `json:"-"`
}

type RunResult struct {
	Host     string
	Username string
	Success  bool
	Result   string
}

type HostJson struct {
	Hosts  []SSHHost
	Global GlobalConfig
}

type GlobalConfig struct {
	Ciphers      string
	KeyExchanges string
}

type ClientConfig struct {
	SshHosts        []SSHHost
	NumLimit        int
	TimeLimit       int
	CipherList      []string
	KeyExchangeList []string
}

type SCPConfig struct {
	SshHosts  []SSHHost
	NumLimit  int
	TimeLimit int
	SrcPath   []string
	DestPath  []string
	Method    string
}

type KubeConfig struct {
	SshHosts  []SSHHost
	PodsList  []string
	Namespace string
	AppName   string
	Container string
	Label     string
}

type KubePods struct {
	MasterHost string
	PodsName   []string
	AppName    string
	Label      string
	Namespace  string
	Container  string
}
