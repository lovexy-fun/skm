package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github/lovexy-fun/skm/storage"
)

func init() {
	rootCmd.AddCommand(delCmd)
}

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a key from manager",
	Long:  "Delete a key from manager",
	Run: func(cmd *cobra.Command, args []string) {
		idx, key := selectKey(storage.List())
		err := storage.Delete(idx, key)
		if err != nil {
			fmt.Printf("Failed to delete %s\n", key.Name)
		} else {
			fmt.Printf("%s deleted successfully\n", key.Name)
		}
	},
}
