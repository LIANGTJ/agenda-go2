package agenda

import (
	"convention/agendaerror"
	"entity"
	"model"
	log "util/logger"
)

type Username = entity.Username
type Auth = entity.Auth

type UserInfo = entity.UserInfo
type User = entity.User
type MeetingInfo = entity.MeetingInfo
type Meeting = entity.Meeting
type MeetingTitle = entity.MeetingTitle
type UserInfoPublicList = entity.UserInfoPublicList

func MakeUserInfo(username Username, password Auth, email, phone string) UserInfo {
	info := UserInfo{}

	info.Name = username
	info.Auth = password
	info.Mail = email
	info.Phone = phone

	return info
}

var NewUser = entity.NewUser
var NewMeeting = entity.NewMeeting

var registeredUsers = entity.GetAllUsersRegistered()
var allMeetings = entity.GetAllMeetings()

func LoadAll() {
	model.Load()
}
func SaveAll() {
	if err := model.Save(); err != nil {
		log.Printf(err.Error())
	}
}

// NOTE: Now, assume the operations' actor are always the `Current User`

func LoginedUser() *User {
	name := Username("root")
	return name.RefInAllUsers()
}

// RegisterUser ...
// func RegisterUser(username, password, email, phone string) error {
//     u := NewUser(MakeUserInfo(username, password, email, phone))
// func RegisterUser(u *entity.User) error {
func RegisterUser(uInfo UserInfo) error {
	u := NewUser(uInfo)
	err := entity.GetAllUsersRegistered().Add(u)
	return err
}

// CancelAccount cancels(deletes) a User's account
func CancelAccount(name Username) error {
	u := name.RefInAllUsers()

	// check if under login status, TODO: check the login status
	if logined := LoginedUser(); logined == nil {
		return agendaerror.ErrUserNotLogined
	} else if logined != u {
		return agendaerror.ErrUserAuthority
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

// QueryAccountAll queries all accounts
func QueryAccountAll() UserInfoPublicList {
	ret := LoginedUser().QueryAccountAll()
	return ret
}
