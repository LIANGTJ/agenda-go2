package entity

import (
	"time"
	// log "util/logger"
	"errors"
)

type Meeting struct {

	Title 			string      `json:"title"`
	Participators 	[]string	`json:"participators"`
	StartTime 		time.Time	`json:"starttime"`
	EndTime 		time.Time	`json:"endtime"`

}

type SerializedMeeting struct {
	
	Title 			string		`json:"title"`
	Participators 	[]string	`json:"participators"`
	StartTime 		string		`json:"starttime"`
	EndTime 		string		`json:"endtime"`
	
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
	
		Username string	`json:"username"`
		Password string	`json:"password"`
		Email 	 string	`json:"email"`
		Phone 	 string	`json:"phone"`
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






