package logger

import (
	"bytes"
	"fmt"
	"log"
)

type lotter struct {
	infobuf, warnbuf, errbuf bytes.Buffer
	info                     *log.Logger
	warn                     *log.Logger
	err                      *log.Logger
}

var tLogger *lotter

func init() {
	lt := new(lotter)
	var infob bytes.Buffer
	var warnb bytes.Buffer
	var errb bytes.Buffer

	lt.infobuf = infob
	lt.warnbuf = warnb
	lt.errbuf = errb

	lt.info = log.New(&lt.infobuf, "INFO: ", log.Lshortfile)
	lt.warn = log.New(&lt.warnbuf, "WARN: ", log.Lshortfile)
	lt.err = log.New(&lt.errbuf, "ERROR: ", log.Lshortfile)

	tLogger = lt
}

func Getlogger() *lotter {
	return tLogger
}

func (lt *lotter) Info(info string) {
	lt.info.Output(2, info)
	fmt.Print(&lt.infobuf)
	lt.infobuf.Reset()
}

func (lt *lotter) Warn(info string) {
	lt.warn.Output(2, info)
	fmt.Print(&lt.warnbuf)
	lt.warnbuf.Reset()
}

func (lt *lotter) Err(info string) {
	lt.err.Output(2, info)
	fmt.Print(&lt.errbuf)
	lt.errbuf.Reset()
}
