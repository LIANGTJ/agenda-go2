package agenda

import (
	errors "convention/agendaerror"
	"entity"
	"model"
	"time"
	log "util/logger"
)

type Username = entity.Username
type Auth = entity.Auth

type UserInfo = entity.UserInfo
type UserInfoPublic = entity.UserInfoPublic
type User = entity.User
type MeetingInfo = entity.MeetingInfo
type Meeting = entity.Meeting
type MeetingTitle = entity.MeetingTitle

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

func LoadAll() {
	model.Load()
	LoadLoginStatus()
}
func SaveAll() {
	if err := model.Save(); err != nil {
		log.Error(err)
	}
	SaveLoginStatus()
}

// NOTE: Now, assume the operations' actor are always the `Current User`

// RegisterUser ...
func RegisterUser(uInfo UserInfo) error {
	if !uInfo.Name.Valid() {
		return errors.ErrInvalidUsername
	}

	u := entity.NewUser(uInfo)
	err := entity.GetAllUsersRegistered().Add(u)
	return err
}

func LogIn(name Username, auth Auth) error {
	u := name.RefInAllUsers()
	if u == nil {
		return errors.ErrNilUser
	}

	log.Printf("User %v logs in.\n", name)

	if LoginedUser() != nil {
		return errors.ErrLoginedUserAuthority
	}

	if verified := u.Auth.Verify(auth); !verified {
		return errors.ErrFailedAuth
	}

	loginedUser = name

	return nil
}

// LogOut log out User's own (current working) account
// TODO:
func LogOut(name Username) error {
	u := name.RefInAllUsers()

	// check if under login status, TODO: check the login status
	if logined := LoginedUser(); logined == nil {
		return errors.ErrUserNotLogined
	} else if logined != u {
		return errors.ErrUserAuthority
	}

	err := u.LogOut()
	if err != nil {
		log.Errorf("Failed to log out, error: %q.\n", err.Error())
	}
	loginedUser = ""
	return err
}

// QueryAccountAll queries all accounts
func QueryAccountAll() []UserInfoPublic {
	// NOTE: FIXME: whatever, temporarily ignore the problem that the actor of query is Nil
	// Hence, now if so, agenda would crash for `Nil.Name`
	ret := LoginedUser().QueryAccountAll()
	return ret
}

// CancelAccount cancels(deletes) LoginedUser's account
func CancelAccount() error {
	u := LoginedUser()
	if u == nil {
		return errors.ErrUserNotLogined
	}

	if err := entity.GetAllMeetings().ForEach(func(m *Meeting) error {
		if m.SponsoredBy(u.Name) {
			return m.Dissolve()
		}
		if m.ContainsParticipator(u.Name) {
			return m.Exclude(u)
		}
		return nil
	}); err != nil {
		log.Error(err)
	}

	if err := entity.GetAllUsersRegistered().Remove(u); err != nil {
		log.Error(err)
	}
	if err := u.LogOut(); err != nil {
		log.Error(err)
	}

	err := u.CancelAccount()
	return err
}

// SponsorMeeting creates a meeting
func SponsorMeeting(mInfo MeetingInfo) (*Meeting, error) {
	u := LoginedUser()
	if u == nil {
		return nil, errors.ErrUserNotLogined
	}

	info := mInfo

	if !info.Title.Valid() {
		return nil, errors.ErrInvalidMeetingTitle
	}

	// NOTE: dev-assert
	if info.Sponsor != nil && info.Sponsor.Name != LoginedUser().Name {
		log.Fatalf("User %v is creating a meeting with Sponsor %v\n", LoginedUser().Name, info.Sponsor.Name)
	}

	// NOTE: repeat in MeetingList.Add ... DEL ?
	if info.Title.RefInAllMeetings() != nil {
		return nil, errors.ErrExistedMeetingTitle
	}

	// if !LoginedUser().Registered() { return nil, errors.ErrUserNotRegistered }

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.Registered() {
			return errors.ErrUserNotRegistered
		}
		return nil
	}); err != nil {
		log.Error(err)
		return nil, err
	}

	if !info.EndTime.After(info.StartTime) {
		return nil, errors.ErrInvalidTimeInterval
	}

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.FreeWhen(info.StartTime, info.EndTime) {
			return errors.ErrConflictedTimeInterval
		}
		return nil
	}); err != nil {
		log.Error(err)
		return nil, err
	}

	m, err := LoginedUser().SponsorMeeting(info)
	if err != nil {
		log.Errorf("Failed to sponsor meeting, error: %q.\n", err.Error())
	}
	return m, err
}

// AddParticipatorToMeeting ...
func AddParticipatorToMeeting(title MeetingTitle, name Username) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return errors.ErrUserNotLogined
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return errors.ErrNilMeeting
	}
	if user == nil {
		return errors.ErrNilUser
	}

	if !meeting.SponsoredBy(u.Name) {
		return errors.ErrSponsorAuthority
	}

	if meeting.ContainsParticipator(name) {
		return errors.ErrExistedUser
	}

	if !user.FreeWhen(meeting.StartTime, meeting.EndTime) {
		return errors.ErrConflictedTimeInterval
	}

	err := u.AddParticipatorToMeeting(meeting, user)
	if err != nil {
		log.Errorf("Failed to add participator into Meeting, error: %q.\n", err.Error())
	}
	return err
}

// RemoveParticipatorFromMeeting ...
func RemoveParticipatorFromMeeting(title MeetingTitle, name Username) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return errors.ErrUserNotLogined
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return errors.ErrMeetingNotFound
	}
	if user == nil {
		return errors.ErrUserNotRegistered
	}

	if !meeting.SponsoredBy(u.Name) {
		return errors.ErrSponsorAuthority
	}

	if !meeting.ContainsParticipator(name) {
		return errors.ErrUserNotFound
	}

	err := u.RemoveParticipatorFromMeeting(meeting, user)
	if err != nil {
		log.Errorf("Failed to remove participator from Meeting, error: %q.\n", err.Error())
	}
	return err
}

func QueryMeetingByInterval(start, end time.Time, name Username) entity.MeetingInfoListPrintable {
	// NOTE: FIXME: whatever, temporarily ignore the problem that the actor of query is Nil
	// Hence, now if so, agenda would crash for `Nil.Name`
	ret := LoginedUser().QueryMeetingByInterval(start, end)
	return ret
}

// CancelMeeting cancels(deletes) the given meeting which sponsored by LoginedUser self
func CancelMeeting(title MeetingTitle) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return errors.ErrUserNotLogined
	}

	meeting := title.RefInAllMeetings()
	if meeting == nil {
		return errors.ErrMeetingNotFound
	}

	if !meeting.SponsoredBy(u.Name) {
		return errors.ErrSponsorAuthority
	}

	err := u.CancelMeeting(meeting)
	if err != nil {
		log.Errorf("Failed to cancel Meeting, error: %q.\n", err.Error())
	}
	return err
}

// QuitMeeting let LoginedUser quit the given meeting
func QuitMeeting(title MeetingTitle) error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return errors.ErrUserNotLogined
	}

	meeting := title.RefInAllMeetings()
	if meeting == nil {
		return errors.ErrMeetingNotFound
	}

	// CHECK: what to do in case User is exactly the sponsor ?
	// for now, refuse that
	if meeting.SponsoredBy(u.Name) {
		return errors.ErrSponsorResponsibility
	}

	if !meeting.ContainsParticipator(u.Name) {
		return errors.ErrUserNotFound
	}

	err := u.QuitMeeting(meeting)
	if err != nil {
		log.Errorf("Failed to quit Meeting, error: %q.\n", err.Error())
	}
	return err
}

// ClearAllMeeting cancels all meeting sponsored by LoginedUser
func ClearAllMeeting() error {
	u := LoginedUser()

	// check if under login status, TODO: check the login status
	if u == nil {
		return errors.ErrUserNotLogined
	}

	if err := entity.GetAllMeetings().ForEach(func(m *Meeting) error {
		if m.SponsoredBy(u.Name) {
			return CancelMeeting(m.Title)
		}
		return nil
	}); err != nil {
		log.Errorf("Failed to clear all Meetings, error: %q.\n", err.Error())
		return err
	}
	return nil
}
