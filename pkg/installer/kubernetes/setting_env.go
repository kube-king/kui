package kubernetes

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
)

func initEnv(config *config.Config, targetHostList ...ssh_client.Config) error {

	t := task.New("Set Kubernetes Env", targetHostList...)
	_, err := t.Run(&task.LineInFile{
		Title: "Set host file",
		Path:  "/etc/hosts",
		Lines: []task.Line{
			{
				Line:    "{{ .ip }} {{ .hostname }}",
				Pattern: "{{ .ip }}",
				State:   "present",
			},
		},
	},
		&task.Command{
			Title: "Set hostname ",
			CmdList: []string{
				"echo {{ .hostname }} > /etc/hostname",
				"hostname {{ .hostname }}",
			},
		},
		&task.File{
			Title: "Create manifests directory",
			Type:  "directory",
			Paths: []string{
				"/etc/kubernetes/manifests",
				"/usr/lib/systemd/system/kubelet.service.d",
			},
		},
		&task.Unarchive{
			Title:         "Copy Binary",
			LocalFilePath: fmt.Sprintf("%v/kubernetes-%v-linux-%v.tar.gz", constant.BinaryPath, config.KubernetsOption.Version, config.Core.Arch),
			RemoteDir:     "/usr/bin/",
			Mode:          0775,
			Force:         true,
		},
		&task.File{
			Title: "Create Kubelet Path",
			Type:  "directory",
			Paths: []string{
				"/var/lib/kubelet",
			},
		}, &task.Command{
			Title: "Config Systemd",
			CmdList: []string{
				KubeletService,
			},
			IgnoreError: true,
		}, &task.Command{
			Title: "Config kubeadm Systemd Service",
			CmdList: []string{
				KubeadmConf,
			},
		}, &task.File{
			Title: "Create .kube Directory",
			Type:  "directory",
			Mode:  0755,
			Paths: []string{"/root/.kube"},
		})
	if err != nil {
		return err
	}

	return err
}
