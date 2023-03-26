package cli

import (
	"encoding/json"
	"fmt"
	"gosail/client"
	"gosail/cycle"
	"gosail/model"
	"gosail/utils"
	"strings"
)

var (
	workPath string
)

func interExec() {
	cmdLine = fmt.Sprintf(`cd %s;%s;echo -e "\n"$(pwd)`, workPath, cmdLine)
	clientConfig, err = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil {
		log.Error(err)
		return
	}

	spinnerConfig.IsSelect = selection

	sshResults := cycle.Exec(clientConfig, spinnerConfig)
	if sshResults[0].Success {
		if !linuxMode {
			workPath = getPWDPath(sshResults[0].Result, 4)
		} else {
			workPath = getPWDPath(sshResults[0].Result, 2)
		}
	}
	interShowExecResult(clientConfig.SshHosts, sshResults, &jsonMode, &linuxMode)

	if selection {
		client.LoginHostByID(clientConfig.SshHosts, sshResults, "")
	}
}

func interShowExecResult(sshHosts []model.SSHHost, sshResults []model.RunResult, jsonMode, linuxMode *bool) {
	if *jsonMode {
		jsonResult, err := json.Marshal(sshResults)
		if err != nil {
			log.Errorf("json Marshal error, %v", err)
		}
		fmt.Println(string(jsonResult))
		return
	}
	for id, sshResult := range sshResults {
		fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", sshResult.Username, sshResult.Host, id)
		if sshResult.Success {
			if !*linuxMode {
				sshResult.Result = utils.SimpleLine(sshResult.Result, len(sshHosts[id].CmdList)+3, 6)
			} else {
				sshResult.Result = splitLastLine(sshResult.Result)
			}
		}
		fmt.Println(sshResult.Result)
	}
}

func getPWDPath(result string, lastline int) string {
	resList := strings.Split(result, "\n")
	var pwdPath string
	if len(resList)-lastline >= 0 {
		pwdPath = resList[len(resList)-lastline]
		pwdPath = strings.TrimRight(pwdPath, "\r\n")
	} else {
		pwdPath = workPath
	}
	return pwdPath
}

func interK8sExec() {
	before := cmdLine
	if workPath != "~" {
		cmdLine = fmt.Sprintf(`cd %s;%s;echo -e "\n"$(pwd)`, workPath, cmdLine)
	} else {
		cmdLine = fmt.Sprintf(`%s;echo -e "\n"$(pwd)`, cmdLine)
	}

	clientConfig, err := cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, true)
	if err != nil {
		log.Error(err)
		return
	}
	kubeConfig := &model.KubeConfig{
		SshHosts:  clientConfig.SshHosts,
		Namespace: namespace,
		App:       app,
		Container: container,
		Label:     label,
		Shell:     shell,
		Highlight: highlight,
		CmdLine:   cmdLine,
	}
	spinnerConfig.IsSelect = selection
	sshResults := cycle.K8sExec(clientConfig, kubeConfig, spinnerConfig)
	if sshResults[0].Success {
		workPath = getPWDPath(sshResults[0].Result, 2)
	}
	kubeConfig.CmdLine = before

	interK8sShowResults(sshResults, kubeConfig, &jsonMode)
	if selection {
		client.LoginPodByID(kubeConfig, clientConfig.SshHosts, sshResults, kubeConfig.Shell)
	}
}

func interK8sShowResults(sshResults []model.RunResult, kubeConfig *model.KubeConfig, jsonMode *bool) {
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
		fmt.Printf("👇===============> %-15s (%s) <===============[%-3d]\n", sshResults[id].Host, kubeConfig.Container, id)
		if kubeConfig.CmdLine != "" {
			fmt.Printf("👉 ------------> %s \n", kubeConfig.CmdLine)
		}
		if sshResult.Success {
			sshResult.Result = splitLastLine(sshResult.Result)
		}
		fmt.Println(sshResult.Result)
	}
}

func splitLastLine(result string) string {
	resList := strings.Split(result, "\n")
	if len(resList) == 2 {
		return ""
	}
	return utils.SimpleLine(result, 0, 2)
}
