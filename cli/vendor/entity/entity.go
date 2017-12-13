package entity

import (
	"time"
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

func (*Meeting) NewMeeting(title string, 
	startTime, endTime time.Time, participators []string) Meeting {

	return Meeting {
		title,
		participators,
		startTime,
		endTime,
	}
}

func (*User) NewUser(username, password, email, phone string) User {

	return User {
		username,
		password,
		email,
		phone,
	}
}

