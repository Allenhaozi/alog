package alog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/allenhaozi/alog/internal/config"
	"github.com/allenhaozi/alog/spew"
)

// define log level
const (
	LevelCritical = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
	levelSize
)

var levels = map[string]int{
	"info":     LevelInfo,
	"trace":    LevelTrace,
	"debug":    LevelDebug,
	"warn":     LevelWarn,
	"error":    LevelError,
	"critical": LevelCritical,
}

var (
	loggers = make([]*logger, levelSize, levelSize)

	discard = log.New(ioutil.Discard, "", 0)
)

// initialize
func init() {

	for index := range loggers {
		loggers[index] = &logger{
			log: discard,
		}
	}
}

// InitFromXMLFile 从一个 XML 文件中初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLFile(path string) error {
	cfg, err := config.ParseXMLFile(path)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// InitFromXMLString 从一个 XML 字符串初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLString(xml string) error {
	cfg, err := config.ParseXMLString(xml)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

func InitALog(data map[string]string) error {
	if len(data["path"]) <= 0 {
		return errors.New("log path not found !")
	}

	if len(data["size"]) <= 0 {
		return errors.New("log rotate size not found !")
	}

	if len(data["level"]) <= 0 {
		data["level"] = "info"
	}
	if _, ok := levels[data["level"]]; !ok {
		return errors.New("invalid log level set !")
	}

	xml := config.InitConfigString(data)
	cfg, err := config.ParseXMLString(xml)
	if err != nil {
		return err
	}

	err1 := initFromConfig(cfg)
	if err1 != nil {
		return err1
	}

	//set log level
	setLogLevel(data["level"])

	return nil
}

//disable unaccepted log level
func setLogLevel(strlevel string) {
	intLevel := levels[strlevel]
	for _, level := range levels {
		if level > intLevel {
			loggers[level].set(nil, "", 0)
		}
	}

}

//
// 若将 w 设置为 nil 等同于 iotuil.Discard，即关闭此类型的输出。
func SetWriter(level int, w io.Writer, prefix string, flag int) error {
	if level < 0 || level > levelSize {
		return errors.New("无效的 level 值")
	}

	loggers[level].set(w, prefix, flag)
	return nil
}

// 从 config.Config 中初始化整个 logs 系统
func initFromConfig(cfg *config.Config) error {
	for name, c := range cfg.Items {
		index, found := levels[name]
		if !found {
			return fmt.Errorf("未知道的二级元素名称:[%s]", name)
		}
		flag, err := parseFlag(c.Attrs["flag"])
		if err != nil {
			return err
		}

		w, err := toWriter(c)
		if err != nil {
			return err
		}
		loggers[index].set(w, c.Attrs["prefix"], flag)
	}

	return nil
}

// Flush 输出所有的缓存内容。
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Flush() 输出可能缓存的日志内容。
func Flush() {
	for _, l := range loggers {
		if l.flush != nil {
			l.flush.Flush()
		}
	}
}

// INFO 获取 INFO 级别的 log.Logger 实例，在未指定 info 级别的日志时，该实例返回一个 nil。
func INFO() *log.Logger {
	return loggers[LevelInfo].log
}

// Info 相当于 INFO().Println(v...) 的简写方式
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func Info(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		INFO().Output(2, string(r))
	}
}

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func Infof(format string, v ...interface{}) {
	INFO().Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例，在未指定 debug 级别的日志时，该实例返回一个 nil。
func DEBUG() *log.Logger {
	return loggers[LevelDebug].log
}

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func Debug(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		DEBUG().Output(2, string(r))
	}
}

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func Debugf(format string, v ...interface{}) {
	DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例，在未指定 trace 级别的日志时，该实例返回一个 nil。
func TRACE() *log.Logger {
	return loggers[LevelTrace].log
}

// Trace 相当于 TRACE().Println(v...) 的简写方式
func Trace(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		TRACE().Output(2, string(r))
	}
}

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func Tracef(format string, v ...interface{}) {
	TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例，在未指定 warn 级别的日志时，该实例返回一个 nil。
func WARN() *log.Logger {
	return loggers[LevelWarn].log
}

// Warn 相当于 WARN().Println(v...) 的简写方式
func Warn(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		WARN().Output(2, string(r))
	}
}

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func Warnf(format string, v ...interface{}) {
	WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例，在未指定 error 级别的日志时，该实例返回一个 nil。
func ERROR() *log.Logger {
	return loggers[LevelError].log
}

// Error 相当于 ERROR().Println(v...) 的简写方式
func Error(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		ERROR().Output(2, string(r))
	}
}

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func Errorf(format string, v ...interface{}) {
	ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例，在未指定 critical 级别的日志时，该实例返回一个 nil。
func CRITICAL() *log.Logger {
	return loggers[LevelCritical].log
}

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) {
	r, err := json.Marshal(v)
	if err == nil {
		CRITICAL().Output(2, string(r))
	}
}

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func Criticalf(format string, v ...interface{}) {
	CRITICAL().Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容。
func All(v ...interface{}) {
	all(v...)
}

// Allf 向所有的日志输出内容。
func Allf(format string, v ...interface{}) {
	allf(format, v...)
}

// Fatal 输出错误信息，然后退出程序。
func Fatal(v ...interface{}) {
	all(v...)
	Flush()
	os.Exit(2)
}

// Fatalf 输出错误信息，然后退出程序。
func Fatalf(format string, v ...interface{}) {
	allf(format, v...)
	Flush()
	os.Exit(2)
}

// Panic 输出错误信息，然后触发 panic。
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	all(s)
	Flush()
	panic(s)
}

// Panicf 输出错误信息，然后触发 panic。
func Panicf(format string, v ...interface{}) {
	allf(format, v...)
	Flush()
	panic(fmt.Sprintf(format, v...))
}

func all(v ...interface{}) {
	for _, l := range loggers {
		l.log.Output(3, fmt.Sprintln(v...))
	}
}

func allf(format string, v ...interface{}) {
	for _, l := range loggers {
		l.log.Output(3, fmt.Sprintf(format, v...))
	}
}

func Dump(v ...interface{}) {
	spew.Dump(v)
}
func TT(v ...interface{}) {
	spew.Dump(v)
}
func Pretty(i ...interface{}) {
	for _, v := range i {
		res, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(res))
	}
}
