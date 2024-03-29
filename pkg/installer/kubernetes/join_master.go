package kubernetes

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"log"
	"os"
)

type JoinMasterNode struct {
}

func (j *JoinMasterNode) Exec(conf *config.Config) (err error) {

	joinMasterHostList := make([]ssh_client.Config, 0)
	if conf.TaskType == config.TaskTypeInit {
		joinMasterHostList = append(joinMasterHostList, conf.Hosts.Masters[1:]...)
	} else if conf.TaskType == config.TaskTypeJoinMaster {
		joinMasterHostList = append(joinMasterHostList, conf.Hosts.Masters...)
	}

	err = initEnv(conf, joinMasterHostList...)
	if err != nil {
		return err
	}

	var vip string
	if conf.KubeVipOption.Enable {
		vip = conf.KubeVipOption.Vip
	} else {
		vip = conf.Hosts.Masters[0].Ip
	}

	t := task.New("Join Master", joinMasterHostList...)
	t.SetEnv(map[string]interface{}{
		"KUBE_VIP_VERSION": conf.KubeVipOption.Version,
		"REGISTRY":         conf.Core.Registry,
		"VIP_ADDRESS":      vip,
		"INTERFACE":        conf.Core.NetworkAdapter,
	})

	masterJoinCmd := fmt.Sprintf("%v/join-master-command", constant.DataPath)
	data, err := os.ReadFile(masterJoinCmd)
	if err != nil {
		return err
	}

	tasks := make([]task.Module, 0)
	tasks = append(tasks, &task.Command{
		Title:   "join master node",
		CmdList: []string{fmt.Sprintf("%v && sleep 5", string(data))},
		Callback: func(module task.Module, status task.Result) task.Result {
			log.Println(fmt.Sprintf("kubernetes join master Node Result: \n%v\n", status.Output))
			return status
		},
	}, &task.Command{
		Title: "Copy kube config file",
		CmdList: []string{
			"cat /etc/kubernetes/admin.conf > /root/.kube/config",
		},
	})

	if conf.KubeVipOption.Enable {
		tasks = append(tasks, &task.Template{
			Title:            "Deploy Kube Vip ",
			TemplateFilePath: fmt.Sprintf("%v/kube-vip-%v.yaml.tpl", constant.TemplatePath, conf.KubeVipOption.Version),
			RemoteFilePath:   "/etc/kubernetes/manifests/kube-vip.yaml",
			Force:            true,
		})
	}

	_, err = t.Run(tasks...)
	return err
}
