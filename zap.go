package zap

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
/  uber-zap provides a logger to replace the default Logger and SugaredLogger
/  It uses zapcore.NewCore() to create a "Core", which is a minimal, faster logger
/  interface. In this example the only logger option set is 'zap.AddCaller()' which configures
/  the logger to annotate each message with the filename, line number, and func name
/  This config will print logs in this format:
/  '| 2021-04-08T18:26:05Z INFO example/main.go:15 Starting the server...'
*/

func InitLogger() {
	createDirectoryIfNotExists()
	writerSync := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
}

func createDirectoryIfNotExists() {
	path, _ := os.Getwd()

	if _, err := os.Stat(fmt.Sprintf("%s/logs", path)); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(path+"/logs/filename.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	})
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
