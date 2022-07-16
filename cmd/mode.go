package cmd

import (
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "mode",
	Short: "Mode offers choices of exec output formats",
	Long: `
-j : use jsonmode to make the outpout with json format 
-l : use linuxmode to make the output without the hostname
-s : use selection to login hosts by their id
`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	execCmd.AddCommand(modelCmd)
	// model
	modelCmd.Flags().BoolVarP(&jsonMode, "jsonmode", "j", false, "json mode")
	modelCmd.Flags().BoolVarP(&linuxMode, "linuxmode", "l", false, "linux mode")
	modelCmd.Flags().BoolVarP(&selection, "selection", "s", false, "select host to login")
}
