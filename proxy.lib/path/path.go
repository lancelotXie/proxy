/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-18 21:12:03
 * @LastEditTime: 2019-12-18 21:43:16
 */

package path

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

// ModType 模式
type ModType string

const (
	// Empty 空
	Empty ModType = ""
	// Server 服务端模式
	Server ModType = "server"
	// Client 客户端模式
	Client ModType = "client"
)

// modNow 程序当前的模式
var modNow ModType

// HomeDir 外部指定的 Home 目录，用于测试
var HomeDir string

var (
	// ErrHomeNotFound HOME 目录找不到
	ErrHomeNotFound = errors.New("home path not found")
	// ErrNeedToInitMod 还没有初始化模式
	ErrNeedToInitMod = errors.New("need to init mod firstly")
)

// Init 初始化
func Init(mod ModType) {
	modNow = mod
}

// Home 获取用户 Home 目录
func Home() (homePath string, err error) {
	if HomeDir != "" {
		homePath = HomeDir
		return
	}
	user, err := user.Current()
	if err == nil {
		homePath = user.HomeDir
		return
	}

	homePath, err = unixHome()
	if err == nil {
		return
	}

	homePath, err = windowsHome()
	return
}

func unixHome() (homePath string, err error) {
	homePath = os.Getenv("HOME")
	if homePath != "" {
		return
	}

	homePath = os.Getenv("~")
	if homePath != "" {
		return
	}

	err = ErrHomeNotFound
	return
}

func windowsHome() (homePath string, err error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	if drive != "" && path != "" {
		homePath = filepath.Join(drive, path)
		return
	}

	homePath = os.Getenv("USERPROFILE")
	if homePath == "" {
		err = ErrHomeNotFound
	}
	return
}

// ConfigDir 配置文件所在目录
func ConfigDir() (dirPath string, err error) {
	const dirname = `proxys`
	home, err := Home()
	if err != nil {
		return
	}

	dirPath = filepath.Join(home, dirname)
	return
}

// ConfigFile 配置文件路径
func ConfigFile() (filePath string, err error) {
	if modNow == Empty {
		err = ErrNeedToInitMod
		return
	}

	configDir, err := ConfigDir()
	if err != nil {
		return
	}

	var filename = string(modNow) + ".json"
	filePath = filepath.Join(configDir, filename)
	return
}
