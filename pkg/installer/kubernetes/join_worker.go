package kubernetes

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"log"
	"os"
)

type JoinWorkerNode struct {
}

func (k *JoinWorkerNode) Exec(config *config.Config) error {

	workerJoinCmd := fmt.Sprintf("%v/join-worker-command", constant.DataPath)
	data, err := os.ReadFile(workerJoinCmd)
	if err != nil {
		return err
	}

	err = initEnv(config, config.Hosts.Workers...)
	if err != nil {
		return err
	}

	t := task.New("Join Worker Node", config.Hosts.Workers...)
	_, err = t.Run(&task.Command{
		Title: "Join Work Node",
		CmdList: []string{
			fmt.Sprintf(string(data)),
		},
		Callback: func(module task.Module, status task.Result) task.Result {
			log.Println(fmt.Sprintf(":\n%v\n", status.Output))
			return status
		},
	})

	return err
}
