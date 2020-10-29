/*
Log.WithFields(log.Fields{
	"name": "maApp",
	"size":   10,
}).Info("log info")
*/
package common

import (
	"fmt"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Logger = log.New()        // 创建log对象
	writer *lumberjack.Logger // 输出源
)

type PrintHook struct {
	Level []log.Level
}
type OutputHook struct {
	Level []log.Level
}

func InitLog() {
	LogSetting()
	LogPrintLine()
	LogOutput()
}

// 设置参数
func LogSetting() {
	Logger.SetFormatter(&log.TextFormatter{
		ForceColors:            true,                  // 绕过tty输出颜色
		DisableColors:          false,                 // 强制禁用颜色
		DisableLevelTruncation: false,                 // 日志级别截断为4字节,false为4字节
		FullTimestamp:          true,                  // 启用完整时间戳
		TimestampFormat:        "2006-01-02 15:04:05", // 时间戳格式
	})
	Logger.SetReportCaller(true) // 添加调用方法作为字段
	// 设置日志等级
	if Trace {
		Logger.SetLevel(log.TraceLevel)
	} else if Debug {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.InfoLevel)
	}
	Logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout)) // 设置转义，解决颜色显示问题
	//Logger.SetOutput(ioutil.Discard)     //设置不输出到控制台
	//Loggers = Logger.WithFields(logrus.Fields{ // 额外添加字段
	//"log_level": Logger.Level,
	//})
}

// 打印行号
func LogPrintLine() {
	h := &PrintHook{}
	_ = h.setLevels(log.AllLevels)
	Logger.AddHook(h)
}

func (hook *PrintHook) Fire(entry *log.Entry) error {
	f := strings.Split(filepath.Base(entry.Caller.File), ".")
	if len(f) > 1 {
		entry.Caller.File = f[0]
	} else if len(f) > 0 {
		entry.Caller.File = f[0]
	}
	s := strings.Split(filepath.Base(entry.Caller.Function), ".")
	if len(s) > 1 {
		entry.Caller.Function = s[len(s)-1]
	}
	return nil
}

func (hook *PrintHook) Levels() []log.Level {
	return hook.Level
}

// 设置触发级别
func (hook *PrintHook) setLevels(levels []log.Level, level ...log.Level) error {
	if len(levels) > 0 {
		hook.Level = levels
		return nil
	}
	for _, v := range level {
		hook.Level = append(hook.Level, v)
	}
	return nil
}

// 输出到日志文件
func LogOutput() {
	fileName := fmt.Sprintf("%s.log", filepath.Base(os.Args[0])) // 应用程序名称为日志文件名
	if filepath.IsAbs(LogPath) {                                 // 是否绝对路径
		LogPath = filepath.Join(LogPath, fileName)
	} else {
		appPath := GetCurrentDirectory()
		LogPath = filepath.Join(appPath, LogPath, fileName)
	}
	writer = &lumberjack.Logger{
		Filename:   LogPath,     // 当前日志文件名
		MaxSize:    100,         // 最大大小MB
		MaxBackups: 0,           // 保留个数，默认不删除
		MaxAge:     LogHoldDays, // 保留天数，默认不删除
		Compress:   false,       // 是否压缩
	}
	h := &OutputHook{}
	_ = h.setLevels(nil, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel)
	Logger.AddHook(h)
}

func (hook *OutputHook) Fire(entry *log.Entry) error {
	_ = readLogUpdateTime(LogPath, entry.Time)
	format := &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 时间戳格式
		PrettyPrint:     false,                 // 缩进所有json日志
	}
	msg, err := format.Format(entry)
	if err != nil {
		log.Errorf("failed to generate string for entry: %+v", err)
		return err
	}
	_, err = writer.Write(msg)
	return nil
}

func (hook *OutputHook) Levels() []log.Level {
	return hook.Level
}

// 设置触发级别
func (hook *OutputHook) setLevels(levels []log.Level, level ...log.Level) error {
	if len(levels) > 0 {
		hook.Level = levels
		return nil
	}
	for _, v := range level {
		hook.Level = append(hook.Level, v)
	}
	return nil
}

// 读取当前log文件修改时间，并判断分割
func readLogUpdateTime(name string, nowTime time.Time) error {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	fi, _ := f.Stat()
	if fi.ModTime().Year() != nowTime.Year() || fi.ModTime().Month() != nowTime.Month() || fi.ModTime().Day() != nowTime.Day() {
		go func() { _ = writer.Rotate() }()
	}
	return nil
}
