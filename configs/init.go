package configs

import (
	"FreeOps/global"
	"FreeOps/pkg/logger"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Init() {
	err := declareRootPath()
	if err != nil {
		panic(fmt.Errorf("获取项目根目录失败: %v ", err))
	}
	viper.SetConfigFile(global.RootPath + "/configs/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置信息失败: %s ", err))
	}
	if err = viper.Unmarshal(global.Conf); err != nil {
		panic(fmt.Errorf("获取配置失败, err:%s ", err))
	}
	logger.BuildLogger(global.Conf.Logger.Level)
	// 监视配置文件的变化。当配置文件发生改变时，它会自动重新加载并更新 Viper 实例中的配置信息。这样可以避免在应用程序运行时手动重新加载配置文件
	viper.WatchConfig()
	// OnConfigChange，回调函数，用于在配置文件发生变化时进行处理。您可以将自定义的函数传递给 OnConfigChange()，在配置文件发生更改时，该函数将被调用
	viper.OnConfigChange(func(in fsnotify.Event) { // 传递配置文件变更事件的参数类型，以便在 OnConfigChange() 回调函数中获取有关配置文件变化的详细信息。
		logger.Log().Warning("config", "Conf", "配置文件触发修改重载"+in.Name)
		if err = viper.Unmarshal(global.Conf); err != nil {
			logger.Log().Panic("config", "Conf", "配置文件写入结构体变量失败")
		}
	})

	if err = declareGlobal(); err != nil {
		panic(fmt.Errorf("declareGlobal初始化失败: %v", err))
	}
}

// FindRootDir 递归向上查找直到找到 go.mod 文件，返回该目录作为项目根目录
func findRootDir(dir string) (string, error) {
	// 检查当前目录是否存在 go.mod 文件
	if _, err := os.Stat(filepath.Join(dir, `go.mod`)); err == nil {
		return dir, nil
	} else if !os.IsNotExist(err) {
		// 发生了其他错误
		return "", err
	}

	// 获取当前目录的父目录
	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		// 如果父目录与当前目录相同，说明已经到达根目录
		return "", os.ErrNotExist
	}

	// 递归向上查找
	return findRootDir(parentDir)
}

func declareRootPath() (err error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("No caller information")
	}
	dir := filepath.Dir(filename)
	global.RootPath, err = findRootDir(dir)
	global.RootPath = strings.Replace(global.RootPath, "\\", "/", -1)
	return err
}

func declareConcurrencyVar() {

	// 设置总并发数
	global.Sem = semaphore.NewWeighted(global.Conf.Concurrency.Number)
	global.MaxWebSSH = uint64(global.Conf.WebSSH.MaxConnNumber)
	global.WebSSHCounter = 0
}

func declareGlobal() (err error) {
	declareConcurrencyVar()
	// ssh密钥
	global.OpsSSHKey, err = os.ReadFile(global.Conf.SshConfig.OpsKeyPath)
	return err
}
