package kubernetes

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"log"
	"os"
	"strings"
)

type JoinMasterNode struct {
}

func (j *JoinMasterNode) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters...)

	err = initEnv(config, config.Hosts.Masters[1:]...)
	if err != nil {
		return err
	}

	var vip string
	if config.KubeVipOption.Enable {
		vip = config.KubeVipOption.Vip
	} else {
		vip = config.Hosts.Masters[0].Ip
	}

	t := task.New("Join Master", hostList...)
	t.SetEnv(map[string]interface{}{
		"KUBE_VIP_VERSION": config.KubeVipOption.Version,
		"REGISTRY":         config.Core.Registry,
		"VIP_ADDRESS":      vip,
		"INTERFACE":        config.Core.NetworkAdapter,
	})

	checkClusterStatusTask := make([]task.Module, 0)
	for _, com := range KubernetesComponentList {
		checkClusterStatusTask = append(checkClusterStatusTask, &task.Command{
			Title: fmt.Sprintf("Check Component: %v ", com),
			CmdList: []string{
				fmt.Sprintf(CheckKubernetesClusterStatus, com),
			},
			Until: func(output string) (isSuccess bool) {
				log.Println(fmt.Sprintf("Check Kubernetes Component Result: \n%v\n", strings.ReplaceAll(output, "|", "\n")))
				resList := strings.Split(output, "|")
				if len(resList) <= 0 {
					isSuccess = false
					return
				}

				for _, res := range resList {
					if strings.Trim(res, " ") == "" {
						continue
					}
					val := strings.Split(res, "=")
					if len(val) <= 0 {
						isSuccess = false
						break
					}
					if strings.Trim(val[1], "") != "Running" {
						isSuccess = false
						break
					} else {
						isSuccess = true
					}
				}
				return isSuccess
			},
			Retries:    120,
			Delay:      3,
			DelegateTo: hostList[0].Ip,
		})
	}

	masterJoinCmd := fmt.Sprintf("%v/join-master-command", constant.DataPath)
	data, err := os.ReadFile(masterJoinCmd)
	if err != nil {
		return err
	}

	tasks := make([]task.Module, 0)
	if len(config.Hosts.Masters) > 1 {
		for i := 1; i < len(config.Hosts.Masters); i++ {
			tasks = append(tasks, &task.Command{
				Title:      fmt.Sprintf("join %v node", hostList[i].Hostname),
				CmdList:    []string{fmt.Sprintf("%v && sleep 5", string(data))},
				DelegateTo: hostList[i].Ip,
				Callback: func(module task.Module, status task.Result) task.Result {
					log.Println(fmt.Sprintf("kubernetes join master Node Result: \n%v\n", status.Output))
					return status
				},
			}, &task.Command{
				Title: "Copy kube config file",
				CmdList: []string{
					"cat /etc/kubernetes/admin.conf > /root/.kube/config",
				},
				DelegateTo: hostList[i].Ip,
			})

			if config.KubeVipOption.Enable {
				tasks = append(tasks, &task.Template{
					Title:            "Deploy Kube Vip ",
					TemplateFilePath: fmt.Sprintf("%v/kube-vip-%v.yaml.tpl", constant.TemplatePath, config.KubeVipOption.Version),
					RemoteFilePath:   "/etc/kubernetes/manifests/kube-vip.yaml",
					Force:            true,
					DelegateTo:       hostList[i].Ip,
				})
			}
			tasks = append(tasks, checkClusterStatusTask...)
		}
	} else {
		tasks = append(tasks, checkClusterStatusTask...)
	}
	_, err = t.Run(tasks...)

	return err
}
