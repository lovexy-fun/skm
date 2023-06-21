package skm

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

var (
	userHome string
	skmHome  string
	dbPath   string
	sshHome  string
)

func Run() {
	var err error

	userHome, err = getHome()
	if err != nil {
		log.Fatalln(err.Error())
	}
	skmHome = filepath.Join(userHome, ".skm")
	dbPath = filepath.Join(skmHome, "skm.db")
	sshHome = filepath.Join(userHome, ".ssh")

	/* 初始化 */
	if !isInitialized() {
		err = initializeSkm()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	app := &cli.App{
		Name:    "skm",
		Usage:   "SSH Key Manager",
		Version: "v0.0.2",
		Action:  current,
	}
	app.Commands = initializeCommands()
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalln(err.Error())
	}

}

// getHome 获取用户目录
func getHome() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.New("Failed to get HOME. ")
	}
	return dir, nil
}

// initializeSkm 初始化skm
func initializeSkm() error {

	_, err := os.Stat(skmHome)
	if os.IsNotExist(err) {
		err = os.Mkdir(skmHome, 0644)
		if err != nil {
			return errors.New("Failed to create directory \".skm\". ")
		}
	}

	db, err := getDB()
	if err != nil {
		return err
	}
	err = createTables(db)
	if err != nil {
		return err
	}

	return nil

}

// checkInitialized 检查是否已经初始化
func isInitialized() bool {
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err)
}
