package cmd

import (
	"github.com/spf13/cobra"
)

var (
	namespace string
	app       string
	container string
	label     string
	highlight string
	shell     string
)
var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "K8s master to do something",
	Long: `
eg. : gosail login k8s -n "<namespace>" -a "<deployment.app>" [-c "<container>"]
`,
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	loginCmd.AddCommand(k8sCmd)
	// k8s
	k8sCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "k8s namespace")
	k8sCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "k8s deployment app")
	k8sCmd.PersistentFlags().StringVarP(&container, "container", "c", "", "deployment container")
	k8sCmd.PersistentFlags().StringVarP(&label, "label", "l", "", "deployment label")
	k8sCmd.PersistentFlags().StringVarP(&shell, "shell", "", "sh", "container shell")
}
