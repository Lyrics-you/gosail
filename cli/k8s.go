package cli

import (
	"gosail/client"
	"gosail/cycle"
	"gosail/model"

	"github.com/desertbit/grumble"
)

var (
	namespace string
	app       string
	container string
	label     string
	shell     = "sh"
	isK8s     = false
)

func setK8sArgs(c *grumble.Context) {
	isK8s = true
	namespace = c.Flags.String("namespace")
	app = c.Flags.String("app")
	container = c.Flags.String("container")
}

func k8sExec() {
	linuxMode = true
	clientConfig, err := cycle.GetClientConfig("", keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
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

func k8sPull() {
	linuxMode = true
	clientConfig, _ = cycle.GetClientConfig("", keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
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
