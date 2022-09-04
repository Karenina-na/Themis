package util

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	Debug logLevel = iota
	Info
	Warn
	Error
)

type logLevel int

var f func(r any)
var flag bool

func LoggerInit(f func(r any), F logLevel) {
	switch F {
	case Debug:
		flag = true
	case Info, Warn, Error:
		flag = false
	}
	setExceptionFunc(f)
	if !exists("./log") {
		_ = os.Mkdir("./log", 0644)
	}
}

func Loglevel(level logLevel, name string, message string) {
	var logger *log.Logger
	defer func() {
		r := recover()
		if r != nil {
			f(r)
		}
	}()
	switch level {
	case Debug:
		logger = log.New(os.Stdout, name+" == "+" ["+"debug"+"] ", log.Ldate|log.Ltime)
	case Info:
		logger = log.New(os.Stdout, name+" == "+" ["+"info"+"] ", log.Ldate|log.Ltime)
	case Warn:
		logger = log.New(os.Stdout, name+" == "+" ["+"warn"+"] ", log.Ldate|log.Ltime)
	case Error:
		logger = log.New(os.Stdout, name+" == "+" ["+"error"+"] ", log.Ldate|log.Ltime)
	}
	switch level {
	case Debug:
		if flag {
			logger.Println(message)
			recordFile(message, level, logger)
		}
	case Info, Warn, Error:
		logger.Println(message)
		recordFile(message, level, logger)
	default:
		log.Panic("无此选项")
	}
}

func setExceptionFunc(exceptionFunc func(r any)) {
	f = exceptionFunc
}

func recordFile(message string, level logLevel, logger *log.Logger) {
	var FileLevel string
	switch level {
	case Debug:
		FileLevel = "debug"
	case Info:
		FileLevel = "Info"
	case Warn:
		FileLevel = "Warn"
	case Error:
		FileLevel = "Error"
	}
	year, month, day := time.Now().Date()
	t := strconv.Itoa(year) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(day)
	filename := "./log/" + t
	if !exists(filename) {
		_ = os.Mkdir(filename, 0644)
	}
	file, err := os.OpenFile(filename+"/"+FileLevel+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		f("打开日志文件异常")
	}
	logger.SetOutput(file)
	logger.Println(message)
	err = file.Close()
	if err != nil {
		f("关闭日志文件异常")
		return
	}
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
