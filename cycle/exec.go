package cycle

import (
	"encoding/json"
	"fmt"
	"gosail/client"
	"gosail/model"
	"gosail/spinner"
	"gosail/utils"
)

func Exec(clientConfig *model.ClientConfig, spinnerConfig *model.SpinConfig) []model.RunResult {
	if spinnerConfig != nil {
		// user spinner
		spinner.Spin.Init(spinnerConfig)
		spinner.Spin.SetTimeOut(clientConfig.TimeLimit)
		spinner.Spin.Start()
		defer spinner.Spin.Stop()
	}
	sshResults, _ := client.LimitShhWithGroup(clientConfig)
	return sshResults
}

func ShowExecResult(sshHosts []model.SSHHost, sshResults []model.RunResult, jsonMode, linuxMode *bool) {
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
		if !*linuxMode && sshResult.Success {
			sshResult.Result = utils.SimpleLine(sshResult.Result, len(sshHosts[id].CmdList)+2, 3)
		}
		fmt.Println(sshResult.Result)
	}
	fmt.Println("👌Finshed!")
}
