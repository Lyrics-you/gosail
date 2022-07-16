package cmd

import (
	"gosail/client"
	"gosail/cycle"
	"gosail/model"

	"github.com/spf13/cobra"
)

var k8sexecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec can execute commands concurrently and in batches on all specified pods",
	Long: `
eg. : gosail k8s exec [-e] "<cmdline>"
eg. : gosail k8s exec -e "<cmdline>" mode [flags]
eg. : gosail k8s exec [-e] "<cmdline>" -b "<highlight>"
	`,
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		linuxMode = true
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		k8sExec()
	},
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 1 {
			cmdLine = args[0]
		}
		if len(args) == 2 {
			cmdLine = args[0]
			highlight = args[1]
		}
	},
}

func init() {
	k8sCmd.AddCommand(k8sexecCmd)
	//
	k8sexecCmd.PersistentFlags().StringVarP(&cmdLine, "cmdline", "e", "", "exec cmdline")
	k8sexecCmd.PersistentFlags().StringVarP(&highlight, "highlight", "b", "", "bold highlight")
}

func k8sExec() {
	clientConfig, err := cycle.GetClientConfig(config, keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil {
		log.Error(err)
		return
	}
	kubeConfig := &model.KubeConfig{
		SshHosts:  clientConfig.SshHosts,
		Namespace: namespace,
		App:       app,
		Container: container,
		Label:     label,
		Shell:     shell,
		Highlight: highlight,
		CmdLine:   cmdLine,
	}
	sshResults := cycle.K8sExec(clientConfig, kubeConfig)
	cycle.K8sShowResults(sshResults, kubeConfig, &jsonMode)
	if selection {
		client.LoginPodByID(kubeConfig, clientConfig.SshHosts, sshResults, kubeConfig.Shell)
	}
}
