package skm

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Key struct {
	Id             string `gorm:"primarykey"`
	Name           string
	Description    string
	Algorithm      string
	PublicKeyFile  []byte
	PrivateKeyFile []byte
}

type CurrentKey struct {
	Id          string
	Name        string
	Description string
}

type WebDAV struct {
	Url      string
	Username string
	Password string
}

var gdb *gorm.DB

// getDB 获取数据库连接
func getDB() (*gorm.DB, error) {

	if gdb != nil {
		return gdb, nil
	}

	var err error

	gdb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to database. Reason: %s", err.Error()))
	}
	return gdb, nil
}

// createTables 创建表
func createTables(db *gorm.DB) error {
	err := db.AutoMigrate(&Key{}, &CurrentKey{}, &WebDAV{})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create table. Reason: %s", err.Error()))
	}
	return nil
}

// getAllKeys 查询所有key
func getAllKeys() ([]Key, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	var keys []Key
	db.Find(&keys)
	return keys, nil
}

// getCurrentKey 获取当前key
func getCurrentKey() (*CurrentKey, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	var currentKey CurrentKey
	db.Find(&currentKey)
	return &currentKey, nil
}

// setCurrentKey 设置当前key
func setCurrentKey(key Key) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	db.Delete(&CurrentKey{}, "1=1").Create(
		&CurrentKey{
			Id:          key.Id,
			Name:        key.Name,
			Description: key.Description,
		})
	return nil
}

func save(key Key) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	result := db.Create(&key)
	if result.Error != nil {
		return errors.New(fmt.Sprintf("Failed to save key. Reason: %s", result.Error.Error()))
	}
	return nil
}

func deleteById(id string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	result := db.Delete(Key{
		Id: id,
	})
	if result.Error != nil {
		return errors.New(fmt.Sprintf("Failed to delete key. Reason: %s", result.Error.Error()))
	}
	return nil
}
