package cmd

import (
	"errors"
	"gosail/cycle"

	"github.com/spf13/cobra"
)

var (
	srcPath  string
	destPath string
	tar      bool
	scp      bool
)
var pullCmd = &cobra.Command{
	Aliases: []string{"download"},
	Use:     "pull",
	Short:   "Pull can copy file from hosts concurrently, and create folders of host to distinguish",
	Long: `
destPath default "."
eg. : gosail login pull "<srcPath>" ["<destPath>"]
eg. : gosail login pull --src "<srcPath>" [--dest "<destPath>"]
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		if scp {
			pull()
		} else {
			download()
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
	loginCmd.AddCommand(pullCmd)
	// copy
	pullCmd.Flags().StringVarP(&srcPath, "src", "", "", "source path")
	pullCmd.Flags().StringVarP(&destPath, "dest", "", "", "destination path")
	pullCmd.Flags().BoolVarP(&scp, "scp", "", false, "pull file by scp")
	pullCmd.Flags().BoolVarP(&tar, "tar", "", false, "tar pull's file")

}
func pull() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.PullAndShow(clientConfig, &srcPath, &destPath, tar)
}

func download() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.DownloadAndShow(clientConfig, &srcPath, &destPath, tar)
}
