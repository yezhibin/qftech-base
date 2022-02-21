package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Sugare export var
var Sugare *zap.SugaredLogger

// LumberJackLogger export var
var LumberJackLogger *lumberjack.Logger

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

type LogConfig struct {
	MaxCount int    // 日志文件保存最大数
	MaxSize  int    // 日志单个文件最大保存大小，单位为M
	Compress bool   // 自导打 gzip包 默认false
	FilePath string // 日志文件输出路径
	Level    LogLevel
}

// Init def
func Init(config *LogConfig) {
	if config == nil {
		config = &LogConfig{
			MaxCount: 30,
			MaxSize:  10,
			Compress: true,
			FilePath: "./log/server.log",
			Level:    InfoLevel,
		}
	}
	levelMap := map[LogLevel]zapcore.Level{
		DebugLevel: zap.DebugLevel,
		InfoLevel:  zap.InfoLevel,
	}
	LumberJackLogger = &lumberjack.Logger{
		Filename: config.FilePath, // 日志输出文件
		MaxSize:  config.MaxSize,  // 日志最大保存10M
		// MaxBackups: 5,  // 就日志保留5个备份
		MaxAge:   config.MaxCount, // 最多保留30个日志 和MaxBackups参数配置1个就可以
		Compress: config.Compress, // 自导打 gzip包 默认false
	}

	writer := zapcore.AddSync(LumberJackLogger)

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间戳的格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别使用大写显示
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writer, levelMap[config.Level])
	logger := zap.New(core, zap.AddCaller()) // 增加caller信息
	Sugare = logger.Sugar()

	Sugare.Infof("zap log init ok.")
}
