package cycle

import (
	"encoding/json"
	"fmt"
	"gosail/client"
	"gosail/gokube"
	"gosail/goscp"
	"gosail/model"
	"gosail/spinner"
)

func getKubePods(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig) []model.KubePods {
	gokube.KubeGetPods(kubeConfig)
	sshResults, _ := client.LimitShhWithGroup(clientConfig)
	kubePods := gokube.GetPodsByResult(kubeConfig, sshResults)
	return kubePods
}

func K8sExec(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig, spinnerConfig *model.SpinConfig) []model.RunResult {
	if spinnerConfig != nil {
		// user spinner
		spinner.Spin.Init(spinnerConfig)
		spinner.Spin.SetTimeOut(clientConfig.TimeLimit)
		spinner.Spin.Start()
		defer spinner.Spin.Stop()
	}
	kubePods := getKubePods(clientConfig, kubeConfig)
	sshResults := []model.RunResult{}
	if kubeConfig.CmdLine != "" {
		kubeHosts := gokube.MakeMultiExecSshHosts(kubePods, clientConfig.SshHosts, kubeConfig.CmdLine)
		clientConfig.SshHosts = kubeHosts
		sshResults, _ = client.LimitShhWithGroup(clientConfig)
	}
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

func K8sPull(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig, srcPath, destPath *string, tar *bool) []model.RunResult {
	if *destPath == "" {
		*destPath = "./"
	}

	spinner.Spin.SetStyle(pstyle)
	spinner.Spin.Start()
	defer spinner.Spin.Stop()

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

func K8sDownloadById(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig, srcPath, destPath *string, tar *bool, id int) []model.RunResult {
	if *destPath == "" {
		*destPath = "./"
	}

	spinner.Spin.SetStyle(pstyle)
	spinner.Spin.Start()
	defer spinner.Spin.Stop()

	kubePods := getKubePods(clientConfig, kubeConfig)

	kubePod := []model.KubePods{}
	count := 0
	for i, pods := range kubePods {
		for j := range pods.PodsName {
			if count == id {
				kubePod = []model.KubePods{kubePods[i]}
				kubePod[0].PodsName = []string{pods.PodsName[j]}
				// kubeConfig.SshHosts = []model.SSHHost{kubeConfig.SshHosts[i]}
				break
			}
			count += 1
		}
	}

	kubeHosts := gokube.MakeMultiCopySshHosts(kubePods, kubeConfig.SshHosts, *srcPath, "./")
	kubePods = kubePod

	// check pods_index
	if id >= len(kubeHosts) {
		log.Errorf("invaild pod index,%d", id)
		return []model.RunResult{}
	}
	kubeHosts = []model.SSHHost{kubeHosts[id]}

	clientConfig.SshHosts = kubeHosts
	client.LimitShhWithGroup(clientConfig)

	scpConfig := model.SCPConfig{
		SshHosts:  kubeHosts,
		TimeLimit: clientConfig.TimeLimit,
		NumLimit:  clientConfig.NumLimit,
		Method:    "PULL",
	}

	kubeConfig.PodsList = []string{kubeConfig.PodsList[id]}
	// srcList :=
	// destList :=
	scpConfig.SrcPath = kubeConfig.PodsList
	scpConfig.DestPath = []string{*destPath}

	// tar file
	if *tar {
		goscp.SecureCopyPullTarFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}
	//copy
	var sftpResults []model.RunResult

	// fmt.Printf("PODS:%+v\n", kubePods)
	// fmt.Printf("HOST:%+v\n", clientConfig.SshHosts)
	// fmt.Printf("SCPCONFIG:%+v\n", &scpConfig)
	// &scpConfig.SrcPath
	sftpResults, _ = client.LimitSftpWithGroup(clientConfig, &scpConfig)

	// delete tar file
	if *tar {
		goscp.SecureCopyPullDelFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}

	// delete k8s master's file
	delHosts := gokube.MakeMultiDeleteSshHosts(kubePods, clientConfig.SshHosts, "./")

	// delHosts = []model.SSHHost{delHosts[id]}

	clientConfig.SshHosts = delHosts
	client.LimitShhWithGroup(clientConfig)
	return sftpResults
}

func K8sDownload(clientConfig *model.ClientConfig, kubeConfig *model.KubeConfig, srcPath, destPath *string, tar *bool) []model.RunResult {
	if *destPath == "" {
		*destPath = "./"
	}

	spinner.Spin.SetStyle(pstyle)
	spinner.Spin.Start()
	defer spinner.Spin.Stop()

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

	// tar file
	if *tar {
		goscp.SecureCopyPullTarFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}
	//copy
	var sftpResults []model.RunResult

	sftpResults, _ = client.LimitSftpWithGroup(clientConfig, &scpConfig)

	// delete tar file
	if *tar {
		goscp.SecureCopyPullDelFile(&scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}

	// delete k8s master's file
	delHosts := gokube.MakeMultiDeleteSshHosts(kubePods, clientConfig.SshHosts, "./")
	clientConfig.SshHosts = delHosts
	client.LimitShhWithGroup(clientConfig)
	return sftpResults
}
