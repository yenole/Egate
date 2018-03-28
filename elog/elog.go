package elog

import (
	"log"
	"os"
	"strings"
)

const (
	debugLevel = 0
	infoLevel  = 1
	errorLevel = 2
	fatalLevel = 3
)
const (
	printDebugLevel = "[DEBUG] "
	printInfoLevel  = "[INFO ] "
	printErrorLevel = "[ERROR] "
	printFatalLevel = "[FATAL] "
)

var gLevel int
var gLogger = log.New(os.Stderr, "", log.LstdFlags)

func SetLevel(strLevel string) {
	switch strings.ToLower(strLevel) {
	case "debug":
		gLevel = debugLevel
	case "info":
		gLevel = infoLevel
	case "error":
		gLevel = errorLevel
	case "fatal":
		gLevel = fatalLevel
	default:
		log.Fatalf("unknown level: %v", strLevel)
	}
}

func doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < gLevel {
		return
	}
	format = printLevel + format
	gLogger.Printf(format, a...)
}

func Debug(format string, a ...interface{}) {
	doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Info(format string, a ...interface{}) {
	doPrintf(infoLevel, printInfoLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	doPrintf(fatalLevel, printFatalLevel, format, a...)
}
