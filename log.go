package infra

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ronannnn/infra/cfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLog(cfg *cfg.Log) (log *zap.SugaredLogger, err error) {
	getOrDefault(cfg)
	var level zapcore.Level
	if level, err = zapcore.ParseLevel(cfg.Level); err != nil {
		return
	}
	var writeSyncer zapcore.WriteSyncer
	if writeSyncer, err = newWriteSyncer(cfg); err != nil {
		return
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(cfg.TimeFormat))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		writeSyncer,
		level,
	)
	log = zap.New(core, zap.AddCaller()).Sugar()
	return
}

// newWriteSyncer get multiple write syncers
// 1. stdout if LogInConsole is enabled
// 2. RotateLogs if LogInRotateFile is enabled
func newWriteSyncer(cfg *cfg.Log) (syncer zapcore.WriteSyncer, err error) {
	var multiWriter []zapcore.WriteSyncer
	if cfg.LogInConsole {
		multiWriter = append(multiWriter, zapcore.AddSync(os.Stdout))
	}
	if cfg.LogInRotateFile {
		// create directory for storing log files
		if err = createDirsIfNotExist(cfg.StoreDir); err != nil {
			return
		}
		var fileWriter = &lumberjack.Logger{
			Filename:   filepath.Join(cfg.StoreDir, cfg.LatestFilename),
			MaxSize:    1, // rotate when the size gets 1MB
			MaxBackups: 0, // 0 backup: keep all old files
			MaxAge:     0, // 0 days: keep all old files
		}
		go func() {
			for {
				now := time.Now()
				// 计算距离明天该时间的时间间隔
				nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
				// utc-4小时就是我们的4点（我们是东八区）
				nextDay = nextDay.Add(-4 * time.Hour)
				duration := nextDay.Sub(now)
				// 使用定时器，在指定的时间间隔后执行函数
				timer := time.NewTimer(duration)
				<-timer.C // 等待定时器到期
				fileWriter.Rotate()
			}
		}()
		multiWriter = append(multiWriter, zapcore.AddSync(fileWriter))
	}
	return zapcore.NewMultiWriteSyncer(multiWriter...), nil
}

func getOrDefault(cfg *cfg.Log) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.StoreDir == "" {
		cfg.StoreDir = "logs"
	}
	if cfg.LatestFilename == "" {
		cfg.LatestFilename = "latest.log"
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = "2006-01-02 15:04:05.000"
	}
}
