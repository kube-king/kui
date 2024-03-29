package cmd

import (
	"errors"
	"github.com/spf13/cobra"
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

	rootCmd.AddCommand(InitCmd, AddNodeCmd, GenCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

}
