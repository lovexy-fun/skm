package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github/lovexy-fun/skm/storage"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all keys",
	Long:  "List all keys",
	Run: func(cmd *cobra.Command, args []string) {
		list := storage.List()
		for i := 0; i < len(list); i++ {
			fmt.Printf("[%d] %s\n", i+1, list[i].Name)
		}
	},
}
