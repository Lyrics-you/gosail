package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gosail/client"
	"gosail/logger"
	"gosail/model"
	"gosail/ssh"
	"gosail/utils"
	"os"
	"sync"
	"time"
)

var (
	log = logger.Logger()
	wg  sync.WaitGroup
)

func main() {
	version := flag.Bool("v", false, "show version")
	hosts := flag.String("hosts", "", "host address list")
	ips := flag.String("ips", "", "ip address list")
	cmdline := flag.String("cmdline", "", "command line")
	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	key := flag.String("k", "", "ssh private key")
	port := flag.Int("port", 22, "ssh port")
	ciphers := flag.String("ciphers", "", "ciphers")
	keyExchanges := flag.String("keyexchanges", "", "keyexchanges")
	cmdFile := flag.String("cmdfile", "", "cmdfile path")
	hostFile := flag.String("hostfile", "", "hostfile path")
	ipFile := flag.String("ipfile", "", "ipfile path")
	config := flag.String("c", "", "config file Path")
	jsonMode := flag.Bool("j", false, "print output in json format")
	outTxt := flag.Bool("otxt", false, "write result into txt")
	filePath := flag.String("path", "", "write file path")
	linuxMode := flag.Bool("l", false, "linux mode : multi command combine with && ,such as date&&cd /opt&&ls")
	timeLimit := flag.Int("t", 30, "max timeout")
	numLimit := flag.Int("n", 20, "max execute number")
	selection := flag.Bool("s", false, "select host to login")

	flag.Parse()

	var cmdList, hostList, cipherList, keyExchangeList []string
	var err error

	sshHosts := []model.SSHHost{}
	var host_Struct model.SSHHost

	if *version {
		fmt.Println("Version : " + model.VERSION)
		fmt.Println("Email : Leyuan.Jia@Outlook.com")
		os.Exit(0)
	}

	if *ipFile != "" {
		hostList, err = utils.GetIpListFromFile(*ipFile)
		if err != nil {
			log.Errorf("load iplist error, %v", err)
			return
		}
	}

	if *hostFile != "" {
		hostList, err = utils.GetString(*hostFile)
		if err != nil {
			log.Errorf("load hostfile error, %v", err)
			return
		}
	}

	if *ips != "" {
		hostList, err = utils.GetIpListFromString(*ips)
		if err != nil {
			log.Errorf("load iplist error, %v", err)
			return
		}
	}

	if *hosts != "" {
		hostList = utils.SplitString(*hosts)
	}

	if *cmdFile != "" {
		cmdList, err = utils.GetString(*cmdFile)
		if err != nil {
			log.Errorf("load cmdfile error, %v", err)
			return
		}
	}

	if *cmdline != "" {
		cmdList = utils.SplitString(*cmdline)
	}

	if *ciphers != "" {
		cipherList = utils.SplitString(*ciphers)
	}

	if *keyExchanges != "" {
		keyExchangeList = utils.SplitString(*keyExchanges)
	}

	if *config == "" {
		if len(hostList) == 0 {
			log.Warnf("hosts is empty")
			return
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
			host_Struct.Port = *port

			// empty command
			if len(cmdList) == 0 {
				log.Warnf("cmdline is empty")
				return
			}

			host_Struct.CmdList = cmdList
			if *password == "" && *key == "" {
				*key = ssh.DefaultKeyPath()
			}
			host_Struct.Key = *key
			host_Struct.LinuxMode = *linuxMode
			sshHosts = append(sshHosts, host_Struct)
		}
	} else {
		sshHostConfig, err := utils.GetJson(*config)
		if err != nil {
			log.Errorf("load config error, %v", err)
			return
		}
		cipherList = utils.SplitString(sshHostConfig.Global.Ciphers)
		keyExchangeList = utils.SplitString(sshHostConfig.Global.KeyExchanges)
		sshHosts = sshHostConfig.Hosts
		for i := 0; i < len(sshHosts); i++ {
			if sshHosts[i].CmdLine != "" {
				sshHosts[i].CmdList = utils.SplitString(sshHosts[i].CmdLine)
			} else {
				cmdList, err = utils.GetString(sshHosts[i].CmdFile)
				if err != nil {
					log.Errorf("load cmdFile error, %v", err)
					return
				}
				sshHosts[i].CmdList = cmdList
			}
		}
	}

	// ssh connect
	startTime := time.Now()
	log.Infof("gosail start.")

	clientConfig := &model.ClientConfig{
		SshHosts:        sshHosts,
		NumLimit:        *numLimit,
		TimeLimit:       *timeLimit,
		CipherList:      cipherList,
		KeyExchangeList: keyExchangeList,
	}

	sshResults, _ := client.LimitShhWithChan(clientConfig)

	endTime := time.Now()
	log.Infof("gosail finished. Process time %s. Number of active ip is %d", endTime.Sub(startTime), len(sshHosts))

	if *outTxt {
		for _, sshResult := range sshResults {
			err = utils.WriteIntoTxt(sshResult, *filePath)
			if err != nil {
				log.Errorf("write into txt error, %v", err)
				return
			}
		}
		return
	}

	if *jsonMode {
		jsonResult, err := json.Marshal(sshResults)
		if err != nil {
			log.Errorf("json Marshal error, %v", err)
		}
		fmt.Println(string(jsonResult))
		return
	}

	for _, sshResult := range sshResults {
		fmt.Println("> host: ", sshResult.Host)
		fmt.Println("=============== Result ===============")
		fmt.Println(sshResult.Result)
	}
	fmt.Println("############### Finshed ###############")

	if *selection {
		client.LoginHostByID(sshHosts)
	}
}
