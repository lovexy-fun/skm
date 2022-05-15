package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	//rootCmd.AddCommand(genCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a key",
	Long:  "Generate a key",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command not completed")
	},
}
