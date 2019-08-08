package gorc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LstdFile   = "log"
	StdLogTime = "2006/01/02 15:04:05"
	StdLogFile = "2006-01-02"
)

type (
	// Logger represents logging interface
	Logger interface {
		Log(msg string, params map[string]interface{}, err error)
		Fatal(v ...interface{})
		Fatalf(format string, v ...interface{})
		GetOutputFile() string
		GetTimeLocation() *time.Location
	}

	// Log represents custom-zap logger
	Log struct {
		logFile string
		config  *zap.Config
		sugar   *zap.SugaredLogger
		time    *time.Location
	}
	// LogOptions represent option to custom-zap logger
	LogOptions struct {
		Development bool
		OutputFile  []string
	}
)

// New instantiates new custom-zap logger
func New(dir string, prefix string, opt *LogOptions) *Log {
	return newLog(dir, prefix, opt)
}

// newLog return new custom-zap logger
// set default logFile to yyyy-mm-dd.log
func newLog(dir string, prefix string, opt *LogOptions) *Log {
	var (
		cfg             zap.Config
		timeLocation, _ = time.LoadLocation("Asia/Jakarta")
	)

	logFile := createLogFile(create(dir), prefix, timeLocation)

	cfg = newConfig(opt.Development, logFile, timeLocation)
	if opt.Development {
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
		cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, "stdout")
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Log{
		logFile: logFile,
		config:  &cfg,
		sugar:   logger.Sugar(),
		time:    timeLocation,
	}
}

// create set log's directory location and,
// create directory if not exist
// default is 				: __path__/log
// returned log's directory : __path__/dir
func create(dir string) string {
	if dir == "" {
		dir = LstdFile
	}
	_, currentFile, _, _ := runtime.Caller(0)
	dir = filepath.Join(filepath.Dir(currentFile), dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("create folder in: %s\n", dir)
		if err = os.MkdirAll(dir, os.ModePerm|os.ModeAppend); err != nil {
			log.Fatalf("[log] failed to create directory: %v", err)
		}
	}

	return dir
}

// createLogFile set log filename
// if there's prefix,
// logfile will be dir/prefix-yyyy-mm-dd.log, else
// logfile will be dir/yyyy-mm-dd.log
func createLogFile(dir, prefix string, localTime *time.Location) string {
	var logFile = Join(time.Now().In(localTime).Format(StdLogFile), ".", LstdFile)
	if prefix != "" {
		logFile = Join(prefix, "-", logFile)
	}
	return filepath.Join(dir, logFile)
}

// newConfig set config for custom-zap logger
// set log's file to filename
// set log's time with timeLocation
func newConfig(dev bool, filename string, localTime *time.Location) (cfg zap.Config) {
	if !dev {
		cfg = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:      "json",
			EncoderConfig: zap.NewProductionEncoderConfig(),
		}
	} else {
		cfg = zap.Config{
			Level:         zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Development:   true,
			DisableCaller: true,
			Encoding:      "console",
			EncoderConfig: zapcore.EncoderConfig{
				// Keys can be anything except the empty string.
				TimeKey:    "T",
				LevelKey:   "L",
				NameKey:    "N",
				MessageKey: "M",
				// StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			},
		}
	}

	cfg.EncoderConfig.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(time.Now().In(localTime).Format(StdLogTime))
	}
	cfg.OutputPaths = []string{filename}
	cfg.ErrorOutputPaths = []string{filename}

	return cfg
}

func (l *Log) GetOutputFile() string {
	return l.logFile
}

// GetTimeLocation return time location
// Used to get time.Local in projects
func (l *Log) GetTimeLocation() *time.Location {
	return l.time
}

// Log logs using zap log.
// msg is custom message
// params contains key-value message. used for tracing
// err is error
func (l *Log) Log(msg string, params map[string]interface{}, err error) {
	var build []interface{}

	if params != nil {
		for k, v := range params {
			build = append(build, k, v)
		}
	}

	if err != nil {
		if params == nil {
			l.sugar.Error(err)
		} else {
			l.sugar.Errorw(err.Error(), build...)
		}
	} else {
		if params == nil {
			l.sugar.Info(msg)
		} else {
			l.sugar.Infow(msg, build...)
		}
	}
}

func (l *Log) Fatal(v ...interface{}) {
	l.sugar.Fatal(v...)
}

func (l *Log) Fatalf(format string, v ...interface{}) {
	l.sugar.Fatal(fmt.Sprintf(format, v...))
}