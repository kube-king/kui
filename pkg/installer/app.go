package installer

import (
	"fmt"
	"github.com/spf13/cobra"
	"kube-invention/pkg/installer/check_system"
	"kube-invention/pkg/installer/cni"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/container_runtime"
	"kube-invention/pkg/installer/download"
	"kube-invention/pkg/installer/etcd"
	"kube-invention/pkg/installer/gen"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/installer/init_system"
	"kube-invention/pkg/installer/kubernetes"
	"kube-invention/pkg/utils/yaml_util"
)

var exec = new(Exec)

var InitConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Init A Kubernetes Config File",
	Long:  "Init Kubernetes Config File",
	Run: func(cmd *cobra.Command, args []string) {
		err := gen.GenDefaultConfig()
		if err != nil {
			global.Log.Error(err.Error())
		}
	}}

var InitClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Init A Kubernetes Cluster",
	Long:  "Init Kubernetes Cluster",
	Run: func(cmd *cobra.Command, args []string) {
		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value
		var conf *config.Config
		var host config.Hosts
		err := yaml_util.UnYamlFile(configYamlFile.String(), &conf)
		if err != nil {
			return
		}
		err = yaml_util.UnYamlFile(hostYamlFile.String(), &host)
		if err != nil {
			return
		}
		conf.Hosts = host
		global.Log.Info(fmt.Sprintf("Run Inint Kubernetes Cluster Version: %v", conf.KubernetsOption.Version))
		exec.InitCluster(conf)
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "init operation",
	Long:  "init operation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

var AddNodeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Node",
	Long:  "Add Node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

var AddWorkerNodeCmd = &cobra.Command{
	Use:   "worker",
	Short: "Add Worker Node",
	Long:  "Add Worker Node",
	Run: func(cmd *cobra.Command, args []string) {

		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value
		var conf *config.Config
		var host config.Hosts
		err := yaml_util.UnYamlFile(configYamlFile.String(), &conf)
		if err != nil {
			return
		}
		err = yaml_util.UnYamlFile(hostYamlFile.String(), &host)
		if err != nil {
			return
		}
		conf.Hosts = host
		global.Log.Info(fmt.Sprintf("Run Inint Kubernetes Cluster Version: %v", conf.KubernetsOption.Version))
		exec.JoinWorkerNode(conf)
	},
}
var AddMasterNodeCmd = &cobra.Command{
	Use:   "master",
	Short: "Add Master Node",
	Long:  "Add Master Node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

func init() {
	InitClusterCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	InitClusterCmd.Flags().StringP("hosts", "s", "host.yaml", "host file")

	AddMasterNodeCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	AddMasterNodeCmd.Flags().StringP("hosts", "s", "host.yaml", "host file")

	AddWorkerNodeCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	AddWorkerNodeCmd.Flags().StringP("hosts", "s", "host.yaml", "host file")

	InitCmd.AddCommand(InitClusterCmd, InitConfigCmd)
	AddNodeCmd.AddCommand(AddMasterNodeCmd, AddWorkerNodeCmd)
}

type Installer interface {
	Exec(config *config.Config) error
}

type Exec struct {
}

func (e *Exec) runTask(config *config.Config, ins ...Installer) {
	for _, i := range ins {
		err := i.Exec(config)
		if err != nil {
			global.Log.Error(err.Error())
			break
		}
	}
}

func (e *Exec) JoinWorkerNode(conf *config.Config) {

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

	e.runTask(conf, installerTasks...)
}

func (e *Exec) InitCluster(conf *config.Config) {

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

	switch conf.CniOption.Type {
	case "calico":
		installerTasks = append(installerTasks, &cni.Calico{})
	}

	e.runTask(conf, installerTasks...)
}
