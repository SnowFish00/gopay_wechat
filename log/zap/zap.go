package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局log变量
var log *zap.SugaredLogger

// zapLevel声明
type Level = zapcore.Level

const (
	//对zapLevel自定义 数字编号
	DebugLevel Level = -1

	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
	PanicLevel
)

func InitZap() *zap.SugaredLogger {
	// 获取一个指定的的EncoderConfig，进行自定义
	encoderConfig := zap.NewProductionEncoderConfig()
	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//日志格式 ConsoleEncoder控制台格式输出 JSONEncoder JSON格式输出
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//自定义日志输出级别
	highlevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})
	lowlevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zap.ErrorLevel && l >= zap.DebugLevel
	})

	//lumberjack日志分割和属性自定义
	lowFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "../log/runlow.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	})
	//Core装载
	/*
	   syncFile := zapcore.AddSync(lumberJackLogger) // 打印到文件
	   syncConsole := zapcore.AddSync(os.Stderr) // 打印到控制台
	   return zapcore.NewMultiWriteSyncer(syncFile, syncConsole) //同时打印
	*/
	lowFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(lowFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowlevel)

	highFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "../log/runhigh.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	})
	highFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(highFileWriteSyncer, zapcore.AddSync(os.Stdout)), highlevel)

	//多个Core装载
	core := zapcore.NewTee(highFileCore, lowFileCore)
	//添加Caller记录
	logger := zap.New(core, zap.AddCaller())
	//Sugar日志格式
	log = logger.Sugar()

	return log
}

// 日志记录函数封装
func Debug(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func Info(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warn(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Error(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Panic(template string, args ...interface{}) {
	log.Panicf(template, args...)
}
