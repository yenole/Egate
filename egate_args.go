package egate

import "flag"

var gWatch bool
var gLogs string

func init() {
	flag.BoolVar(&gWatch, "watch", false, "Enabling Daemons!")
	flag.StringVar(&gLogs, "logs", "", "Logs file dir!")
	flag.Parse()
}

func WatchMode() bool {
	return gWatch
}

func LogDirPath() string {
	return gLogs
}
