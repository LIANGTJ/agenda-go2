package entity

import (
	"time"
	"util"
)

var logln = util.Log
var logf = util.Logf

type Notification struct {
	Type      string
	Text      string
	Timestamp time.Time
	Source    *interface{}
}

func Notify(text string) {
	logf(text)
}
