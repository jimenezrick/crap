package log

import (
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
)

import "crap/config"

var Info *log.Logger
var Debug *log.Logger

func Init(config config.Config) {
	debug := config.GetBool("log.debug")
	sys := config.GetBool("log.syslog")

	Debug = log.New(ioutil.Discard, "", 0)
	if sys {
		Info = newSyslogger(syslog.LOG_INFO)
		if debug {
			Debug = newSyslogger(syslog.LOG_DEBUG)
		}
	} else {
		Info = newLogger()
		if debug {
			Debug = newLogger()
		}
	}
}

func newSyslogger(p syslog.Priority) *log.Logger {
	logger, err := syslog.NewLogger(p, 0)
	if err != nil {
		panic(err)
	}
	return logger
}

func newLogger() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}
