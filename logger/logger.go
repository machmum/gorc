package logger

import (
	"fmt"
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

var (
	filename     = "stdout"     // log's filename. Relative path + filename
	timeLocation *time.Location // log's time location
)

type (
	// Log represents logger / custom zap-logger
	// A Logger provides fast, leveled, structured logging. All methods are safe
	// for concurrent use.
	//
	// The Logger is designed for contexts in which every microsecond and every
	// allocation matters, so its API intentionally favors performance and type
	// safety over brevity. For most applications, the SugaredLogger strikes a
	// better balance between performance and ergonomics.
	Logger struct {
		// In the rare contexts where every microsecond and every allocation matter,
		// use the Logger. It's even faster than the SugaredLogger and allocates far less,
		// but it only supports strongly-typed, structured logging.
		logger *zap.Logger
	}

	// A SugaredLogger wraps the base Logger functionality in a slower, but less
	// verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar
	// method.
	//
	// Unlike the Logger, the SugaredLogger doesn't insist on structured logging.
	// For each log level, it exposes three methods: one for loosely-typed
	// structured logging, one for println-style formatting, and one for
	// printf-style formatting. For example, SugaredLoggers can produce InfoLevel
	// output with Infow ("info with" structured context), Info, or Infof.
	SugaredLogger struct {
		// In contexts where performance is nice, but not critical, use the SugaredLogger.
		sugar *zap.SugaredLogger
	}

	// LogOptions represent option to custom-zap logger
	//
	// Level set log's level logger, either development or production
	// Time set log's time location being used, default is "Asia/Jakarta".
	// Use according to Time Zone database, such as "America/New_York".
	// WithTrace set trace-id to logs output.
	// RefID will set ref-id to logs output.
	// Output file is another output file. If you want logger to write log
	// to multiple file, add other source here. add "stdout" for console log.
	LogOptions struct {
		Level      int
		Time       *time.Location
		WithTrace  bool
		RefID      string
		OutputFile []string
	}
)

// GetOutputFile returns log's filename name.
func GetOutputFile() string {
	return filename
}

// GetTimeLocation return time location used in logs
func GetTimeLocation() *time.Location {
	return timeLocation
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
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.StacktraceKey = ""
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

// newLogger return new custom zap-logger
//
// set default logger: logs to os.stdout, production level, with Asia/Jakarta time
// 		printed log : os.stdout
// 		filename 	: yyyy-mm-dd.log
// 		log level 	: production
// 		log time 	: "Asia/Jakarta"
func newLogger(dir string, prefix string, opt *LogOptions) *zap.Logger {
	if opt == nil {
		opt = &LogOptions{}
	}

	if opt.Level < 1 {
		opt.Level = LLvlProduction
	}

	if opt.Time == nil {
		opt.Time, _ = time.LoadLocation(LStdLocation)
	}
	timeLocation = opt.Time

	if dir != "" {
		filename = opt.makeLogFile(create(dir), prefix)
	}

	logger, err := opt.newConfig(filename).Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

// create set log's directory location and,
// create directory if not exist
// default is 				: __path__/log
// returned log's directory : __path__/dir
func create(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("create folder in: %s\n", dir)
		if err = os.MkdirAll(dir, os.ModePerm|os.ModeAppend); err != nil {
			panic(fmt.Sprintf("[log] failed to create directory: %v", err))
		}
	}
	return dir
}

func NewSugaredLogger(dir string, prefix string, opt *LogOptions) *SugaredLogger {
	return &SugaredLogger{sugar: newLogger(dir, prefix, opt).Sugar()}
}

// Log logs using zap log.
// msg is custom message
// params contains key-value message. used for tracing
// err is error
func (l *SugaredLogger) Log(msg string, params map[string]interface{}, err error) {
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
func (l *SugaredLogger) Fatal(v ...interface{}) {
	l.sugar.Fatal(v...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *SugaredLogger) Fatalf(format string, v ...interface{}) {
	l.sugar.Fatal(fmt.Sprintf(format, v...))
}

// NewLogger initiate new zap logger,
// by satisfy log's dir, prefix and options
func NewLogger(dir string, prefix string, opt *LogOptions) *Logger {
	return &Logger{logger: newLogger(dir, prefix, opt)}
}

// Debug logs the message at debug level with additional fields, if any
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Error logs the message at error level with additional fields, if any
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs the message at fatal level with additional fields, if any and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Info logs the message at info level with additional fields, if any
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs the message at warn level with additional fields, if any
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}
