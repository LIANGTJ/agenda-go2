package agenda

import (
	"convention/agendaerror"
	"entity"
	"model"
	"time"
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
func MakeMeetingInfo(title MeetingTitle, sponsor Username, participators []Username, startTime, endTime time.Time) MeetingInfo {
	info := MeetingInfo{}

	info.Title = title
	info.Sponsor = sponsor.RefInAllUsers()
	info.Participators.InitFrom(participators)
	info.StartTime = startTime
	info.EndTime = endTime

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

// QueryAccountAll queries all accounts
func QueryAccountAll() UserInfoPublicList {
	ret := LoginedUser().QueryAccountAll()
	return ret
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

// SponsorMeeting creates a meeting
func SponsorMeeting(mInfo MeetingInfo) (*Meeting, error) {
	info := mInfo

	// NOTE: dev-assert
	if info.Sponsor != nil && info.Sponsor.Name != LoginedUser().Name {
		log.Fatalf("User %v is creating a meeting with Sponsor %v\n", LoginedUser().Name, info.Sponsor.Name)
	}

	// NOTE: repeat in MeetingList.Add ... DEL ?
	if info.Title.RefInAllMeetings() != nil {
		return nil, agendaerror.ErrExistedMeetingTitle
	}

	if !LoginedUser().Registered() {
		return nil, agendaerror.ErrUserNotRegistered
	}

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.Registered() {
			return agendaerror.ErrUserNotRegistered
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	if !info.EndTime.After(info.StartTime) {
		return nil, agendaerror.ErrInvalidTimeInterval
	}

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.FreeWhen(info.StartTime, info.EndTime) {
			return agendaerror.ErrConflictedTimeInterval
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	m, err := LoginedUser().SponsorMeeting(info)
	if err != nil {
		log.Printf("Failed to sponsor meeting, error: %q", err.Error())
	}
	return m, err
}

// AddParticipatorToMeeting ...
func AddParticipatorToMeeting(title MeetingTitle, name Username) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return agendaerror.ErrUserNotLogined
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return agendaerror.ErrNilMeeting
	}
	if user == nil {
		return agendaerror.ErrNilUser
	}

	if !meeting.SponsoredBy(u.Name) {
		return agendaerror.ErrSponsorAuthority
	}

	if meeting.ContainsParticipator(name) {
		return agendaerror.ErrExistedUser
	}

	if !user.FreeWhen(meeting.StartTime, meeting.EndTime) {
		return agendaerror.ErrConflictedTimeInterval
	}

	err := u.AddParticipatorToMeeting(meeting, user)
	if err != nil {
		log.Printf("Failed to add participator into Meeting, error: %q", err.Error())
	}
	return err
}

// RemoveParticipatorFromMeeting ...
func RemoveParticipatorFromMeeting(title MeetingTitle, name Username) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return agendaerror.ErrUserNotLogined
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return agendaerror.ErrMeetingNotFound
	}
	if user == nil {
		return agendaerror.ErrUserNotRegistered
	}

	if !meeting.SponsoredBy(u.Name) {
		return agendaerror.ErrSponsorAuthority
	}

	if !meeting.ContainsParticipator(name) {
		return agendaerror.ErrUserNotFound
	}

	err := u.RemoveParticipatorFromMeeting(meeting, user)
	if err != nil {
		log.Printf("Failed to remove participator from Meeting, error: %q", err.Error())
	}
	return err
}
