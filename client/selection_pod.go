package client

import (
	"fmt"
	"gosail/gokube"
	"gosail/model"
)

func LoginPodByID(kubeConfig *model.KubeConfig, sshHosts []model.SSHHost, sshResult []model.RunResult, cmdLine string) {
	var id int
	if len(sshHosts) == 0 {
		id = -1
	} else {
		id = 0
	}

	// can be cyclically selected
	for id >= 0 {
		showPodList(sshResult)
		id := selectHost()
		if id == -1 {
			break
		} else if id >= len(sshHosts) {
			fmt.Println()
			fmt.Println("Enter the appropriate range of ids!")
			fmt.Println()

		} else {
			execline := gokube.KubeExceLine(kubeConfig.PodsList[id], kubeConfig.Namespace, kubeConfig.Container, cmdLine)
			loginHost(sshHosts, sshResult, id, execline)
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Println("👌End Selection!")
}

func showPodList(sshResult []model.RunResult) {
	fmt.Println()
	fmt.Println("✋Pods List:")
	if len(sshResult) != 1 {
		fmt.Printf("Enter the 0~%d to select the pod, other input will exit!\n", len(sshResult)-1)
	} else {
		fmt.Println("Enter the 0 to select the pod, other input will exit!")
	}
	// var status = map[bool]string{false: "\u001b[01;31m[x]\u001b[0m", true: "\u001b[01;32m[√]\u001b[0m"}
	var status = map[bool]string{false: "[x]", true: "[√]"}
	for idx, host := range sshResult {
		fmt.Printf("%-3d :  %-15s  %s\n", idx, host.Host, status[host.Success])
	}
}
