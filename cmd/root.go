package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"kube-invention/pkg/installer"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("args is empty")
		}

		return nil
	},
}

func Execute() {

	rootCmd.AddCommand(installer.InitCmd, installer.AddNodeCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
