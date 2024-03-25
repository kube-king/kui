package init_system

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
)

type InitSystem struct {
}

func (c *InitSystem) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters...)
	hostList = append(hostList, config.Hosts.Workers...)

	t := task.New("Init System", hostList...)
	_, err = t.Run(
		&task.Script{
			Title:          "Copy Init System Script",
			ScriptFilePath: fmt.Sprintf("%v/init_system.sh", constant.TemplatePath),
		},
	)
	return err
}
