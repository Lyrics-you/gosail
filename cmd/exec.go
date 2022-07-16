package cmd

import (
	"gosail/client"
	"gosail/cycle"

	"github.com/spf13/cobra"
)

var (
	cmdLine string
	cmdFile string
)
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec can execute commands concurrently and in batches on all hosts",
	Long: `
eg. : gosail login exec [-e] "<cmdline>"
eg. : gosail login exec -e "<cmdline>" mode [flags]
eg. : gosail login exec --cmdfile "<cmdfile>"
`,
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		exec()
	},
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			cmdLine = args[0]
		}
	},
}

func init() {
	loginCmd.AddCommand(execCmd)
	// model
	execCmd.PersistentFlags().StringVarP(&cmdLine, "cmdline", "e", "", "exec cmdline")
	execCmd.PersistentFlags().StringVarP(&cmdFile, "cmdfile", "", "", "exec cmdfile")
}

func exec() {
	clientConfig, err := cycle.GetClientConfig(config, keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
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
