package log

import (
	"io"
	"log"
	"os"
	"sync"
)

// 日志记录器相关变量
var (
	// errorLog 用于记录错误级别的日志，输出到标准输出，使用红色标识
	// \033[31m 是 ANSI 转义码，设置文字颜色为红色
	// log.LstdFlags 包含日期和时间的标准格式
	// log.Lshortfile 显示简短的文件名和行号信息
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)

	// infoLog 用于记录信息级别的日志，输出到标准输出，使用蓝色标识
	// \033[34m 是 ANSI 转义码，设置文字颜色为蓝色
	infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)

	// loggers 存储所有日志记录器的切片，便于统一管理和遍历
	loggers = []*log.Logger{errorLog, infoLog}

	// mu 互斥锁，用于保护并发环境下的日志写入操作，防止多个 goroutine 同时写入导致冲突
	mu sync.Mutex
)

// log methods 包级别的导出函数，用于提供简洁的日志API
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

//设置日志的层级（infoLevel,ErrorLevel,Disabled）

// log levels
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
