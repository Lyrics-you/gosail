package cmd

import (
	"errors"
	"gosail/cycle"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Aliases: []string{"upload"},
	Use:     "push",
	Short:   "Pull can copy file to hosts concurrently, and create folders that do not exist",
	Long: `
destPath default "."
eg. : gosail login push "<srcPath>" ["<destPath>"]
eg. : gosail login push --src "<srcPath>" [--dest "<destPath>"]
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		if scp {
			push()
		} else {
			upload()
		}
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
	pushCmd.Flags().StringVarP(&srcPath, "src", "", "", "source path")
	pushCmd.Flags().StringVarP(&destPath, "dest", "", "", "destination path")
	pushCmd.Flags().BoolVarP(&scp, "scp", "", false, "push file by scp")
}
func push() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.PushAndShow(clientConfig, &srcPath, &destPath)
}

func upload() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.UploadAndShow(clientConfig, &srcPath, &destPath)
}
