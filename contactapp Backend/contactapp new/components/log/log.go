package log

import (
	"fmt"
	"time"
)

type Logger interface {
	Print(value ...string)
	New()
}

type Log struct {
}

func GetLogger() *Log {
	return &Log{}
}

func (l *Log) Print(value ...interface{}) {
	curTime := time.Now()
	dateTime := curTime.Format(time.DateTime)
	fmt.Println(dateTime, value)
}
