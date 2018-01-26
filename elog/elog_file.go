package elog

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	CYCLE_DAY  time.Duration = 1
	CYCLE_WEEK time.Duration = 7

	time_layout = "2006-01-02"
)

type logFile struct {
	cycle   time.Duration
	dirPath string
	logger  *log.Logger
}

func (l *logFile) Sync() {
	if _, err := os.Stat(l.dirPath); err != nil && os.IsNotExist(err) {
		os.Mkdir(l.dirPath, 0755)
	}
	l.gSync()
}

func (l *logFile) gSync() {
	defer time.AfterFunc(l.duration(), l.gSync)
	logFile := fmt.Sprintf("%v/app.%v.log", l.dirPath, time.Now().Format(time_layout))
	if file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755); err == nil {
		l.logger.SetOutput(file)
	}
}

func (l *logFile) duration() time.Duration {
	dt, _ := time.ParseInLocation(time_layout, time.Now().Format(time_layout), time.Local)
	return dt.Add(l.cycle * 24 * time.Hour).Sub(time.Now())
}

func FileCycleMode(dir string, cycle time.Duration) {
	log := &logFile{cycle: cycle, dirPath: dir, logger: gLogger}
	log.Sync()
}
