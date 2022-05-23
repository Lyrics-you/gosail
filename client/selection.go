package client

import (
	"fmt"
	"gosail/model"
	"gosail/ssh"
	"strconv"
)

// todo:pointer
func LoginHostByID(sshHosts []model.SSHHost, sshResult []model.SSHResult) {
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
		// fmt.Printf("id --------->%v\n", id)
		if id == -1 {
			break
		} else if id >= len(sshHosts) {
			fmt.Println()
			fmt.Println("Enter the appropriate range of ids!")
		} else {
			loginHost(sshHosts, sshResult, id)
			fmt.Println()
		}
	}
	fmt.Println("### End Selection ###")
}

func showHostsList(sshResult []model.SSHResult) {
	fmt.Println("Server List:")
	if len(sshResult) != 1 {
		fmt.Printf("Enter the 0~%d to select the host, other input will exit!\n", len(sshResult)-1)
	} else {
		fmt.Println("Enter the 0 to select the host, other input will exit!")
	}

	var status = map[bool]string{false: "[x]", true: "[√]"}
	for idx, host := range sshResult {
		fmt.Printf("%d : %s %s\n", idx, host.Host, status[host.Success])
	}
}

func selectHost() int {
	var str string
	fmt.Println("Not the first time, you need to press the Enter before selecting")
	fmt.Print("Input id : ")
	fmt.Scanln(&str)

	id, err := strconv.Atoi(str)

	if err != nil {
		// log.Println(err)
		id = -1
	}
	return id
}

func loginHost(sshHosts []model.SSHHost, sshResult []model.SSHResult, id int) {
	host := sshHosts[id]
	fmt.Println()
	if !sshResult[id].Success {
		fmt.Print(sshResult[id].Result)
		return
	}
	fmt.Printf(">>> login into %v ...\n", host.Host)
	err := ssh.GetInteractiveTerminal(host.Username, host.Password, host.Host, host.Key, host.Port, []string{}, []string{})
	if err != nil {
		log.Errorf("err:%v", err)
	}
	fmt.Printf("<<< logout from %v .\n", sshHosts[id].Host)
}
