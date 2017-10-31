package agenda

import (
	"convention/agendaerror"
	"entity"
	"log"
)

type Username = entity.Username
type User = entity.User
type Meeting = entity.Meeting
type MeetingTitle = entity.MeetingTitle

var registeredUsers = entity.GetAllUsersRegistered()
var allMeetings = entity.GetAllMeetings()

// NOTE: Now, assume the operations' actor are always the `Current User`

func LoginedUser() *User {

}

// RegisterUser ...
func RegisterUser(u *entity.User) error {

	err := entity.GetAllUsersRegistered().Add(u)
	return err
}

// CancelAccount cancels(deletes) a User's account
func CancelAccount(name Username) error {
	u := name.RefInAllUsers()

	// check if under login status, TODO: check the login status
	if logined := LoginedUser(); logined == nil {
		return agendaerror.UserNotLogined
	} else if logined != u {
		return agendaerror.UserAuthority
	}

	// del all meeting that this user is sponsor
	// remove this user from participators of all meeting that this user participate
	//      if removing cause people count < 0, del the meeting
	if err := allMeetings.ForEach(func(m *Meeting) error {
		if m.SponsoredBy(u.Name) {
			return m.Dissolve()
		}
		if m.ContainsParticipator(u.Name) {
			return m.Exclude(u)
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
	}

	if err := registeredUsers.Remove(u); err != nil {
		log.Printf(err.Error())
	}
	if err := u.LogOut(); err != nil {
		log.Printf(err.Error())
	}

	// Notify("CancelAccount: OK.")  TODEL: this should be in `cmd` module, notify dependon error

	err := u.CancelAccount()
	return err
}
