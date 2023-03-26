package model

type GlobalConfig struct {
	NumLimit     int
	TimeLimit    int
	Ciphers      string
	KeyExchanges string
}

type ClientConfig struct {
	SshHosts        []SSHHost
	NumLimit        int `ini:"number_limit"`
	TimeLimit       int `ini:"timeout_limit"`
	CipherList      []string
	KeyExchangeList []string
}

type SpinConfig struct {
	SpinType int    `ini:"spin_type"`
	SpinTips string `ini:"spin_tips"`
	TimeOut  int
	IsSelect bool
}

type LoginConfig struct {
	HostFile string `ini:"host_file"`
	IpFile   string `ini:"ip_file"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Port     int    `ini:"port"`
}

type K8sConfig struct {
	Namespace string `ini:"namespace"`
	AppName   string `ini:"app"`
	Container string `ini:"container"`
	Label     string `ini:"label"`
	Shell     string `ini:"shell"`
}

type ModeConfig struct {
	JsonMode  bool `ini:"json_mode"`
	LinuxMode bool `ini:"linux_mode"`
	Selection bool `ini:"selection"`
}

type GosailConfig struct {
	Client *ClientConfig
	Spin   *SpinConfig
	Login  *LoginConfig
	K8s    *K8sConfig
	Mode   *ModeConfig
}
