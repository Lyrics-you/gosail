package cycle

import (
	"fmt"
	"gosail/gossh"
	"gosail/logger"
	"gosail/model"
	"gosail/utils"

	"gopkg.in/ini.v1"
)

var (
	LoginHost string
	LoginUser string
	LoginPwd  string
)
var (
	log = logger.Logger()
)

func GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key string, port, numLimit, timeLimit int, linuxMode bool) (*model.ClientConfig, error) {
	clientConfig := &model.ClientConfig{}
	clientConfig.NumLimit = numLimit
	clientConfig.TimeLimit = timeLimit
	clientConfig.KeyExchangeList = utils.SplitString(keyExchanges)
	clientConfig.CipherList = utils.SplitString(ciphers)
	hostList, err := utils.GetHostList(&hostLine, &hostFile, &ipLine, &ipFile)
	if err != nil {
		return clientConfig, err
	}
	var cmdList []string
	if linuxMode {
		cmdList = []string{cmdLine}
	} else {
		cmdList, err = utils.GetCmdList(&cmdLine, &cmdFile)
	}
	clientConfig.SshHosts = MakeSshHosts(hostList, cmdList, &username, &password, &key, port, &linuxMode)
	if err != nil {
		return clientConfig, err
	}
	return clientConfig, nil
}

func MakeSshHosts(hostList, cmdList []string, username, password, key *string, port int, linuxMode *bool) []model.SSHHost {
	sshHosts := []model.SSHHost{}
	var host_Struct model.SSHHost
	if len(hostList) == 0 {
		log.Errorf("hosts is empty")
		return []model.SSHHost{}
	}
	for _, host := range hostList {
		// user@host
		usr, hst := utils.SplitUserHost(host)
		if usr != "" {
			host_Struct.Host = hst
			host_Struct.Username = usr
		} else {
			host_Struct.Host = host
			host_Struct.Username = *username
		}

		host_Struct.Password = *password
		host_Struct.Port = port

		host_Struct.CmdList = cmdList
		if *password == "" && *key == "" {
			*key = gossh.DefaultKeyPath()
		}
		host_Struct.Key = *key
		host_Struct.LinuxMode = *linuxMode
		sshHosts = append(sshHosts, host_Struct)
	}
	return sshHosts
}

func MakeClientConfigFromJson(config string) (*model.ClientConfig, error) {
	sshHostConfig, err := utils.GetJson(config)
	if err != nil {
		log.Errorf("load config error, %v", err)
		return &model.ClientConfig{}, err
	}
	sshHosts := sshHostConfig.Hosts

	for i := 0; i < len(sshHosts); i++ {
		if sshHosts[i].CmdLine != "" {
			sshHosts[i].CmdList = utils.SplitString(sshHosts[i].CmdLine)
		} else {
			cmdList, err := utils.GetString(sshHosts[i].CmdFile)
			if err != nil {
				log.Errorf("load cmdFile error, %v", err)
				return &model.ClientConfig{}, err
			}
			sshHosts[i].CmdList = cmdList
		}
	}

	cipherList := utils.SplitString(sshHostConfig.Global.Ciphers)
	keyExchangeList := utils.SplitString(sshHostConfig.Global.KeyExchanges)

	return &model.ClientConfig{
		SshHosts:        sshHosts,
		CipherList:      cipherList,
		KeyExchangeList: keyExchangeList,
		NumLimit:        sshHostConfig.Global.NumLimit,
		TimeLimit:       sshHostConfig.Global.TimeLimit,
	}, nil
}

func GetIniConfig() *ini.File {
	config, err := ini.Load("./gosail.conf")
	if err != nil {
		log.Info("can't load gosial configuration from gosail.conf, err:", err)
		return nil
	}
	return config
}

func getSectionConfig(config *ini.File, section string, v interface{}) interface{} {
	if _, err := config.SectionsByName(section); err == nil {
		config.Section(section).MapTo(v)
	}
	return v
}

func GetGosailConfiguration() *model.GosailConfig {
	// get gosail configuration from ini file
	config := GetIniConfig()
	if config == nil {
		return nil
	}
	gosailConfig := &model.GosailConfig{
		Client: getSectionConfig(config, "client", &model.ClientConfig{}).(*model.ClientConfig),
		Spin:   getSectionConfig(config, "spin", &model.SpinConfig{}).(*model.SpinConfig),
		Login:  getSectionConfig(config, "login", &model.LoginConfig{}).(*model.LoginConfig),
		K8s:    getSectionConfig(config, "k8s", &model.K8sConfig{}).(*model.K8sConfig),
		Mode:   getSectionConfig(config, "mode", &model.ModeConfig{}).(*model.ModeConfig),
	}
	return gosailConfig
}

func PrintGosailConfiguration() {
	gosail := GetGosailConfiguration()
	fmt.Printf("%+v\n", gosail.Client)
	fmt.Printf("%+v\n", gosail.Spin)
	fmt.Printf("%+v\n", gosail.Login)
	fmt.Printf("%+v\n", gosail.K8s)
	fmt.Printf("%+v\n", gosail.Mode)
}
