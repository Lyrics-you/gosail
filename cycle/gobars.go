package cycle

import (
	"encoding/json"
	"fmt"
	"gosail/client"
	"gosail/gokube"
	"gosail/goscp"
	"gosail/model"
)

func getKubePods(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig) []model.KubePods {
	gokube.KubeGetPods(kubeConfig)
	sshResults, _ := client.LimitShhWithGroup(clientConfig)
	kubePods := gokube.GetPodsByResult(kubeConfig, sshResults)
	return kubePods
}

func K8sExec(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig) []model.RunResult {
	kubePods := getKubePods(clientConfig, kubeConfig)
	sshResults := []model.RunResult{}
	if kubeConfig.CmdLine != "" {
		kubeHosts := gokube.MakeMultiExecSshHosts(kubePods, clientConfig.SshHosts, kubeConfig.CmdLine)
		clientConfig.SshHosts = kubeHosts
		sshResults, _ = client.LimitShhWithGroup(clientConfig)
	}
	return sshResults
}

func K8sPull(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig, srcPath, destPath *string, tar *bool) []model.RunResult {
	if *destPath == "" {
		*destPath = "./"
	}
	kubePods := getKubePods(clientConfig, kubeConfig)
	kubeHosts := gokube.MakeMultiCopySshHosts(kubePods, kubeConfig.SshHosts, *srcPath, "./")
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
		destList = append(destList, *destPath)
	}
	scpConfig.SrcPath = srcList
	scpConfig.DestPath = destList

	mkdirResults := goscp.SecureCopyPullMakeDir(&scpConfig)

	// tar file
	if *tar {
		goscp.SecureCopyPullTarFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}
	//copy
	var sshResults []model.RunResult
	sshResults, _ = client.LimitScpWithGroup(&scpConfig, mkdirResults)
	// delete tar file
	if *tar {
		goscp.SecureCopyPullDelFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}

	// delete k8s master's file
	delHosts := gokube.MakeMultiDeleteSshHosts(kubePods, clientConfig.SshHosts, "./")
	clientConfig.SshHosts = delHosts
	client.LimitShhWithGroup(clientConfig)
	return sshResults
}

func K8sShowResults(sshResults []model.RunResult, kubeConfig *model.KubeConfig, jsonMode *bool) {
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
		fmt.Println(sshResult.Result)

	}
	fmt.Println("👌Finshed!")
}
