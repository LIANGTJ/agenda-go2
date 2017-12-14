package entity

import (
	"time"
	// log "util/logger"
	"errors"
)

type Meeting struct {

	Title 			string
	Participators 	[]string
	StartTime 		time.Time
	EndTime 		time.Time

}

type MeetingList []Meeting


type User struct {

	Username string
	Password string
	Email 	 string
	Phone 	 string
}

func NewMeeting(title string, participators []string, startTime, endTime time.Time) (*Meeting, error) {
	var err error = nil
	if startTime.After(endTime) {
		// log.Fatal("the start-time should be before the end-time")
		err = errors.New("the start-time must go ahead of the end-time")
	}

	meeting := Meeting {
		title,
		participators,
		startTime,
		endTime,
	}
	return &meeting, err
}

func NewUser(username, password, email, phone string) *User {

	u := User {
		username,
		password,
		email,
		phone,
	}
	return &u
}

func ( u *User) Invalid() bool {
	return u.Username == "" || u.Password == ""
}





