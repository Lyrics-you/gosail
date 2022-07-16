package cmd

import (
	"github.com/spf13/cobra"
)

var k8smodelCmd = &cobra.Command{
	Use:   "mode",
	Short: "Mode ",
	Long: `
-j : use jsonmode to make the outpout with json format 
-s : use selection to login pods by their id
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	k8sexecCmd.AddCommand(k8smodelCmd)
	// model
	k8smodelCmd.Flags().BoolVarP(&jsonMode, "jsonmode", "j", false, "json mode")
	k8smodelCmd.Flags().BoolVarP(&selection, "selection", "s", false, "select host to login")
}
