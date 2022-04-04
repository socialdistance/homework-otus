package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
)

type Logger struct {
	logger *logrus.Logger
}

func (l *Logger) Debug(msg string, params ...interface{}) {
	l.logger.Debugf(msg, params...)
}

func (l *Logger) Info(msg string, params ...interface{}) {
	l.logger.Infof(msg, params...)
}

func (l *Logger) Error(msg string, params ...interface{}) {
	l.logger.Errorf(msg, params...)
}

func (l *Logger) Warn(msg string, params ...interface{}) {
	l.logger.Warnf(msg, params...)
}

func (l *Logger) LogHTTP(r *http.Request, code, length int) {
	l.logger.Infof(
		"%s [%s] %s %s %s %d %d %q",
		r.RemoteAddr,
		time.Now().Format("01/Jan/2003:10:10:10 MST"),
		r.Method,
		r.RequestURI,
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

func New(loggerConf config.LoggerConf) (*Logger, error) {
	logger := logrus.New()

	loggerOut, err := parseFile(loggerConf.Filename)
	if err != nil {
		return nil, fmt.Errorf("invalid log file %w", err)
	}
	logger.SetOutput(loggerOut)

	logLevel, err := logrus.ParseLevel(string(loggerConf.Level))
	if err != nil {
		return nil, fmt.Errorf("failed parse level %w", err)
	}

	logrus.SetLevel(logLevel)

	logger.SetFormatter(&logrus.JSONFormatter{})

	return &Logger{
		logger,
	}, nil
}

func parseFile(filename string) (io.Writer, error) {
	switch filename {
	case "stderr":
		return os.Stderr, nil
	case "stdout":
		return os.Stdout, nil
	default:
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, err
		}

		return file, nil
	}
}
