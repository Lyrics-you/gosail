package client

import (
	"fmt"
	"gosail/model"
	"gosail/ssh"
	"strconv"
)

func LoginHostByID(sshHosts []model.SSHHost) {
	id := 0
	for id >= 0 && len(sshHosts) > 0 {
		showHostsList(sshHosts)
		id := selectHost()
		if id == -1 {
			break
		} else if id >= len(sshHosts) {
			fmt.Println()
			fmt.Println("Enter the appropriate range of ids!")
		} else {
			loginHost(sshHosts, id)
			fmt.Printf("<<< logout from %v .\n", sshHosts[id].Host)
			fmt.Println()
		}
	}
	fmt.Println("### End Selection ###")
}

func showHostsList(sshHosts []model.SSHHost) {
	fmt.Println("Server List:")
	if len(sshHosts) != 1 {
		fmt.Printf("Enter the 0~%d to select the host, other input will exit!\n", len(sshHosts)-1)
	} else {
		fmt.Println("Enter the 0 to select the host, other input will exit!")
	}

	for idx, host := range sshHosts {
		fmt.Printf("%d : %s\n", idx, host.Host)
	}
}

func selectHost() int {
	var str string
	fmt.Print("Input id : ")
	fmt.Scanln(&str)
	id, err := strconv.Atoi(str)
	if err != nil {
		id = -1
	}
	return id
}

func loginHost(sshHosts []model.SSHHost, id int) {
	host := sshHosts[id]
	fmt.Println()
	fmt.Printf(">>> login into %v ...\n", host.Host)
	err := ssh.GetInteractiveTerminal(host.Username, host.Password, host.Host, host.Key, host.Port, []string{}, []string{})
	if err != nil {
		fmt.Printf("err:%v", err)
	}
}
