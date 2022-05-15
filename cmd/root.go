package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github/lovexy-fun/skm/storage"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skm",
	Short: "SSH key manager",
	Long:  "SSH key manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Effective key: %s\n", storage.Currennt().Name)
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

//选择key
func selectKey(list []storage.Key) (int, storage.Key) {
	if len(list) == 0 {
		fmt.Println("No key")
		os.Exit(0)
	}
	prompt := promptui.Select{
		Label: "Please select key",
		Items: list,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "> {{ .Name | red }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "{{ .Name | green }}",
		},
	}
	i, _, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}
	return i, list[i]
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
