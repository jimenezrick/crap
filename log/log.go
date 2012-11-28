package log

import (
	"os"
	"io/ioutil"
	"log"
	"log/syslog"
)

import "crap/kvmap"

var Info *log.Logger
var Debug *log.Logger

func Init(config *kvmap.KVMap) {
	debug, err := config.GetBool("log.debug")
	if err != nil {
		panic(err)
	}
	sys, err := config.GetBool("log.syslog")
	if err != nil {
		panic(err)
	}

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
