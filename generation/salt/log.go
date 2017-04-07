package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// logging utilities for the salt package, to provide clean and easy logging
// with simple visual tags with which to pick salt messages out of the log as
// a whole

var (
	debugLog veldt.Logger
	infoLog  veldt.Logger
	warnLog  veldt.Logger
	errorLog veldt.Logger
)

const (
	// Yellow "SALT" log prefix
	// Code derived from github.com/mgutz/ansi, but I wanted it as a const, so
	// couldn't se that directly
	preLog = "\033[1;38;5;3mSALT\033[0m: "
	// And codes to make the message red, similarly as constants
	preMsg  = "\033[1;97;3m"
	postMsg = "\033[0m"
)


// SetDebugLogger sets the debug level logger for the batch package
func SetDebugLogger (log veldt.Logger) {
	debugLog = log
}
// SetInfoLogger sets the info level logger for the batch package
func SetInfoLogger (log veldt.Logger) {
	infoLog = log
}
// SetWarnLogger sets the info level logger for the batch package
func SetWarnLogger (log veldt.Logger) {
	warnLog = log
}
// SetErrorLogger sets the info level logger for the batch package
func SetErrorLogger (log veldt.Logger) {
	errorLog = log
}

func getLogger (level int) veldt.Logger {
	if veldt.Error == level {
		if nil == errorLog {
			if nil == warnLog {
				if nil == infoLog {
					return debugLog
				}
				return infoLog
			}
			return warnLog
		}
		return errorLog
	} else if veldt.Warn == level {
		if nil == warnLog {
			if nil == infoLog {
				return debugLog
			}
			return infoLog
		}
		return warnLog
	} else if veldt.Info == level {
		if nil == infoLog {
			return debugLog
		}
		return infoLog
	} else if veldt.Debug == level {
		return debugLog
	}
	return nil
}

// Errorf logs to the error log
func Errorf(format string, args ...interface{}) {
	logger := getLogger(veldt.Error)
	if nil != logger {
		errorLog.Errorf(preLog + format, args...)
	} else {
		veldt.Errorf(preLog + format, args...)
	}
}

// Warnf logs to the warn log
func Warnf(format string, args ...interface{}) {
	logger := getLogger(veldt.Warn)
	if nil != logger {
		warnLog.Warnf(preLog + format, args...)
	} else {
		veldt.Warnf(preLog + format, args...)
	}
}

// Infof logs to the info log
func Infof(format string, args ...interface{}) {
	logger := getLogger(veldt.Info)
	if nil != logger {
		infoLog.Infof(preLog + format, args...)
	} else {
		veldt.Infof(preLog + format, args...)
	}
}

// Debugf logs to the debug log
func Debugf(format string, args ...interface{}) {
	logger := getLogger(veldt.Debug)
	if nil != logger {
		debugLog.Debugf(preLog + format, args...)
	} else {
		veldt.Debugf(preLog + format, args...)
	}
}
