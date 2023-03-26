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
				getPodsLineByLabel(kubeConfig.Namespace, kubeConfig.Label),
			}
		} else {
			kubeConfig.SshHosts[id].CmdList = []string{
				getPodsLineByGrep(kubeConfig.Namespace, kubeConfig.App),
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
			AppName:    kubeConfig.App,
			Container:  kubeConfig.Container,
			Shell:      kubeConfig.Shell,
			Highlight:  kubeConfig.Highlight,
		})
		kubeConfig.PodsList = append(kubeConfig.PodsList, kubePods[id].PodsName...)
	}
	return kubePods
}

func getPodsLineByGrep(namespace, appname string) string {
	// offer : lujun
	line := fmt.Sprintf("kubectl get pods -n %s | grep Running | awk '{print $1}' | grep -P '^%s(-\\w+){1,2}$'", namespace, appname)
	return line
}

func getPodsLineByLabel(namespace, label string) string {
	// line := fmt.Sprintf("kubectl get pods -n %s -l %s | awk 'NR!=1 {print $1}'", namespace, label)
	line := fmt.Sprintf("kubectl get pods -n %s -l %s | grep Running | awk '{print $1}'", namespace, label)
	return line
}

// func GetPodsLineByLabelApp(namespace, appname string) string {
// 	line := fmt.Sprintf("kubectl get pods -n %s -l app=%s | awk 'NR!=1 {print $1}'", namespace, appname)
// 	return line
// }

func MakeMultiExecSshHosts(kubePods []model.KubePods, masterHosts []model.SSHHost, cmdline string) []model.SSHHost {
	kubeHosts := []model.SSHHost{}
	for id, pod := range kubePods {
		masterHost := masterHosts[id]
		for _, name := range pod.PodsName {
			masterHost.CmdLine = KubeExceLine(name, pod.Namespace, pod.Container, pod.Shell, cmdline)
			if pod.Highlight != "" {
				masterHost.CmdLine = PerlHightlight(masterHost.CmdLine, pod.Highlight)
			}
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
				kubeCopyLine(pod.Namespace, name, pod.Container, srcPath, tagPath),
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
				deleteFileLine(tagPath),
			}
			kubeHosts = append(kubeHosts, masterHost)
		}
	}
	return kubeHosts
}

func KubeExceLine(podname, namespace, container, shell, cmdline string) string {
	// kubectl exec -it -n namespace -c container -- /bin/bash -c 'command'
	line := fmt.Sprintf("kubectl exec -it %s -n %s -c %s -- /bin/%s -c '%s'", podname, namespace, container, shell, cmdline)
	return line
}

func PerlHightlight(line, key string) string {
	// line | perl -pe "s/(${key})/\e[1;31m\$1\e[0m/g"
	hline := fmt.Sprintf(`%s | perl -pe "s/(%s)/\e[1;31m\$1\e[0m/g"`, line, key)
	return hline
}

func kubeCopyLine(namespace, podname, container, srcPath, destPath string) string {
	// kubectl cp -c container srcPath destPath
	// src=[namespace/[pod]]:file/path dest
	filename := utils.GetPathLastName(srcPath)
	line := fmt.Sprintf("kubectl cp -c %s %s/%s:%s %s/%s", container, namespace, podname, srcPath, destPath, filename)
	return line
}

func deleteFileLine(filename string) string {
	line := fmt.Sprintf("rm -rf %s", filename)
	return line
}
