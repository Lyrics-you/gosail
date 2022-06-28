package gokube

import (
	"fmt"
	"gosail/model"
	"gosail/utils"
)

func KubeGetPods(kubeConfig *model.KubeConfig) {
	for id := range kubeConfig.SshHosts {
		if kubeConfig.Label != "" {
			kubeConfig.SshHosts[id].CmdList = []string{
				GetPodsLineByLabel(kubeConfig.Namespace, kubeConfig.Label),
			}
		} else {
			kubeConfig.SshHosts[id].CmdList = []string{
				GetPodsLineByGrep(kubeConfig.Namespace, kubeConfig.AppName),
			}
		}
	}
}

func GetPodsByResult(kubeConfig *model.KubeConfig, sshResults []model.RunResult) []model.KubePods {
	kubePods := []model.KubePods{}
	for id, sshResult := range sshResults {
		kubePods = append(kubePods, model.KubePods{
			MasterHost: sshResult.Host,
			PodsName:   utils.SplitStringLine(sshResult.Result),
			Namespace:  kubeConfig.Namespace,
			AppName:    kubeConfig.AppName,
			Container:  kubeConfig.Container,
		})
		kubeConfig.PodsList = append(kubeConfig.PodsList, kubePods[id].PodsName...)
	}
	return kubePods
}

func GetPodsLineByGrep(namespace, appname string) string {
	// offer : lujun
	line := fmt.Sprintf("kubectl get pods -n %s | awk '{print $1}' | grep -P '%s(-\\w+){1,2}$'", namespace, appname)
	return line
}

func GetPodsLineByLabel(namespace, label string) string {
	line := fmt.Sprintf("kubectl get pods -n %s -l %s | awk 'NR!=1 {print $1}'", namespace, label)
	return line
}

func GetPodsLineByLabelApp(namespace, appname string) string {
	line := fmt.Sprintf("kubectl get pods -n %s -l app=%s | awk 'NR!=1 {print $1}'", namespace, appname)
	return line
}

func MakeMultiExecSshHosts(kubePods []model.KubePods, masterHosts []model.SSHHost, cmdline string) []model.SSHHost {
	kubeHosts := []model.SSHHost{}
	for id, pod := range kubePods {
		masterHost := masterHosts[id]
		for _, name := range pod.PodsName {
			masterHost.CmdLine = KubeExceLine(name, pod.Namespace, pod.Container, cmdline)
			masterHost.CmdList = []string{masterHost.CmdLine}
			kubeHosts = append(kubeHosts, masterHost)
		}
	}
	return kubeHosts
}

func MakeMultiCopySshHosts(kubePods []model.KubePods, masterHosts []model.SSHHost, srcPath, destPath string) []model.SSHHost {
	if destPath == "" {
		destPath = "./"
	}
	kubeHosts := []model.SSHHost{}
	for id, pod := range kubePods {
		masterHost := masterHosts[id]
		for _, name := range pod.PodsName {
			tagPath := fmt.Sprintf("%s/%s", destPath, name)
			masterHost.CmdList = []string{
				KubeCopyLine(pod.Namespace, name, pod.Container, srcPath, tagPath),
			}
			kubeHosts = append(kubeHosts, masterHost)
		}
	}
	return kubeHosts
}

func MakeMultiDeleteSshHosts(kubePods []model.KubePods, masterHosts []model.SSHHost, destPath string) []model.SSHHost {
	kubeHosts := []model.SSHHost{}
	for id, pod := range kubePods {
		masterHost := masterHosts[id]
		for _, name := range pod.PodsName {
			tagPath := fmt.Sprintf("%s/%s", destPath, name)
			masterHost.CmdList = []string{
				DeleteFileLine(tagPath),
			}
			kubeHosts = append(kubeHosts, masterHost)
		}
	}
	return kubeHosts
}

func KubeExceLine(podname, namespace, container, cmdline string) string {
	// kubectl exec -it -n namespace -c container -- /bin/bash -c 'command'
	line := fmt.Sprintf("kubectl exec -it %s -n %s -c %s -- /bin/bash -c '%s'", podname, namespace, container, cmdline)
	return line
}

func KubeCopyLine(namespace, podname, container, srcPath, destPath string) string {
	// kubectl cp -c container srcPath destPath
	// src=[namespace/[pod]]:file/path dest
	line := fmt.Sprintf("kubectl cp -c %s %s/%s:%s %s", container, namespace, podname, srcPath, destPath)
	return line
}

func DeleteFileLine(filename string) string {
	line := fmt.Sprintf("rm -rf %s", filename)
	return line
}
