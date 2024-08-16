package logger

import (
	"FreeOps/global"
	"FreeOps/internal/consts"
	"fmt"
	"log"
	"time"

	"os"
)

type Logger struct {
	level int
}

var logger *Logger

func BuildLogger(level string) {
	var logLevel int
	switch level {
	case "Error":
		logLevel = consts.LevelError
	case "Warning":
		logLevel = consts.LevelWarning
	case "Info":
		logLevel = consts.LevelInfo
	case "Debug":
		logLevel = consts.LevelDebug
	}
	logvar := Logger{
		level: logLevel,
	}
	logger = &logvar
}

func Log() *Logger {
	if logger == nil {
		logvar := Logger{
			level: consts.LevelInfo,
		}
		logger = &logvar
	}
	return logger
}

// 日志写入函数
func (logvar *Logger) Panic(service string, handler string, m ...any) {
	if consts.LevelError > logvar.level {
		return
	}
	msg := fmt.Sprint("[Panic] "+"["+handler+"] ", fmt.Sprint(m...))
	commonLog(service, msg)
	os.Exit(0)
}

func (logvar *Logger) Error(service string, handler string, m ...any) bool {
	if consts.LevelError > logvar.level {
		return true
	}
	msg := fmt.Sprint("[Error] "+"["+handler+"] ", fmt.Sprint(m...))
	if !commonLog(service, msg) {
		return false
	}
	return true
}

func (logvar *Logger) Warning(service string, handler string, m ...any) bool {
	if consts.LevelWarning > logvar.level {
		return true
	}
	msg := fmt.Sprint("[Warning] "+"["+handler+"] ", fmt.Sprint(m...))
	if !commonLog(service, msg) {
		return false
	}
	return true
}

func (logvar *Logger) Info(service string, handler string, m ...any) bool {
	if consts.LevelInfo > logvar.level {
		return true
	}
	msg := fmt.Sprint("[Info] "+"["+handler+"] ", fmt.Sprint(m...))
	if !commonLog(service, msg) {
		return false
	}
	return true
}

func (logvar *Logger) Debug(service string, handler string, m ...any) bool {
	if consts.LevelDebug > logvar.level {
		return true
	}
	msg := fmt.Sprint("[Debug] "+"["+handler+"] ", fmt.Sprint(m...))
	if !commonLog(service, msg) {
		return false
	}
	return true
}

func commonLog(service string, msg string) bool {
	var dirPath, file string
	if service == "" {
		dirPath = global.RootPath + "/logs" + "/common"
		file = dirPath + "/" + time.Now().Format("2006-01-02") + ".log"
	} else {
		dirPath = global.RootPath + "/logs/" + service
		file = dirPath + "/" + time.Now().Format("2006-01-02") + ".log"
	}

	if dirStat, err := os.Stat(global.RootPath + dirPath); err != nil || !dirStat.IsDir() {
		if err := os.MkdirAll(dirPath, 0775); err != nil {
			fmt.Println("存放日志的目录未创建成功")
			return false
		}
	}

	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		fmt.Println("打开日志文件失败")
		return false
	}
	if logFile == nil {
		return false
	}
	log.SetOutput(logFile)
	log.SetPrefix("[" + time.Now().Local().Format("2006-01-02 15:04:05") + "] ")
	log.Println(msg)
	return true
}
