package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github/lovexy-fun/skm/storage"
	"os"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&name, "name", "n", "", "key name")
	addCmd.MarkFlagRequired("name")
	addCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "key filepath")
	addCmd.MarkFlagRequired("filepath")
}

var name string
var filepath string
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a key to manager",
	Long:  "Add a key to manager",
	Run: func(cmd *cobra.Command, args []string) {
		err := storage.Add(name, filepath)
		if err == nil {
			fmt.Println("Success")
		} else if os.IsNotExist(err) {
			fmt.Println("File does not exist")
			os.Exit(1)
		}
	},
}
