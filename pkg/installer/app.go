package installer

import (
	"kube-invention/pkg/installer/check_system"
	"kube-invention/pkg/installer/cni"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/container_runtime"
	"kube-invention/pkg/installer/download"
	"kube-invention/pkg/installer/etcd"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/installer/init_system"
	"kube-invention/pkg/installer/kubernetes"
)

type Installer interface {
	Exec(config *config.Config) error
}

var InstallJobApp = new(InstallJob)

type InstallJob struct {
}

func (e *InstallJob) runTask(config *config.Config, ins ...Installer) {
	for _, i := range ins {
		err := i.Exec(config)
		if err != nil {
			global.Log.Error(err.Error())
			break
		}
	}
}

func (e *InstallJob) JoinWorkerNode(conf *config.Config) {

	conf.TaskType = config.TaskTypeJoinWorker
	installerTasks := make([]Installer, 0)
	installerTasks = append(installerTasks, &download.DownloadData{})
	installerTasks = append(installerTasks, &check_system.CheckSystem{})
	installerTasks = append(installerTasks, &init_system.InitSystem{})
	switch conf.ContainerRuntimeOption.Type {
	case "docker":
		installerTasks = append(installerTasks, &container_runtime.Docker{})
	case "containerd":
		installerTasks = append(installerTasks, &container_runtime.Containerd{})
	}

	installerTasks = append(installerTasks, &kubernetes.JoinWorkerNode{})

	e.runTask(conf, installerTasks...)
}

func (e *InstallJob) JoinMasterNode(conf *config.Config) {

	conf.TaskType = config.TaskTypeJoinMaster
	installerTasks := make([]Installer, 0)
	installerTasks = append(installerTasks, &download.DownloadData{})
	installerTasks = append(installerTasks, &check_system.CheckSystem{})
	installerTasks = append(installerTasks, &init_system.InitSystem{})
	switch conf.ContainerRuntimeOption.Type {
	case "docker":
		installerTasks = append(installerTasks, &container_runtime.Docker{})
	case "containerd":
		installerTasks = append(installerTasks, &container_runtime.Containerd{})
	}
	installerTasks = append(installerTasks, &kubernetes.JoinMasterNode{})

	e.runTask(conf, installerTasks...)
}

func (e *InstallJob) InitCluster(conf *config.Config) {

	conf.TaskType = config.TaskTypeInit
	installerTasks := make([]Installer, 0)
	installerTasks = append(installerTasks, &download.DownloadData{})
	installerTasks = append(installerTasks, &check_system.CheckSystem{})
	installerTasks = append(installerTasks, &init_system.InitSystem{})
	switch conf.ContainerRuntimeOption.Type {
	case "docker":
		installerTasks = append(installerTasks, &container_runtime.Docker{})
	case "containerd":
		installerTasks = append(installerTasks, &container_runtime.Containerd{})
	}
	installerTasks = append(installerTasks, &etcd.Etcd{})
	installerTasks = append(installerTasks, &kubernetes.InitCluster{})
	installerTasks = append(installerTasks, &kubernetes.JoinMasterNode{})
	installerTasks = append(installerTasks, &kubernetes.JoinWorkerNode{})

	if conf.CniOption.Enable {
		switch conf.CniOption.Type {
		case "calico":
			installerTasks = append(installerTasks, &cni.Calico{})
		}
	}

	e.runTask(conf, installerTasks...)
}
