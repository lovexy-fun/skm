package storage

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"os"
)

type storageJson struct {
	Current Key
	Keys    []Key
}

type Key struct {
	Id   string
	Name string
}

var userHome string
var skmHomePath string
var storagePath string
var keysPath string
var sj storageJson

func init() {
	//获取用户程序数据目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get user home")
		os.Exit(1)
	}
	userHome = homeDir
	skmHomePath = userHome + "/.skm"
	storagePath = skmHomePath + "/storage.json"
	keysPath = skmHomePath + "/keys"

	//打开用户文件，不存在用户文件就进行创建
	jsonFile, err := os.Open(storagePath)
	defer jsonFile.Close()
	if os.IsNotExist(err) {
		os.Mkdir(skmHomePath, 0766)
		os.Mkdir(keysPath, 0766)
		jsonFile, err = os.Create(storagePath)
		if err != nil {
			fmt.Printf("Failed to create %s\n", storagePath)
			os.Exit(1)
		} else {
			ioutil.WriteFile(storagePath, []byte("{}"), 0644)
		}
	}
	sj = read(jsonFile)
}

//读取文件
func read(file *os.File) storageJson {
	//读取storage.json文件字节
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to load storage.json")
		os.Exit(1)
	}

	//解析json
	var data storageJson
	if json.Unmarshal(bytes, &data) != nil {
		fmt.Println("Failed to parse storage.json")
		os.Exit(1)
	}

	return data
}

//写入文件
func write() {
	bytes, err := json.Marshal(sj)
	if err != nil {
		fmt.Println("Failed to convert data to json")
		os.Exit(1)
	}
	err = ioutil.WriteFile(storagePath, bytes, 0644)
	if err != nil {
		fmt.Println("Failed to save")
		os.Exit(1)
	}
}

//应用key
func Apply(key Key) {
	src := keysPath + "/" + key.Id
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0600)
	defer srcFile.Close()
	if os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", src)
		os.Exit(1)
	}

	//先全部读出
	bytes, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("Failed to apply key")
		os.Exit(1)
	}

	//再写到目标文件
	dest := userHome + "/.ssh/id_rsa"
	err = ioutil.WriteFile(dest, bytes, 0600)
	if err != nil {
		fmt.Println("Failed to apply key")
		os.Exit(1)
	}

	sj.Current = key

	write()
}

//获取当前生效的key
func Currennt() Key {
	return sj.Current
}

//获取key列表
func List() []Key {
	return sj.Keys
}

//添加key
func Add(name string, filepath string) error {

	src, err := os.Open(filepath)
	if err != nil {
		return err
	} else if os.IsNotExist(err) {
		return err
	}

	id := uuid.New().String()
	dest, err := os.Create(keysPath + "/" + id)
	if err != nil {
		fmt.Println("Failed to add key")
		os.Exit(1)
	}
	io.Copy(dest, src)
	dest.Close()
	src.Close()

	var k Key
	k.Name = name
	k.Id = id
	sj.Keys = append(sj.Keys, k)

	write()

	return nil
}

//删除key
func Delete(index int, key Key) error {
	sj.Keys = append(sj.Keys[:index], sj.Keys[index+1:]...)
	err := os.Remove(keysPath + "/" + key.Id)
	write()
	return err
}
