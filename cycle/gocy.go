package cycle

import (
	"fmt"
	"gosail/client"
	"gosail/goscp"
	"gosail/model"
)

func PullAndShow(clientConfig *model.ClientConfig, srcPath, destPath *string, tar bool) {
	scpConfig := &model.SCPConfig{
		SshHosts:  clientConfig.SshHosts,
		TimeLimit: clientConfig.TimeLimit,
		NumLimit:  clientConfig.NumLimit,
	}
	var scpResults = []model.RunResult{}
	var srcList, destList []string
	for i := 0; i < len(clientConfig.SshHosts); i++ {
		srcList = append(srcList, *srcPath)
		destList = append(destList, *destPath)
	}
	scpConfig.SrcPath = srcList
	scpConfig.DestPath = destList
	scpConfig.Method = "PULL"

	mkdirResults := goscp.SecureCopyPullMakeDir(scpConfig)

	var tarResults = []model.RunResult{}
	// tar file
	if tar {
		goscp.SecureCopyPullTarFile(scpConfig)
		tarResults, _ = client.LimitShhWithGroup(clientConfig)
	}
	//copy
	scpResults, _ = client.LimitScpWithGroup(scpConfig, mkdirResults)
	// delete tar file
	if tar {
		goscp.SecureCopyPullDelFile(scpConfig)
		client.LimitShhWithGroup(clientConfig)
	}

	for id, scpResult := range scpResults {
		fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", clientConfig.SshHosts[id].Username, "localhost", id)
		if mkdirResults[id].Success {
			fmt.Print(mkdirResults[id].Result)
		}
		if tar && tarResults[id].Success {
			fmt.Print(tarResults[id].Result)
		}
		fmt.Print(scpResult.Result)
		fmt.Println()
	}
	fmt.Println("👌Finshed!")
}

func PushAndShow(clientConfig *model.ClientConfig, srcPath, destPath *string) {
	scpConfig := &model.SCPConfig{
		SshHosts:  clientConfig.SshHosts,
		TimeLimit: clientConfig.TimeLimit,
		NumLimit:  clientConfig.NumLimit,
	}
	var scpResults = []model.RunResult{}
	var srcList, destList []string
	for i := 0; i < len(clientConfig.SshHosts); i++ {
		srcList = append(srcList, *srcPath)
		destList = append(destList, *destPath)
	}
	scpConfig.SrcPath = srcList
	scpConfig.DestPath = destList
	scpConfig.Method = "PUSH"

	goscp.SecureCopyPushMakeDir(scpConfig)
	// scpConfig'SshHosts is equal clientConfig'SshHosts
	sshResults, _ := client.LimitShhWithGroup(clientConfig)

	scpResults, _ = client.LimitScpWithGroup(scpConfig, sshResults)
	for id, scpResult := range scpResults {
		fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", clientConfig.SshHosts[id].Username, "localhost", id)
		if sshResults[id].Success {
			fmt.Print(sshResults[id].Result)
		}
		fmt.Print(scpResult.Result)
		fmt.Println()
	}
	fmt.Println("👌Finshed!")
}
