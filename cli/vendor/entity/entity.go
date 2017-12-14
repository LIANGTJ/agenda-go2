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

type SerializedMeeting struct {
	
	Title 			string
	Participators 	[]string
	StartTime 		string
	EndTime 		string
	
}

type MeetingList []Meeting
type SerializedMeetingList []SerializedMeeting

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

func (m *Meeting)Serialized() *SerializedMeeting{
	return &SerializedMeeting{
		m.Title,
		m.Participators,
		m.StartTime.String(),
		m.EndTime.String(),
	}
}

type User struct {
	
		Username string
		Password string
		Email 	 string
		Phone 	 string
	}
	
type UserList []User

func NewUser(username, password, email, phone string) *User {

	return 	&User {
		username,
		password,
		email,
		phone,
	}
	
}

func ( u *User) Invalid() bool {
	return u.Username == "" || u.Password == ""
}






