package container_runtime

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
)

type Containerd struct {
}

func (c *Containerd) GetPauseVersion(kubernetesVersion string) (pauseVersion string) {

	switch kubernetesVersion {
	case "v1.21.5":
		pauseVersion = "3.4.1"
	case "v1.22.17":
		pauseVersion = "3.5"
	case "v1.23.17":
		pauseVersion = "3.6"
	case "v1.24.0":
		pauseVersion = "3.7"
	case "v1.25.0":
		pauseVersion = "3.8"
	case "v1.26.0":
		pauseVersion = "3.9"
	case "v1.27.0":
		pauseVersion = "3.9"
	case "v1.28.0":
		pauseVersion = "3.9"
	}

	return
}

func (c *Containerd) Exec(config *config.Config) (err error) {

	hostList := make([]ssh_client.Config, 0)
	hostList = append(hostList, config.Hosts.Masters...)
	hostList = append(hostList, config.Hosts.Workers...)

	t := task.New("Install Containerd", hostList...)

	t.SetEnv(map[string]interface{}{
		"PAUSE_VERSION":          c.GetPauseVersion(config.KubernetsOption.Version),
		"INSECURE_REGISTRY_LIST": config.ContainerRuntimeOption.InsecureRegistryList,
		"VERSION":                config.ContainerRuntimeOption.Version,
		"REGISTRY":               config.Core.Registry,
	})

	_, err = t.Run(
		&task.Unarchive{
			Title:         "Copy Binary File",
			LocalFilePath: fmt.Sprintf("%v/containerd-%v-linux-%v.tar.gz", constant.BinaryPath, config.ContainerRuntimeOption.Version, config.Core.Arch),
			RemoteDir:     "/usr/bin/",
			Mode:          0755,
			Force:         true,
		},
		&task.Script{
			Title:          "Deploy Containerd",
			ScriptFilePath: fmt.Sprintf("%v/install_containerd.sh", constant.TemplatePath),
		}, &task.Command{
			Title:   "Check Containerd Service",
			CmdList: []string{`systemctl status containerd | grep Active | awk '{print $3}' | cut -d "(" -f2 | cut -d ")" -f1`},
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
