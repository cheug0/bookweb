// logger.go
// 日志工具
// 全局统一的日志记录器，支持文件轮换和清理
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

var levelFromString = map[string]LogLevel{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"error": ERROR,
}

// Logger 统一日志器
type Logger struct {
	mu         sync.Mutex
	level      LogLevel
	output     io.Writer
	file       *os.File
	filePath   string
	enableHTTP bool
	maxSize    int64  // bytes
	maxAge     int    // days
	outputMode string // stdout, file, both
}

var (
	defaultLogger *Logger
	loggerOnce    sync.Once
)

// InitLogger 初始化全局日志器
func InitLogger(level string, output string, filePath string, enableHTTP bool, maxSizeMB int, maxAgeDays int) error {
	loggerOnce.Do(func() {
		defaultLogger = &Logger{}
	})

	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()

	// 设置日志级别
	if lvl, ok := levelFromString[level]; ok {
		defaultLogger.level = lvl
	} else {
		defaultLogger.level = INFO
	}

	defaultLogger.enableHTTP = enableHTTP
	defaultLogger.filePath = filePath
	defaultLogger.outputMode = output

	if maxSizeMB <= 0 {
		maxSizeMB = 10 // 默认 10MB
	}
	defaultLogger.maxSize = int64(maxSizeMB) * 1024 * 1024

	if maxAgeDays <= 0 {
		maxAgeDays = 7 // 默认 7天
	}
	defaultLogger.maxAge = maxAgeDays

	// 如果有旧的文件句柄，先关闭
	if defaultLogger.file != nil {
		defaultLogger.file.Close()
		defaultLogger.file = nil
	}

	// 设置输出目标
	if err := defaultLogger.setupOutput(); err != nil {
		return err
	}

	// 启动清理旧日志协程 (只启动一次)
	if filePath != "" {
		go defaultLogger.cleanOldLogs()
	}

	return nil
}

// setupOutput 根据模式设置输出
func (l *Logger) setupOutput() error {
	switch l.outputMode {
	case "file":
		if l.filePath != "" {
			if err := l.setupFile(); err != nil {
				return err
			}
			l.output = l.file
		} else {
			l.output = os.Stdout
		}
	case "both":
		if l.filePath != "" {
			if err := l.setupFile(); err != nil {
				return err
			}
			l.output = io.MultiWriter(os.Stdout, l.file)
		} else {
			l.output = os.Stdout
		}
	default: // stdout
		l.output = os.Stdout
	}
	return nil
}

// setupFile 打开日志文件
func (l *Logger) setupFile() error {
	dir := filepath.Dir(l.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.file = f
	return nil
}

// GetLogger 获取全局日志器
func GetLogger() *Logger {
	if defaultLogger == nil {
		// 默认初始化
		InitLogger("info", "stdout", "", true, 10, 7)
	}
	return defaultLogger
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if lvl, ok := levelFromString[level]; ok {
		l.level = lvl
	}
}

// SetHTTPEnabled 设置HTTP日志开关
func (l *Logger) SetHTTPEnabled(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enableHTTP = enabled
}

// IsHTTPEnabled 检查HTTP日志是否启用
func (l *Logger) IsHTTPEnabled() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.enableHTTP
}

// cleanOldLogs 清理旧日志
func (l *Logger) cleanOldLogs() {
	if l.filePath == "" || l.maxAge <= 0 {
		return
	}
	dir := filepath.Dir(l.filePath)
	base := filepath.Base(l.filePath)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// 检查文件名是否匹配模式: base.YYYYMMDD-HHmmss
		// 简单检查前缀
		if filepath.Base(path) == base {
			return nil // 当前文件不删
		}
		// 也是为了简单，只要前缀匹配且长度更长，就认为是旧日志
		if len(filepath.Base(path)) > len(base) && filepath.Base(path)[0:len(base)] == base {
			// 检查修改时间
			if time.Since(info.ModTime()) > time.Duration(l.maxAge)*24*time.Hour {
				os.Remove(path)
				fmt.Printf("Cleaned old log: %s\n", path)
			}
		}
		return nil
	})
}

// log 通用日志方法
func (l *Logger) log(level LogLevel, module string, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	// 检查轮换 (仅当输出到文件时)
	if l.file != nil {
		if err := l.checkAndRotate(); err != nil {
			fmt.Printf("Log rotate error: %v\n", err)
		}
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	var logLine string
	if module != "" {
		logLine = fmt.Sprintf("[%s] %s [%s] %s\n", levelNames[level], timestamp, module, msg)
	} else {
		logLine = fmt.Sprintf("[%s] %s %s\n", levelNames[level], timestamp, msg)
	}

	if l.output != nil {
		l.output.Write([]byte(logLine))
	}
}

// checkAndRotate 检查并执行轮换
func (l *Logger) checkAndRotate() error {
	if l.file == nil {
		return nil
	}
	info, err := l.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= l.maxSize {
		// 1. 关闭当前文件
		l.file.Close()
		l.file = nil

		// 2. 重命名
		timestamp := time.Now().Format("20060102-150405")
		newName := fmt.Sprintf("%s.%s", l.filePath, timestamp)
		os.Rename(l.filePath, newName)

		// 3. 重新建立输出 (重新打开文件并根据模式设置 writer)
		if err := l.setupOutput(); err != nil {
			// 失败降级
			l.output = os.Stdout
			return err
		}

		// 4. 触发一次清理 (可选)
		go l.cleanOldLogs()
	}
	return nil
}

// Debug 调试日志
func (l *Logger) Debug(module string, format string, args ...interface{}) {
	l.log(DEBUG, module, format, args...)
}

// Info 信息日志
func (l *Logger) Info(module string, format string, args ...interface{}) {
	l.log(INFO, module, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(module string, format string, args ...interface{}) {
	l.log(WARN, module, format, args...)
}

// Error 错误日志
func (l *Logger) Error(module string, format string, args ...interface{}) {
	l.log(ERROR, module, format, args...)
}

// HTTP HTTP请求日志
func (l *Logger) HTTP(method, path string, status int, duration time.Duration, clientIP string) {
	if !l.IsHTTPEnabled() {
		return
	}
	l.log(INFO, "HTTP", "%s %s %d %v %s", method, path, status, duration, clientIP)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

// 便捷全局方法
func LogDebug(module string, format string, args ...interface{}) {
	GetLogger().Debug(module, format, args...)
}

func LogInfo(module string, format string, args ...interface{}) {
	GetLogger().Info(module, format, args...)
}

func LogWarn(module string, format string, args ...interface{}) {
	GetLogger().Warn(module, format, args...)
}

func LogError(module string, format string, args ...interface{}) {
	GetLogger().Error(module, format, args...)
}

func LogHTTP(method, path string, status int, duration time.Duration, clientIP string) {
	GetLogger().HTTP(method, path, status, duration, clientIP)
}
