package cycle

import (
	"encoding/json"
	"fmt"
	"gosail/client"
	"gosail/model"
	"strings"
)

func Exec(clientConfig *model.ClientConfig) []model.RunResult {
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
		fmt.Printf("👇===============> %4s@%-15s <===============[%-3d]\n", sshHosts[id].Username, sshResult.Host, id)
		if !*linuxMode && sshResult.Success {
			sshResult.Result = simpleline(sshResult.Result, len(sshHosts[id].CmdList)+2, 3)
		}
		fmt.Println(sshResult.Result)
	}
	fmt.Println("👌Finshed!")
}

func simpleline(str string, n int, m int) string {
	for i := 0; i < n; i++ {
		s := strings.Index(str, "\n")
		if s != -1 {
			str = str[s+1:]
		} else {
			break
		}

	}
	for i := 0; i < m; i++ {
		e := strings.LastIndex(str, "\n")
		if e != -1 {
			str = str[:e]
		} else {
			break
		}
	}
	str += "\n"
	return str
}
