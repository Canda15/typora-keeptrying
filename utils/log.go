package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getLogFileName() string {
	currentTime := time.Now()
	// yyyy-MM-dd_HH-mm-ss.log
	logFileName := fmt.Sprintf("log_%s.log", currentTime.Format("2006-01-02_15-04-05"))
	return logFileName
}
func Log(v ...interface{}) {
	var logDir string
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] %s\n", timestamp, fmt.Sprint(v...))
	//未初始化时不保存日志记录
	if CFG.InstallDirPath != "" {
		logDir = CFG.InstallDirPath + "\\logs"

		_ = os.MkdirAll(logDir, 0755)
		logFileName := getLogFileName()
		f, err := os.OpenFile(filepath.Join(logDir, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		f.WriteString(msg)
	}
	fmt.Println(msg)
}
