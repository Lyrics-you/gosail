package main

import (
	"flag"
	"fmt"
	"gosail/client"
	"gosail/goscp"
	ssh "gosail/gossh"
	"gosail/logger"
	"gosail/model"
	"gosail/utils"

	"os"
	"time"
)

var (
	log = logger.Logger()
)

func main() {
	version := flag.Bool("v", false, "show version")
	hosts := flag.String("hosts", "", "host address list")
	hostFile := flag.String("hostfile", "", "hostfile path")
	ips := flag.String("ips", "", "ip address list")
	ipFile := flag.String("ipfile", "", "ipfile path")

	// gocy's arguments
	// to := flag.String("to", "", "push to destination path")
	pull := flag.String("pull", "", "pull's source path")
	push := flag.String("push", "", "push's source path")
	path := flag.String("path", "", "pull or push's destination path")

	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	key := flag.String("k", "", "ssh private key")
	port := flag.Int("port", 22, "ssh port")
	ciphers := flag.String("ciphers", "", "ciphers")
	keyExchanges := flag.String("keyexchanges", "", "keyexchanges")
	config := flag.String("config", "", "config file Path")
	timeLimit := flag.Int("tl", 30, "max timeout")
	numLimit := flag.Int("nl", 20, "max execute number")

	selection := flag.Bool("s", false, "select host to login")
	// linuxmode is true

	flag.Parse()

	var cmdList, hostList, cipherList, keyExchangeList []string
	var err error

	sshHosts := []model.SSHHost{}
	scpResults := []model.RunResult{}
	var host_Struct model.SSHHost

	// scp need linuxMode avoid timeout
	linuxMode := true

	if *version {
		fmt.Println("Welcome  : " + model.EMOJI["gocy"])
		fmt.Println("ToolName : " + "gocy")
		fmt.Println("Version  : " + model.VERSION)
		fmt.Println("Email    : Leyuan.Jia@Outlook.com")
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

			if *password == "" && *key == "" {
				*key = ssh.DefaultKeyPath()
			}
			host_Struct.Key = *key
			host_Struct.LinuxMode = linuxMode
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
	log.Infof("gocy start.")

	clientConfig := &model.ClientConfig{
		SshHosts:        sshHosts,
		NumLimit:        *numLimit,
		TimeLimit:       *timeLimit,
		CipherList:      cipherList,
		KeyExchangeList: keyExchangeList,
	}

	if *pull != "" && *push != "" {
		log.Errorf("push and pull cannot be used at the same time")
		return
	}

	if err != nil {
		log.Errorf("SCP pull failed, err : %v", err)
		return
	}

	endTime := time.Now()
	log.Infof("gocy finished. Process time %s. Number of active ip is %d.", endTime.Sub(startTime), len(sshHosts))

	var srcList, destList []string

	scpConfig := model.SCPConfig{
		SshHosts:  clientConfig.SshHosts,
		TimeLimit: clientConfig.TimeLimit,
		NumLimit:  clientConfig.NumLimit,
	}

	if *push != "" {
		for i := 0; i < len(sshHosts); i++ {
			srcList = append(srcList, *push)
			destList = append(destList, *path)
		}
		scpConfig.SrcPath = srcList
		scpConfig.DestPath = destList

		scpConfig.Method = "PUSH"
		goscp.SecureCopyPushMakeDir(&scpConfig)
		// scpConfig'SshHosts is equal clientConfig'SshHosts
		sshResults, _ := client.LimitShhWithGroup(clientConfig)

		scpResults, _ = client.LimitScpWithGroup(&scpConfig, sshResults)
		for id, scpResult := range scpResults {
			fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", sshHosts[id].Username, "localhost", id)
			if sshResults[id].Success {
				fmt.Print(sshResults[id].Result)
			}
			fmt.Print(scpResult.Result)
			fmt.Println()
		}
		fmt.Println("👌Finshed!")
	}

	if *pull != "" {
		for i := 0; i < len(sshHosts); i++ {
			srcList = append(srcList, *pull)
			destList = append(destList, *path)
		}
		scpConfig.SrcPath = srcList
		scpConfig.DestPath = destList

		scpConfig.Method = "PULL"
		mkdirResults := goscp.SecureCopyPullMakeDir(&scpConfig)

		scpResults, _ = client.LimitScpWithGroup(&scpConfig, mkdirResults)

		for id, scpResult := range scpResults {
			fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", sshHosts[id].Username, "localhost", id)
			if mkdirResults[id].Success {
				fmt.Print(mkdirResults[id].Result)
			}
			fmt.Print(scpResult.Result)
			fmt.Println()
		}
		fmt.Println("👌Finshed!")
	}

	if err != nil {
		log.Errorf("SCP push failed, err : %v", err)
		return
	}

	if *selection {
		client.LoginHostByID(sshHosts, scpResults, "")
	}

}
