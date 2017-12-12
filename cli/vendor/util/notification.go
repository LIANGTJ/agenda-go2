package util

import (
	"time"
)

type Notification struct {
	Type      string
	Text      string
	Timestamp time.Time
	Source    *interface{}
}

func Notify(text string) {
	Logf(text)
}
