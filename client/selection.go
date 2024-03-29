package client

import (
	"fmt"
	"gosail/gossh"
	"gosail/model"
	"strconv"
)

func LoginHostByID(sshHosts []model.SSHHost, sshResult []model.RunResult, cmdLine string) {
	var id int
	if len(sshHosts) == 0 {
		id = -1
	} else {
		id = 0
	}

	// can be cyclically selected
	for id >= 0 {
		showHostsList(sshResult)
		id := selectHost()
		if id == -1 {
			break
		} else if id >= len(sshHosts) {
			fmt.Println()
			fmt.Println("Enter the appropriate range of ids!")
			fmt.Println()

		} else {
			loginHost(sshHosts, sshResult, id, cmdLine)
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Println("👌End Selection!")
}

func showHostsList(sshResult []model.RunResult) {
	fmt.Println()
	fmt.Println("✋Server List:")
	if len(sshResult) == 0 {
		fmt.Println("No No available servers!")
		return
	}
	if len(sshResult) != 1 {
		fmt.Printf("Enter the 0~%d to select the host, other input will exit!\n", len(sshResult)-1)
	} else {
		fmt.Println("Enter the 0 to select the host, other input will exit!")
	}
	// var status = map[bool]string{false: "\u001b[01;31m[x]\u001b[0m", true: "\u001b[01;32m[√]\u001b[0m"}
	var status = map[bool]string{false: "[x]", true: "[√]"}
	for idx, host := range sshResult {
		fmt.Printf("%-3d : %5s@%-15s  %s\n", idx, host.Username, host.Host, status[host.Success])
	}
}

func selectHost() int {
	var str string
	fmt.Print("Input id : ")
	fmt.Scanln(&str)

	id, err := strconv.Atoi(str)

	if err != nil {
		// log.Println(err)
		id = -1
	}
	return id
}

func loginHost(sshHosts []model.SSHHost, sshResults []model.RunResult, id int, cmdLine string) {
	host := sshHosts[id]

	if !sshResults[id].Success {
		fmt.Println()
		fmt.Print(sshResults[id].Result)
		return
	}

	fmt.Println()
	fmt.Printf(">>> login into %v ...\n", sshResults[id].Host)
	err := gossh.GetInteractiveTerminal(host.Username, host.Password, host.Host, host.Key, host.Port, cmdLine, []string{}, []string{})
	if err != nil {
		log.Errorf("err:%v", err)
	}
	fmt.Printf("\n<<< logout from %v .", sshResults[id].Host)
}
