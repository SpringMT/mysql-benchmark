package main

import (
	"fmt"
	"log"
	"strings"
)

type loglevel int

const (
	mute loglevel = iota
	warn
	info
	debug
)

func (mb *mysqlBench) logLevel() loglevel {
	return loglevel(len(mb.Verbose))
}

func (mb *mysqlBench) logf(lv loglevel, format string, a ...interface{}) {
	mb.log(lv, fmt.Sprintf(format, a...))
}

func (mb *mysqlBench) log(lv loglevel, str string) {
	if mb.logLevel() < lv || lv <= mute {
		return
	}
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	log.Print(str)
}
