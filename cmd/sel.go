package cmd

import (
	"github.com/spf13/cobra"
	"github/lovexy-fun/skm/storage"
)

func init() {
	rootCmd.AddCommand(selCmd)
}

var selCmd = &cobra.Command{
	Use:   "sel",
	Short: "Choose a key to make it effective",
	Long:  "Choose a key to make it effective",
	Run: func(cmd *cobra.Command, args []string) {
		_, key := selectKey(storage.List())
		storage.Apply(key)
	},
}
