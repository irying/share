package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

var RequestLog,ErrorLog *log.Logger

func init() {
	filePath := getLogFileFullPath()
	errFilePath := getErrorLogFileFullPath()
	file := openLogFile(filePath)
	errFile := openLogFile(errFilePath)

	RequestLog = log.New(file,"",log.Ldate | log.Ltime | log.Lshortfile)
	ErrorLog = log.New(errFile,"",log.Ldate | log.Ltime | log.Lshortfile)
}

// 日志 打印
func LogPrint(format string, values ...interface{}) {
	RequestLog.Printf("[RPC-debug] "+format, values...)
}

func LogPanic(format string, values ...interface{}) {
	ErrorLog.Printf("[RPC-debug] "+format, values...)
}
// 错误 输出到文件
func LogPrintError(err error) {

	if err != nil {
		ErrorLog.Printf("[ERROR] %v\n", err)
	}
}

// 错误 打印&&退出
func LogFatalfError(err error) {

	if err != nil {
		ErrorLog.Fatalf("[FATAL] %v\n", err)
	}
}

// 错误 转异常抛出
func ErrToPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func getLogFilePath() string {
	dir, _ := os.Getwd()
	path := dir + "/" + LogSavePath

	return fmt.Sprintf("%s", path)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func getErrorLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("errors.%s", LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// todo openfile
func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)

	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to open file :%v", err)
	}

	return handle
}

func mkDir(filePath string) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}


