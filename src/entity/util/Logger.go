package util

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	Debug = iota
	Info
	Warn
	Error
)

var f func(r any)

func LoggerInit(f func(r any)) {
	setExceptionFunc(f)
	if !exists("./log") {
		_ = os.Mkdir("./log", 0644)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Loglevel(level int, name string, message string) {
	defer func() {
		r := recover()
		if r != nil {
			f(r)
		}
	}()
	log.SetPrefix("[" + name + "] ")
	switch level {
	case Debug:
		printStdio(message)
	case Info, Warn, Error:
		printStdio(message)
		recordFile(message, level)
	default:
		log.Panic("无此选项")
	}
}

func setExceptionFunc(exceptionFunc func(r any)) {
	f = exceptionFunc
}

func printStdio(message string) {
	log.Println(message)
}

func recordFile(message string, level int) {
	var FileLevel string
	switch level {
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
	f, err := os.OpenFile(filename+"/"+FileLevel+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	log.SetOutput(f)
	log.Println(message)
	log.SetOutput(os.Stdout)
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
