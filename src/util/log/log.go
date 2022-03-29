package log

import (
	"fmt"
	"github.com/assimon/luuu/config"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Sugar *zap.SugaredLogger

func Init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	Sugar = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file := fmt.Sprintf("%s/log_%s.log",
		config.LogSavePath,
		time.Now().Format("20060102"))
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    viper.GetInt("log_max_size"),
		MaxBackups: viper.GetInt("max_backups"),
		MaxAge:     viper.GetInt("log_max_age"),
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
