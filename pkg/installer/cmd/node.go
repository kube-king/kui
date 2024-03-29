package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kube-invention/pkg/installer"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/utils/yaml_util"
)

var AddNodeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Node",
	Long:  "Add Node",
}

var AddWorkerNodeCmd = &cobra.Command{
	Use:   "worker",
	Short: "Add Worker Node",
	Long:  "Add Worker Node",
	Args: func(cmd *cobra.Command, args []string) error {
		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value

		if configYamlFile == nil {
			return errors.New("configYamlFile is empty")
		}
		if hostYamlFile == nil {
			return errors.New("hostYamlFile is empty")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value
		var conf *config.Config
		var host config.Hosts
		err := yaml_util.UnYamlFile(configYamlFile.String(), &conf)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		err = yaml_util.UnYamlFile(hostYamlFile.String(), &host)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		conf.Hosts = host
		global.Log.Info(fmt.Sprintf("Run Inint Kubernetes Cluster Version: %v", conf.KubernetsOption.Version))
		installer.InstallJobApp.JoinWorkerNode(conf)
	},
}

var AddMasterNodeCmd = &cobra.Command{
	Use:   "master",
	Short: "Add Master Node",
	Long:  "Add Master Node",
	Args: func(cmd *cobra.Command, args []string) error {
		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value

		if configYamlFile == nil {
			return errors.New("configYamlFile is empty")
		}
		if hostYamlFile == nil {
			return errors.New("hostYamlFile is empty")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		configYamlFile := cmd.Flag("config").Value
		hostYamlFile := cmd.Flag("hosts").Value
		var conf *config.Config
		var host config.Hosts
		err := yaml_util.UnYamlFile(configYamlFile.String(), &conf)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		err = yaml_util.UnYamlFile(hostYamlFile.String(), &host)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		conf.Hosts = host
		global.Log.Info(fmt.Sprintf("Run Inint Kubernetes Cluster Version: %v", conf.KubernetsOption.Version))
		installer.InstallJobApp.JoinMasterNode(conf)
	},
}

func init() {
	AddMasterNodeCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	AddMasterNodeCmd.Flags().StringP("hosts", "s", "host.yaml", "host file")

	AddWorkerNodeCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	AddWorkerNodeCmd.Flags().StringP("hosts", "s", "host.yaml", "host file")

	AddNodeCmd.AddCommand(AddMasterNodeCmd, AddWorkerNodeCmd)
}
