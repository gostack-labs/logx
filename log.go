package logx

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logger struct {
	zap *zap.Logger
}

type options struct {
	Level    LogLevel
	fileOpts *LogFileOptions
}

type LogFileOptions struct {
	Filename     string
	MaxSize      int // MB
	MaxAge       int // max days
	MaxBackups   int // max files
	IsJSONFormat bool
}

type Option func(o *options)

//WithLevel options with LogLevel
func WithLevel(l LogLevel) Option {
	return func(o *options) {
		o.Level = l
	}
}

//WithFile options with write log to files.
//maxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.
//maxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename. Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc. The default is not to remove old log files based on age.
//maxBackups is the maximum number of old log files to retain. The default is to retain all old log files (though MaxAge may still cause them to get deleted.)
func WithFile(opt LogFileOptions) Option {
	return func(o *options) {
		o.fileOpts = &opt
	}
}

func NewLogger(opts ...Option) Logger {
	opt := &options{
		Level: LEVEL_INFO,
	}
	for _, o := range opts {
		o(opt)
	}
	encodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	var writer zapcore.WriteSyncer
	var cores []zapcore.Core
	//console
	{
		writer := zapcore.Lock(os.Stdout)
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = encodeTime
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, zapcore.Level(opt.Level))
		cores = append(cores, core)
	}
	//file
	if opt.fileOpts != nil && len(opt.fileOpts.Filename) > 0 {
		writer = zapcore.AddSync(lumberjackLogger(opt.fileOpts))
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = encodeTime
		var encoder zapcore.Encoder
		if opt.fileOpts.IsJSONFormat {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
		core := zapcore.NewCore(encoder, writer, zapcore.Level(opt.Level))
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)
	l := zap.New(combinedCore,
		zap.AddCallerSkip(1),
		zap.AddCaller())
	return &logger{
		zap: l,
	}
}

func lumberjackLogger(opt *LogFileOptions) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxAge:     opt.MaxAge,
		MaxBackups: opt.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}
}

func (l *logger) SetLevel(level string) {
	l.zap.Core().Enabled(zapcore.Level(ParseLogLevel(level)))
}

func (l *logger) Debug(v ...interface{}) {
	l.zap.Debug(fmt.Sprint(v...))
}

func (l *logger) Info(v ...interface{}) {
	l.zap.Info(fmt.Sprint(v...))
}

func (l *logger) Warn(v ...interface{}) {
	l.zap.Warn(fmt.Sprint(v...))
}

func (l *logger) Error(v ...interface{}) {
	l.zap.Error(fmt.Sprint(v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.zap.Info(fmt.Sprintf(format, v...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.zap.Warn(fmt.Sprintf(format, v...))
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.zap.Error(fmt.Sprintf(format, v...))
}
