package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"kube-invention/pkg/installer/gen"
	"kube-invention/pkg/installer/global"
)

var GenConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generator A Config File",
	Long:  "Generator Kubernetes Config File",
	Args: func(cmd *cobra.Command, args []string) error {

		kubernetesVersion := cmd.Flag("kubernetes-version").Value
		containerRuntimeType := cmd.Flag("container-runtime-type").Value
		arch := cmd.Flag("arch").Value
		vip := cmd.Flag("vip").Value

		if kubernetesVersion == nil {
			return errors.New("kubernetes-version is empty")
		}

		if containerRuntimeType == nil {
			return errors.New("container-runtime-type is empty")
		}

		if arch == nil {
			return errors.New("arch is empty")
		}

		if vip == nil {
			return errors.New("vip is empty")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kubernetesVersion := cmd.Flag("kubernetes-version").Value
		containerRuntimeType := cmd.Flag("container-runtime-type").Value
		arch := cmd.Flag("arch").Value
		vip := cmd.Flag("vip").Value

		err := gen.GenDefaultConfig(kubernetesVersion.String(), arch.String(), containerRuntimeType.String(), vip.String())
		if err != nil {
			global.Log.Error(err.Error())
		}
	}}

var GenHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Generator A host File",
	Long:  "Generator host File",
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value
		err := gen.GenHostConfig(name.String())
		if err != nil {
			global.Log.Error(err.Error())
		}
	}}

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen operation",
	Long:  "gen operation",
}

func init() {

	GenConfigCmd.Flags().StringP("kubernetes-version", "", "v1.28.0", "kubernetes version")
	GenConfigCmd.Flags().StringP("container-runtime-type", "", "docker", "container runtime type (docker,containerd)")
	GenConfigCmd.Flags().StringP("arch", "", "amd64", "cpu arch (amd64,arm64)")
	GenConfigCmd.Flags().StringP("vip", "", "1.1.1.1", "vip address")

	GenHostCmd.Flags().StringP("name", "n", "hosts", "host yaml")

	GenCmd.AddCommand(GenConfigCmd, GenHostCmd)
}
