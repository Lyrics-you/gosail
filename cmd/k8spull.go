package cmd

import (
	"errors"
	"gosail/cycle"
	"gosail/model"

	"github.com/spf13/cobra"
)

var k8spullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull",
	Long: `
destPath default "."
eg. : gosail k8s pull "<srcPath>" ["<destPath>"]
eg. : gosail k8s pull --src "<srcPath>" [--dest "<destPath>"]
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		if scp {
			k8sPull()
		} else {
			k8sDownload()
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
	k8sCmd.AddCommand(k8spullCmd)
	// copy
	k8spullCmd.Flags().StringVarP(&srcPath, "src", "", "", "source path")
	k8spullCmd.Flags().StringVarP(&destPath, "dest", "", "", "destination path")
	k8spullCmd.Flags().BoolVarP(&scp, "scp", "", false, "pull file by exec scp")
	k8spullCmd.Flags().BoolVarP(&tar, "tar", "", false, "tar pull's file")
}

func k8sPull() {
	if container == "" {
		container = app
	}
	if app == "" {
		app = container
	}
	if container == "" && app == "" && label == "" {
		log.Errorf("container or app name is not specified")
		return
	}
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, true)
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
	sshResults := cycle.K8sPull(clientConfig, kubeConfig, &srcPath, &destPath, &tar)
	cycle.K8sShowResults(sshResults, kubeConfig, &jsonMode)
}

func k8sDownload() {
	if container == "" {
		container = app
	}
	if app == "" {
		app = container
	}
	if container == "" && app == "" && label == "" {
		log.Errorf("container or app name is not specified")
		return
	}
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, cmdFile, hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, true)
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
	sftpResults := cycle.K8sDownload(clientConfig, kubeConfig, &srcPath, &destPath, &tar)
	cycle.K8sShowResults(sftpResults, kubeConfig, &jsonMode)
}
