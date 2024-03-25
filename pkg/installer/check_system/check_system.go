package check_system

import (
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
)

type CheckSystem struct {
}

func (c *CheckSystem) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters...)
	hostList = append(hostList, config.Hosts.Workers...)

	t := task.New("Check System Environment", hostList...)
	_, err = t.Run(
		&task.Command{
			Title: "Check Username",
			CmdList: []string{
				"id -u",
			},
			Callback: func(module task.Module, status task.Result) task.Result {
				if status.Output != "0" {
					status.State = task.StateFailed
					status.Message = "error: current user is not root!"
				}
				return status
			},
		},
	)
	return
}
