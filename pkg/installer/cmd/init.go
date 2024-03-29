package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kube-invention/pkg/installer"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/utils/yaml_util"
)

var InitCmd = &cobra.Command{
	Use:   "init",
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
		installer.InstallJobApp.InitCluster(conf)
	},
}

func init() {
	InitCmd.Flags().StringP("config", "c", "config.yaml", "config file")
	InitCmd.Flags().StringP("hosts", "s", "hosts.yaml", "host file")
}
