package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	req "github.com/machmum/gorc/request"
	str "github.com/machmum/gorc/stringc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LLvlDevelopment = 1 << iota
	LLvlProduction
	LStdLocation = "Asia/Jakarta"        // initial value for log's location
	LStdDir      = "log"                 // initial value for log's directory
	LStdTime     = "2006/01/02 15:04:05" // log's default datetime
	LStdFile     = "2006-01-02"          // log's default filename
)

type (
	// Log represents logger / custom zap-logger
	Logger struct {
		file  string             // filename
		time  *time.Location     // time used in logs
		sugar *zap.SugaredLogger // in contexts where performance is nice, but not critical, use the SugaredLogger.
		// In the rare contexts where every microsecond and every allocation matter,
		// use the Logger. It's even faster than the SugaredLogger and allocates far less,
		// but it only supports strongly-typed, structured logging.
		logger *zap.Logger
	}

	// LogOptions represent option to custom-zap logger
	// Level set log's level logger, either development or production
	// Time set log's time location being used,
	// default is "Asia/Jakarta". Use according to Time Zone database, such as "America/New_York".
	// WithTrace set trace-id to logs output.
	// RefID will set ref-id to logs output.
	// Output file is another output file. If you want
	// logger to write log to multiple file, add other source here.
	// e.g : if you want logger to log to file and console, add "stdout" to LogOptions.OutputFile
	LogOptions struct {
		Level      int
		Time       *time.Location
		WithTrace  bool
		RefID      string
		OutputFile []string
	}
)

// newLogger return new custom-zap logger
// set default logFile to yyyy-mm-dd.log
func (opt *LogOptions) newLogger(dir string, prefix string) *Logger {
	// default log level is Production
	if opt.Level < 1 {
		opt.Level = LLvlProduction
	}
	// default log time is "Asia/Jakarta"
	if opt.Time == nil {
		opt.Time, _ = time.LoadLocation(LStdLocation)
	}
	logFile := opt.makeLogFile(create(dir), prefix)

	logger, err := opt.newConfig(logFile).Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Logger{
		file:   logFile,
		time:   opt.Time,
		sugar:  logger.Sugar(),
		logger: logger,
	}
}

// makeLogFile set log filename
// if there's prefix,
// logfile will be dir/prefix-yyyy-mm-dd.log, else
// logfile will be dir/yyyy-mm-dd.log
func (opt *LogOptions) makeLogFile(dir, prefix string) string {
	logFile := str.StringBuilder(time.Now().In(opt.Time).Format(LStdFile), ".", LStdDir)
	if prefix != "" {
		return str.StringBuilder(dir, "/", prefix, "-", logFile)
	}
	return str.StringBuilder(dir, "/", logFile)
}

// newConfig set config for custom-zap logger
// set log's file to logFile
// set log's time with timeLocation
func (opt *LogOptions) newConfig(logFile string) (cfg zap.Config) {
	if opt.Level > LLvlDevelopment {
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

	var trace bool
	if opt.RefID != "" {
		trace = true
		if opt.WithTrace {
			cfg.InitialFields = map[string]interface{}{"trace-id": req.RequestID(), "ref-id": opt.RefID}
		} else {
			cfg.InitialFields = map[string]interface{}{"ref-id": opt.RefID}
		}
	}

	if opt.WithTrace && !trace {
		cfg.InitialFields = map[string]interface{}{"trace-id": req.RequestID()}
	}

	cfg.EncoderConfig.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(time.Now().In(opt.Time).Format(LStdTime))
	}
	cfg.OutputPaths = []string{logFile}
	cfg.ErrorOutputPaths = []string{logFile}

	if len(opt.OutputFile) > 0 {
		for _, out := range opt.OutputFile {
			cfg.OutputPaths = append(cfg.OutputPaths, out)
			cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, out)
		}
	}

	return cfg
}

// create set log's directory location and,
// create directory if not exist
// default is 				: __path__/log
// returned log's directory : __path__/dir
func create(dir string) string {
	if dir == "" {
		dir = LStdDir
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("create folder in: %s\n", dir)
		if err = os.MkdirAll(dir, os.ModePerm|os.ModeAppend); err != nil {
			log.Fatalf("[log] failed to create directory: %v", err)
		}
	}
	return dir
}

// NewLogger initiate new custom-zap logger,
// by satisfy log's directoryName, prefix and options
func NewLogger(directoryName string, prefix string, opt *LogOptions) *Logger {
	return opt.newLogger(directoryName, prefix)
}

// Log logs using zap log.
// msg is custom message
// params contains key-value message. used for tracing
// err is error
func (l *Logger) Log(msg string, params map[string]interface{}, err error) {
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

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(v ...interface{}) {
	l.sugar.Fatal(v...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.sugar.Fatal(fmt.Sprintf(format, v...))
}

// GetOutputFile returns log's filename name.
func (l *Logger) GetOutputFile() string {
	return l.file
}

// GetTimeLocation return time location used in logs
func (l *Logger) GetTimeLocation() *time.Location {
	return l.time
}
