package skm

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

// exitItem 退出项
var exitItem = Key{
	Id:            "E",
	Name:          "Exit",
	Description:   "Do nothing and exit. ",
	PublicKeyFile: []byte("Exit"),
}

var noKey = "No key!"

// initializeCommands 初始化命令行命令
func initializeCommands() []*cli.Command {
	commands := []*cli.Command{
		{
			Name:    "current",
			Aliases: []string{"cur"},
			Usage:   "Show current key name",
			Action:  current,
		},
		{
			Name:   "cat",
			Usage:  "Output public key file contents to the console",
			Action: cat,
		},
		{
			Name:    "use",
			Aliases: []string{"u"},
			Usage:   "Use a key",
			Action:  use,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a key",
			Action:  add,
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a key",
			Action:  create,
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete a key",
			Action:  delete,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all the keys",
			Action:  list,
		},
		{
			Name:   "dav",
			Usage:  "WebDav Setting",
			Action: dav,
		},
	}
	return commands
}

// current 输出当前生效的密钥名称
func current(cCtx *cli.Context) error {
	key, err := getCurrentKey()
	if err != nil {
		return err
	}

	var out string
	if key.Name == "" && key.Description == "" {
		out = noKey
	} else {
		out = fmt.Sprintf("Current: %s(%s)", key.Name, key.Description)
	}
	fmt.Println(promptui.Styler(promptui.FGGreen)(out))

	return nil
}

// cat 输出选中的密钥的内容
func cat(cCtx *cli.Context) error {

	//获取所有密钥
	keys, err := getAllKeys()
	if err != nil {
		return err
	}

	//添加退出选项
	keys = append(keys, exitItem)

	promptui.FuncMap["pubKey2string"] = func(bytes []byte) string {
		if len(bytes) == 4 {
			return string(bytes)
		} else {
			return string(bytes[:len(bytes)-1])
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | blue }}",
		Active:   `{{ ">" | green }}{{ .Name | green }}`,
		Inactive: " {{ .Name | blue }}",
		Selected: "{{ .PublicKeyFile | pubKey2string | green }}",
		Details:  `{{ "Description: " | blue }}{{ .Description | green }}`,
	}

	prompt := &promptui.Select{
		Label:     "Please select a key:",
		Templates: templates,
		Items:     keys,
		HideHelp:  true,
	}

	_, _, err = prompt.Run()
	if err != nil {
		return err
	}

	return nil

}

// use 应用一个密钥
func use(cCtx *cli.Context) error {

	//获取所有密钥
	keys, err := getAllKeys()
	if err != nil {
		return err
	}

	//添加退出选项
	keys = append(keys, exitItem)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | blue }}",
		Active:   `{{ ">" | green }}{{ .Name | green }}`,
		Inactive: " {{ .Name | blue }}",
		Selected: `{{ if ne .Id "E" }}{{ "Used" | green }} {{ end }}{{ .Name | green }}`,
		Details:  `{{ "Description: " | blue }}{{ .Description | green }}`,
	}

	prompt := &promptui.Select{
		Label:     "Please select a key:",
		Templates: templates,
		Items:     keys,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	key := keys[i]
	if key.Id == "E" {
		return nil
	}

	for _, algorithm := range algorithms {
		os.Remove(filepath.Join(sshHome, fmt.Sprintf("id_%s", algorithm)))
		os.Remove(filepath.Join(sshHome, fmt.Sprintf("id_%s.pub", algorithm)))
	}

	idPath := filepath.Join(sshHome, fmt.Sprintf("id_%s", key.Algorithm))
	idPubPath := filepath.Join(sshHome, fmt.Sprintf("id_%s.pub", key.Algorithm))
	privateKeyFile, err := os.OpenFile(idPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New(`Failed to write file "id_rsa". `)
	}
	defer privateKeyFile.Close()
	publicKeyFile, err := os.OpenFile(idPubPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New(`Failed to write file "id_rsa.pub". `)
	}
	defer publicKeyFile.Close()

	privateKeyFile.Write(key.PrivateKeyFile)
	publicKeyFile.Write(key.PublicKeyFile)

	err = setCurrentKey(key)
	if err != nil {
		return err
	}

	return nil
}

// add 添加一个密钥
func add(cCtx *cli.Context) error {
	// TODO
	fmt.Println(promptui.Styler(promptui.FGYellow)("TODO"))
	return nil
}

// create 创建一个密钥
func create(cCtx *cli.Context) error {
	selectTemplates := &promptui.SelectTemplates{
		Label:    "{{ . | blue }}",
		Active:   `{{ ">" | green }}{{ . | green }}`,
		Inactive: " {{ . | blue }}",
	}
	promptTemplates := &promptui.PromptTemplates{
		Valid:   "{{ . | blue }} ",
		Invalid: "{{ . | blue }} ",
	}

	//选择加密算法
	step1 := promptui.Select{
		Label:        "Please select an algorithm. ",
		Items:        algorithms,
		Templates:    selectTemplates,
		HideHelp:     true,
		HideSelected: true,
	}
	_, alg, err := step1.Run()
	if err != nil {
		return err
	}

	//输入名称
	step2 := promptui.Prompt{
		Label:       "Please enter key name:",
		Templates:   promptTemplates,
		HideEntered: true,
		Validate: func(input string) error {
			if input == "" {
				return errors.New("Not empty. ")
			}
			return nil
		},
	}
	name, err := step2.Run()
	if err != nil {
		return err
	}

	//输入描述
	step3 := promptui.Prompt{
		Label:       "Please enter key description:",
		Templates:   promptTemplates,
		HideEntered: true,
	}
	desc, err := step3.Run()
	if err != nil {
		return err
	}

	privateKey, publicKey, err := genKey(alg)
	if err != nil {
		return err
	}

	key := Key{
		Id:             uuid.New().String(),
		Name:           name,
		Description:    desc,
		Algorithm:      alg,
		PublicKeyFile:  publicKey,
		PrivateKeyFile: privateKey,
	}

	err = save(key)
	if err != nil {
		return err
	}

	fmt.Println(promptui.Styler(promptui.FGGreen)("Created"))

	return nil
}

// delete 删除一个密钥
func delete(cCtx *cli.Context) error {
	//获取所有密钥
	keys, err := getAllKeys()
	if err != nil {
		return err
	}

	//添加退出选项
	keys = append(keys, exitItem)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | blue }}",
		Active:   `{{ ">" | green }}{{ .Name | green }}`,
		Inactive: " {{ .Name | blue }}",
		Details:  `{{ "Description: " | blue }}{{ .Description | green }}`,
	}

	prompt := &promptui.Select{
		Label:        "Please select a key:",
		Templates:    templates,
		Items:        keys,
		HideHelp:     true,
		HideSelected: true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	key := keys[i]
	if key.Id == "E" {
		return nil
	}

	confirm := promptui.Select{
		Label:        fmt.Sprintf("Delete [%s]?", key.Name),
		Items:        []string{"Yes", "No"},
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | blue }}",
			Active:   `{{ ">" | green }}{{ . | green }}`,
			Inactive: " {{ . | blue }}",
		},
	}
	_, cr, err := confirm.Run()
	if err != nil {
		return err
	}

	if cr == "Yes" {
		deleteById(key.Id)
	}

	return nil
}

// list 列出所有密钥
func list(cCtx *cli.Context) error {
	keys, err := getAllKeys()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		fmt.Print(promptui.Styler(promptui.FGGreen)(noKey))
		return nil
	}

	for i, key := range keys {
		fmt.Print(promptui.Styler(promptui.FGGreen)(fmt.Sprintf("%d. %s(%s)\n", i+1, key.Name, key.Description)))
	}

	return nil
}

// dav WebDav设置
func dav(cCtx *cli.Context) error {
	// TODO
	fmt.Println(promptui.Styler(promptui.FGYellow)("TODO"))
	return nil
}
