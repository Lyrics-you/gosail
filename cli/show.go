package cli

import (
	"fmt"

	"github.com/desertbit/grumble"
)

var (
	hostList []string
)

func init() {
	showCommand := &grumble.Command{
		Name: "show",
		Help: "Show the hosts",
		Args: func(a *grumble.Args) {

		},
		Flags: func(f *grumble.Flags) {

		},
		Run: func(_ *grumble.Context) error {
			show()
			return nil
		},
	}
	Gosail.AddCommand(showCommand)
}

func show() {
	if !isK8s {
		fmt.Println("Hosts:")
	} else {
		fmt.Println("K8s master:")
	}
	for _, host := range sshHosts {
		fmt.Printf("%s@%s\n", host.Username, host.Host)
	}
}
