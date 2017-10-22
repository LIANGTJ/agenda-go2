package entity

import "time"

type Meeting struct {
	// ID            int
	Title         string
	Sponsor       *User
	Participators UserList
	// Time          time.Time
	StartTime time.Time
	EndTime   time.Time
}

func NewMeeting() {

}
