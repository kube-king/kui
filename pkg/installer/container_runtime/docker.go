package container_runtime

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"kube-invention/pkg/utils/json_util"
)

type Docker struct {
}

func (d *Docker) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters...)
	hostList = append(hostList, config.Hosts.Workers...)

	t := task.New("Install Docker", hostList...)
	t.SetEnv(map[string]interface{}{
		"REGISTRY":               config.Core.Registry,
		"INSECURE_REGISTRY_LIST": json_util.ToJsonString(config.ContainerRuntimeOption.InsecureRegistryList),
	})

	_, err = t.Run(
		&task.Unarchive{
			Title:         "Copying Docker Binary",
			LocalFilePath: fmt.Sprintf("%v/docker-%v-linux-%v.tar.gz", constant.BinaryPath, config.ContainerRuntimeOption.Version, config.Core.Arch),
			RemoteDir:     "/usr/bin/",
			Mode:          0755,
			Force:         true,
		},
		&task.Script{
			Title:          "Install Docker",
			ScriptFilePath: fmt.Sprintf("%v/install_docker.sh", constant.TemplatePath),
		}, &task.Command{
			Title:   "Check Docker Service",
			CmdList: []string{`systemctl status docker | grep Active | awk '{print $3}' | cut -d "(" -f2 | cut -d ")" -f1`},
			Until: func(output string) bool {
				if output == "running" {
					return true
				}
				return false
			},
			Retries: 5,
			Delay:   3,
		})
	return err
}
