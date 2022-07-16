package cmd

import (
	"errors"
	"fmt"
	"gosail/model"

	"github.com/spf13/cobra"
)

var (
	isDesc bool
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version subcommand show gosail version info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized command"))
			return
		}
		if !isDesc {
			showVersion()
		} else {
			showDesciption()
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVarP(&isDesc, "description", "d", false, "history description")
}

func printTag(tag, value string) {
	fmt.Printf("%-13s : %s\n", tag, value)
}

func showVersion() {
	version := model.Historys[len(model.Historys)-1]
	printTag("Name", fmt.Sprintf("gosail%s", model.EMOJI["gosail"]))
	printTag("Version", version.Version)
	printTag("Email", "Leyuan.Jia@Outlook.com")
}

func showDesciption() {
	version := model.Historys[len(model.Historys)-1]
	printTag("Name", fmt.Sprintf("gosail%s", model.EMOJI["gosail"]))
	printTag("Version", version.Version)
	printTag("Description", version.Description)
}
