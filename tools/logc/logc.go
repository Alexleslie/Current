package logc

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	flag     = log.Ldate | log.Ltime | log.Lshortfile
	preInfo  = "[INFO] "
	preDebug = "[DEBUG] "
	preWarn  = "[WARN] "
	preError = "[ERROR] "
)

var (
	logFile        io.Writer
	infoLogger     *log.Logger
	debugLogger    *log.Logger
	warnLogger     *log.Logger
	errorLogger    *log.Logger
	stdLogger      *log.Logger
	defaultLogFile = "C:\\web.log"
)

func toFileAndStd(logger *log.Logger, format string, values ...interface{}) {
	logger.Output(3, fmt.Sprintf(format, values))
	stdLogger.Output(3, fmt.Sprintf(format, values))
}

func Info(format string, values ...interface{}) {
	toFileAndStd(infoLogger, format, values)
}

func Debug(format string, values ...interface{}) {
	toFileAndStd(debugLogger, format, values)
}
func Warn(format string, values ...interface{}) {
	toFileAndStd(warnLogger, format, values)
}
func Error(format string, values ...interface{}) {
	toFileAndStd(errorLogger, format, values)
}

func SetOutPutPath(logFile io.Writer) {
	infoLogger.SetOutput(logFile)
	debugLogger.SetOutput(logFile)
	warnLogger.SetOutput(logFile)
	errorLogger.SetOutput(logFile)

}

func init() {
	var err error
	logFile, err = os.OpenFile(defaultLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		defaultLogFile = "./web.log"
		logFile, err = os.OpenFile(defaultLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Panicf("Create Log File error, err=[%+v]", err)
		}
	}
	infoLogger = log.New(logFile, preInfo, flag)
	debugLogger = log.New(logFile, preDebug, flag)
	warnLogger = log.New(logFile, preWarn, flag)
	errorLogger = log.New(logFile, preError, flag)
	stdLogger = log.New(os.Stderr, "", flag)

	SetOutPutPath(logFile)
}
