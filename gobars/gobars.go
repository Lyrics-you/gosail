package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gosail/client"
	"gosail/gokube"
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
	cmdLine := flag.String("cmdline", "", "command line")
	// not support cmdfile

	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	key := flag.String("k", "", "ssh private key")
	port := flag.Int("port", 22, "ssh port")
	ciphers := flag.String("ciphers", "", "ciphers")
	keyExchanges := flag.String("keyexchanges", "", "keyexchanges")
	config := flag.String("config", "", "config file Path")
	timeLimit := flag.Int("tl", 30, "max timeout")
	numLimit := flag.Int("nl", 20, "max execute number")

	// gobars' arguments
	container := flag.String("c", "", "k8s container")
	namespace := flag.String("n", "", "k8s namespace")
	appName := flag.String("app", "", "k8s app name")
	label := flag.String("l", "", "k8s label")

	// gobars' scp arguments
	copy := flag.Bool("copy", false, "k8s cp function")
	pull := flag.String("pull", "", "pull's source path")
	path := flag.String("path", "", "pull's destination path")

	selection := flag.Bool("s", false, "select host to login")
	// linuxmode is true

	jsonMode := flag.Bool("j", false, "print output in json format")
	outTxt := flag.Bool("otxt", false, "write result into txt")
	filePath := flag.String("fpath", "", "write file path")

	flag.Parse()

	var cmdList, hostList, cipherList, keyExchangeList []string
	var err error

	sshHosts := []model.SSHHost{}
	var host_Struct model.SSHHost
	linuxMode := true

	if *version {
		fmt.Println("Welcome  : " + model.EMOJI["gobars"])
		fmt.Println("ToolName : " + "gobars")
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
		fmt.Println(hostList)
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

	// if *cmdLine != "" {
	// 	cmdList = utils.SplitString(*cmdLine)
	// }

	if *ciphers != "" {
		cipherList = utils.SplitString(*ciphers)
	}

	if *keyExchanges != "" {
		keyExchangeList = utils.SplitString(*keyExchanges)
	}

	if *namespace == "" {
		log.Errorf("namespace is not specified")
		return
	}

	if *container == "" {
		*container = *appName
	}

	if *appName == "" {
		*appName = *container
	}

	if *container == "" && *appName == "" && *label == "" {
		log.Errorf("container or app name is not specified")
		return
	}

	if *config == "" {
		if len(hostList) == 0 {
			log.Errorf("hosts is empty")
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
			// if len(cmdList) == 0 {
			// 	log.Errorf("command is empty")
			// 	return
			// }

			host_Struct.CmdList = cmdList
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
			sshHosts[i].CmdList = utils.SplitString(sshHosts[i].CmdLine)
		}
	}

	// ssh connect
	startTime := time.Now()
	log.Infof("gobars start.")

	clientConfig := &model.ClientConfig{
		SshHosts:        sshHosts,
		NumLimit:        *numLimit,
		TimeLimit:       *timeLimit,
		CipherList:      cipherList,
		KeyExchangeList: keyExchangeList,
	}

	kubeConfig := &model.KubeConfig{
		SshHosts:  sshHosts,
		Namespace: *namespace,
		AppName:   *appName,
		Container: *container,
		Label:     *label,
	}

	gokube.KubeGetPods(kubeConfig)

	sshResults, _ := client.LimitShhWithGroup(clientConfig)
	// sshResults, _ := client.LimitShhWithChan(clientConfig)
	kubePods := gokube.GetPodsByResult(kubeConfig, sshResults)

	if *cmdLine != "" {
		kubeHosts := gokube.MakeMultiExecSshHosts(kubePods, sshHosts, *cmdLine)
		clientConfig.SshHosts = kubeHosts
		sshResults, _ = client.LimitShhWithGroup(clientConfig)
	}

	endTime := time.Now()
	log.Infof("gobars finished. Process time %s. Number of active ip is %d.", endTime.Sub(startTime), len(sshHosts))

	if *copy {
		if *pull == "" {
			log.Errorf("gobars pull's is empty")
			return
		} else {
			if *path == "" {
				*path = "./"
			}
			kubeHosts := gokube.MakeMultiCopySshHosts(kubePods, sshHosts, *pull, "./")
			clientConfig.SshHosts = kubeHosts
			client.LimitShhWithGroup(clientConfig)

			scpConfig := model.SCPConfig{
				SshHosts:  kubeHosts,
				TimeLimit: clientConfig.TimeLimit,
				NumLimit:  clientConfig.NumLimit,
				Method:    "PULL",
			}

			var srcList, destList []string
			for i := 0; i < len(kubeHosts); i++ {
				srcList = append(srcList, kubeConfig.PodsList[i])
				destList = append(destList, *path)
			}
			scpConfig.SrcPath = srcList
			scpConfig.DestPath = destList

			mkdirResults := goscp.SecureCopyPullMakeDir(&scpConfig)
			sshResults, _ = client.LimitScpWithGroup(&scpConfig, mkdirResults)

			delHosts := gokube.MakeMultiDeleteSshHosts(kubePods, sshHosts, "./")
			clientConfig.SshHosts = delHosts
			client.LimitShhWithGroup(clientConfig)
		}
	}

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

	for id, sshResult := range sshResults {
		sshResults[id].Host = kubeConfig.PodsList[id]
		fmt.Printf("👇===============> %-15s (%s) <===============[%-3d]\n", sshResults[id].Host, *container, id)
		if *cmdLine != "" {
			fmt.Printf("👉 ------------> %s \n", *cmdLine)
		}
		fmt.Println(sshResult.Result)

	}
	fmt.Println("👌Finshed!")

	if *selection {
		client.LoginPodByID(kubeConfig, sshHosts, sshResults, "/bin/bash")
	}
}
