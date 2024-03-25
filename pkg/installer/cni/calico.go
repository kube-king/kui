package cni

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
)

type Calico struct {
}

func (c *Calico) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters[0])

	t := task.New("Deploy Calico", hostList...)
	t.SetEnv(map[string]interface{}{
		"POD_SUBNET": config.KubernetsOption.PodCIDR,
		"VERSION":    config.CniOption.Version,
		"INTERFACE":  config.Core.NetworkAdapter,
		"REGISTRY":   config.Core.Registry,
	})
	_, err = t.Run(
		&task.Template{
			Title:            "Copy Calico yaml",
			TemplateFilePath: fmt.Sprintf("%v/calico-%v.yaml.tpl", constant.TemplatePath, config.CniOption.Version),
			RemoteFilePath:   fmt.Sprintf("/tmp/calico-%v.yaml", config.CniOption.Version),
			Force:            true,
		}, &task.Command{
			Title: "Deploy Calico",
			CmdList: []string{
				fmt.Sprintf("kubectl apply -f /tmp/calico-%v.yaml", config.CniOption.Version),
			},
		},
	)

	return
}
