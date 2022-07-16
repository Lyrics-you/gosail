package cli

import (
	"gosail/client"
	"gosail/cycle"
	"strings"

	"github.com/desertbit/grumble"
)

var (
	cmdLine   string
	highlight string
	command   []string
)

func init() {
	execCommand := &grumble.Command{
		Name: "exec",
		Help: "Exec can execute commands concurrently and in batches on all hosts and k8s pods",
		Args: func(a *grumble.Args) {
			a.StringList("command", "command line", grumble.Default([]string{}))
		},
		Flags: func(f *grumble.Flags) {
			f.String("e", "cmdline", "", "command line")
			f.String("b", "highlight", "", "bold highlight")
		},
		Run: func(c *grumble.Context) error {
			setExecArgs(c)
			if isK8s {
				k8sExec()
			} else {
				exec()
			}
			return nil
		},
	}
	Gosail.AddCommand(execCommand)
}

func setExecArgs(c *grumble.Context) {
	cmdLine = c.Flags.String("cmdline")
	if cmdLine == "" {
		command = c.Args.StringList("command")
		cmdLine = strings.Join(command, " ")
	}
	highlight = c.Flags.String("highlight")
}

func exec() {
	clientConfig, err = cycle.GetClientConfig("", keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil {
		log.Error(err)
		return
	}
	sshResults := cycle.Exec(clientConfig)
	cycle.ShowExecResult(clientConfig.SshHosts, sshResults, &jsonMode, &linuxMode)
	if selection {
		client.LoginHostByID(clientConfig.SshHosts, sshResults, "")
	}
}
