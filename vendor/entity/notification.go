package entity

import (
	"time"
)

type Notification struct {
    Type string
    Text string
    Timestamp time.Time
    Source *interface{}
}
