package logs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	writer   io.Writer = os.Stdout // 默认输出到控制台
	file     *os.File
	mu       sync.Mutex
	initDone bool
)

// Init 初始化日志系统
func Init(logFilePath string) error {
	mu.Lock()
	defer mu.Unlock()

	if initDone {
		return nil // 已经初始化过
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	writer = io.MultiWriter(os.Stdout, f)
	file = f
	initDone = true

	return nil
}

// Close 关闭日志文件
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if file != nil {
		err := file.Close()
		file = nil
		initDone = false
		writer = os.Stdout
		return err
	}
	return nil
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	output("info ", format, v...)
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	output("error", format, v...)
}

// output 内部日志输出函数
func output(level, format string, v ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)
	logEntry := fmt.Sprintf("%s [%s] %s\n", now, level, msg)

	mu.Lock()
	defer mu.Unlock()

	_, _ = writer.Write([]byte(logEntry))
}
