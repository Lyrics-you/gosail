package cmd

import (
	"errors"
	"gosail/cycle"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pull can copy file to hosts concurrently, and create folders that do not exist",
	Long: `
destPath default "."
eg. : gosail login push "<srcPath>" ["<destPath>"]
eg. : gosail login push --src "<srcPath>" [--dest "<destPath>"]
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		clientConfig, _ = cycle.GetClientConfig(config, keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
		cycle.PushAndShow(clientConfig, &srcPath, &destPath)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if srcPath == "" && destPath == "" && len(args) < 1 {
			return errors.New("pull requires srcpath and destpath 2 args, destpath default is '.'")
		}
		if len(args) > 2 {
			return errors.New("pull receive more than 2 args")
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		if srcPath == "" && len(args) == 1 {
			srcPath = args[0]
		}
		if srcPath == "" && destPath == "" && len(args) == 2 {
			srcPath = args[0]
			destPath = args[1]
		}
	},
}

func init() {
	loginCmd.AddCommand(pushCmd)
	// copy
	pushCmd.Flags().StringVarP(&srcPath, "src", "", "", "exec cmdline")
	pushCmd.Flags().StringVarP(&destPath, "dest", "", "", "exec cmdfile")
}
