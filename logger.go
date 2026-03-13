package ciproxy

import (
	"io"
	"log"
	"os"
	"sync"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger 结构化日志器
type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	prefix string
	logger *log.Logger
}

var defaultLogger *Logger
var loggerOnce sync.Once

// GetLogger 获取默认日志器
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		defaultLogger = &Logger{
			level:  LogLevelInfo,
			prefix: "[CiProxy]",
			logger: log.New(os.Stdout, "[CiProxy] ", log.LstdFlags|log.Lshortfile),
		}
	})
	return defaultLogger
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput 设置输出目标
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetOutput(w)
}

// Debug 输出调试日志
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+msg, args...)
	}
}

// Info 输出信息日志
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logger.Printf("[INFO] "+msg, args...)
	}
}

// Warn 输出警告日志
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logger.Printf("[WARN] "+msg, args...)
	}
}

// Error 输出错误日志
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.logger.Printf("[ERROR] "+msg, args...)
	}
}
