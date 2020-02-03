package logger

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// Logger wrap logrus.Logger
type Logger struct {
	*logrus.Logger
}

var (
	// Log is default logger
	// only output to os.Stdout
	Log = logrus.New()
)

// initialize default logger
func init() {
	Log.Formatter = &logrus.JSONFormatter{}
}

// NewLogger ...
// init logger to set filepath and format
func NewLogger(logPath, filename, lv string) (*Logger, error) {
	filenameHook := NewHook()
	filenameHook.Field = "source"

	logger := logrus.New()

	// set formatter
	logger.Formatter = &logrus.JSONFormatter{}

	// set level
	if level, err := logrus.ParseLevel(lv); err != nil {
		return nil, err
	} else {
		logger.SetLevel(level)
	}

	// set file output
	//
	// [fixed]: if folder not existed, will mkdir all path,
	// but the log file will be a folder.
	//
	fullPath := path.Join(logPath, filename)
	fd, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	// check logfile and ptah exist or not
	// if not exist , make dir and new file
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(logPath, 0777); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	logger.AddHook(filenameHook)
	logger.Out = io.MultiWriter(os.Stdout, fd)
	return &Logger{logger}, nil
}
