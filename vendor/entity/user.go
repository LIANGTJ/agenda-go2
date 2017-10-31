package entity

import (
	"auth"
	"convention/codec"
	"convention/err"
	"log"
	"time"
	agnedaLog "util/logger"
)

var agnedaLogger = agnedaLog.Logger

// var logln = util.Log
// var logf = util.Logf

// Username represents username, a unique identifier, of User
// Identifier
type Username string

// Empty checks if Username empty
func (name Username) Empty() bool {
	return name == ""
}

// Valid checks if Username valid
func (name Username) Valid() bool {
	// FIXME: not only !empty
	return !name.Empty()
}

// TODO: Not sure where to place ...
var allUsersRegistered = *(NewUserList())

// RefInAllUsers returns the ref of a Registered User depending on the Username
func (name Username) RefInAllUsers() *User {
	return allUsersRegistered.Ref(name)
}

// GetAllUsersRegistered returns the reference of the UserList of all Registered Users
func GetAllUsersRegistered() *UserList {
	return &allUsersRegistered
}

// UserInfo represents the informations of a User
type UserInfo struct {
	Name  Username
	Auth  auth.Auth
	Mail  string
	Phone string
}

// UserInfoSerializable represents serializable UserInfo
type UserInfoSerializable = UserInfo

// User represents a User, which is the actor of the operations like sponsor/join/cancel a meeting, etc
type User struct {
	UserInfo
}

// NewUser creates a User object with given UserInfo
func NewUser(info UserInfo) *User {
	if info.Name.Empty() {
		// FIXME: more elegant ?
		// TODO: provide a ready-user that allowe to be empty to be loaded info into ?
		log.Printf("An empty UserInfo is passed to new a User. Just return `nil`.")
		return nil
	}
	u := new(User)
	u.UserInfo = info
	return u
}

// LoadUsersAllRegistered concretely loads all Registered Users
func LoadUsersAllRegistered(decoder codec.Decoder) {
	users := &(allUsersRegistered)
	LoadUserList(decoder, users)
}

// SaveUsersAllRegistered concretely saves all Registered Users
func SaveUsersAllRegistered(encoder codec.Encoder) error {
	users := &(allUsersRegistered)
	return users.Save(encoder)
}

// LoadUser load a User into given container(u) from given decoder
func LoadUser(decoder codec.Decoder, u *User) {
	uInfo := new(UserInfo)
	err := decoder.Decode(uInfo)
	if err != nil {
		log.Fatal(err)
	}
	u.UserInfo = *uInfo
}

// LoadedUser returns loaded User from given decoder
func LoadedUser(decoder codec.Decoder) *User {
	u := new(User)
	LoadUser(decoder, u)
	return u
}

// Save saves User with given encoder
func (u *User) Save(encoder codec.Encoder) error {
	return encoder.Encode(u.UserInfo)
}

func (u *User) Registered() bool {
	if u == nil {
		return false
	}
	return GetAllUsersRegistered().Contains(u.Name)
}

func (u *User) involvedMeetings() *MeetingList {
	return GetAllMeetings().Filter(func(m Meeting) bool {
		return m.SponsoredBy(u.Name) || m.ContainsParticipator(u.Name)
	})
}

func (u *User) FreeWhen(start, end time.Time) bool {
	if u == nil {
		return false
	}

	// NOTE: need improve:
	if err := u.involvedMeetings().ForEach(func(m *Meeting) error {
		s1, e1 := m.StartTime, m.EndTime
		s2, e2 := start, end
		if s1.Before(e2) && e1.After(s2) {
			return err.ConflictedTimeInterval
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
		return false
	}

	return true
}

// CancelAccount cancels(deletes) the User's own account
func (u *User) CancelAccount() error {
	agnedaLogger.Printf("User %v cancels account.", u.Name)
	return nil
}

// QueryAccount queries an account, where User as the actor
func (u *User) QueryAccount() error {
	return err.NeedImplement
}

// QueryAccountAll queries all accounts, where User as the actor
func (u *User) QueryAccountAll() UserInfoList {
	return GetAllUsersRegistered().Infos()
}

// CreateMeeting creates a meeting, where User as the actor
func (u *User) CreateMeeting(info MeetingInfo) (*Meeting, error) {
	// NOTE: dev-assert
	if info.Sponsor != nil && info.Sponsor.Name != u.Name {
		log.Fatalf("User %v is creating a meeting with Sponsor %v\n", u.Name, info.Sponsor.Name)
	}

	// NOTE: repeat in MeetingList.Add ... DEL ?
	if info.Title.RefInAllMeetings() != nil {
		return nil, err.ExistedMeetingTitle
	}

	if !u.Registered() {
		return nil, err.UserNotRegistered
	}

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.Registered() {
			return err.UserNotRegistered
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	if !info.EndTime.After(info.StartTime) {
		return nil, err.InvalidTimeInterval
	}

	if err := info.Participators.ForEach(func(u *User) error {
		if !u.FreeWhen(info.StartTime, info.EndTime) {
			return err.ConflictedTimeInterval
		}
		return nil
	}); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	m := NewMeeting(info)
	err := GetAllMeetings().Add(m)
	return m, err
}

// AddParticipatorToMeeting just as its name
func (u *User) AddParticipatorToMeeting(title MeetingTitle, name Username) error {
	if u == nil {
		return err.NilUser
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return err.NilMeeting
	}
	if user == nil {
		return err.NilUser
	}

	if !meeting.SponsoredBy(u.Name) {
		return err.SponsorAuthority
	}

	if meeting.ContainsParticipator(name) {
		return err.ExistedUser
	}

	if !user.FreeWhen(meeting.StartTime, meeting.EndTime) {
		return err.ConflictedTimeInterval
	}

	return meeting.Involve(user)
}

// RemoveParticipatorFromMeeting just as its name
func (u *User) RemoveParticipatorFromMeeting(title MeetingTitle, name Username) error {
	if u == nil {
		return err.NilUser
	}

	meeting, user := title.RefInAllMeetings(), name.RefInAllUsers()
	if meeting == nil {
		return err.MeetingNotFound
	}
	if user == nil {
		return err.UserNotRegistered
	}

	if !meeting.SponsoredBy(u.Name) {
		return err.SponsorAuthority
	}

	if !meeting.ContainsParticipator(name) {
		return err.UserNotFound
	}

	return meeting.Exclude(user)
}

// LogOut log out User's own (current working) account
func (u *User) LogOut() error {
	return err.NeedImplement
}

func (u *User) QueryMeetingByInterval(start, end time.Time) MeetingInfoListPrintable {
	return u.involvedMeetings().Textualize()
}

func (u *User) meetingsSponsored() ([]*Meeting, error) {
	return nil, err.NeedImplement
}

// CancelMeeting cancels(deletes) the given meeting which sponsored by User self, where User as the actor
func (u *User) CancelMeeting(title MeetingTitle) error {
	if u == nil {
		return err.NilUser
	}
	meeting := title.RefInAllMeetings()
	if meeting == nil {
		return err.MeetingNotFound
	}

	if !meeting.SponsoredBy(u.Name) {
		return err.SponsorAuthority
	}

	return meeting.Dissolve()
}

// QuitMeeting quits from the given meeting, where User as the actor
// CHECK: what to do in case User is the sponsor ?
func (u *User) QuitMeeting(title MeetingTitle) error {
	if u == nil {
		return err.NilUser
	}
	meeting := title.RefInAllMeetings()
	if meeting == nil {
		return err.MeetingNotFound
	}

	if meeting.SponsoredBy(u.Name) {
		return err.SponsorResponsibility // NOTE: ???
	}

	if !meeting.ContainsParticipator(u.Name) {
		return err.UserNotFound
	}

	return meeting.Exclude(u)
}

// ................................................................

// UserList represents a list/group of User (of the form of pointers of Users)
type UserList struct {
	Users map[Username](*User)
}

// UserListRaw also represents a list of User, but it is more trivial and more simple, i.e. it basically is ONLY a list of User, besides this, nothing
// NOTE: these type may be modified/removed in future
type UserListRaw = []*User

type UserInfoList []UserInfo

// UserInfoListSerializable represents a list of serializable UserInfo
type UserInfoListSerializable []UserInfoSerializable

// Serialize just serializes from UserList to UserInfoListSerializable
func (ul *UserList) Serialize() UserInfoListSerializable {
	users := ul.Slice()
	ret := make(UserInfoListSerializable, 0, ul.Size())

	// logln("ul.Size(): ", ul.Size())
	// logf("Serialize: %+v \n", users)
	for _, u := range users {

		// FIXME: these are introduced since up to now, it is possible that UserList contains nil User
		if u == nil {
			log.Printf("A nil User is to be used. Just SKIP OVER it.")
			continue
		}

		ret = append(ret, u.UserInfo)
	}
	return ret
}

// Size just returns the size
func (ulSerial UserInfoListSerializable) Size() int {
	return len(ulSerial)
}

// Deserialize deserializes from serialized UserInfoList to UserList
func (ulSerial UserInfoListSerializable) Deserialize() *UserList {
	ret := NewUserList()

	for _, uInfo := range ulSerial {

		// FIXME: these are introduced since up to now, it is possible that UserList contains nil User
		// FIXME: Not use `== nil` because `uInfo` is a  struct
		if uInfo.Name.Empty() {
			log.Printf("A No-Name UserInfo is to be used. Just SKIP OVER it.")
			continue
		}

		u := NewUser(uInfo)
		if err := ret.Add(u); err != nil {
			log.Printf(err.Error()) // CHECK:
		}
	}
	return ret
}

// NewUserList creates a UserList object
func NewUserList() *UserList {
	ul := new(UserList)
	ul.Users = make(map[Username](*User))
	return ul
}

// NOTE: these API (about loading) may be modified in future
// CHECK: Need in-place load method ?

// LoadUserList loads a UserList into given container(ul) from given decoder
func LoadUserList(decoder codec.Decoder, ul *UserList) {
	// CHECK: Need clear ul ?
	// for decoder.More() {
	// 	user := LoadedUser(decoder)
	// 	if err := ul.Add(user); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	ulSerial := new(UserInfoListSerializable)
	if err := decoder.Decode(ulSerial); err != nil {
		log.Fatal(err)
	}
	for _, uInfoSerial := range *ulSerial {
		u := NewUser(uInfoSerial)
		if err := ul.Add(u); err != nil {
			log.Printf(err.Error())
		}
	}
}

// LoadFrom loads UserList in-place from given decoder; Just like `init`
func (ul *UserList) LoadFrom(decoder codec.Decoder) {
	LoadUserList(decoder, ul)
}

// LoadedUserList returns loaded UserList from given decoder
func LoadedUserList(decoder codec.Decoder) *UserList {
	ul := NewUserList()
	LoadUserList(decoder, ul)
	return ul
}

func (ul *UserList) identifiers() []Username {
	ret := make([]Username, 0, ul.Size())
	for _, u := range ul.Infos() {
		ret = append(ret, u.Name)
	}
	return ret
}
func (ul *UserList) Infos() UserInfoList {
	users := ul.Slice()
	ret := make(UserInfoList, 0, ul.Size())

	for _, u := range users {

		// FIXME: these are introduced since up to now, it is possible that UserList contains nil User
		if u == nil {
			log.Printf("A nil User is to be used. Just SKIP OVER it.")
			continue
		}

		ret = append(ret, u.UserInfo)
	}
	return ret
}

// Save use given encoder to Save UserList
func (ul *UserList) Save(encoder codec.Encoder) error {
	sl := ul.Serialize()
	// logf("sl: %+v\n", sl)
	return encoder.Encode(sl)
}

// Size just returns the number of User reference in UserList
func (ul *UserList) Size() int {
	return len(ul.Users)
}

// Ref just returns the reference of user with the given name
func (ul *UserList) Ref(name Username) *User {
	return ul.Users[name] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}

// Contains just check if contains
func (ul *UserList) Contains(name Username) bool {
	u := ul.Ref(name)
	return u != nil
}

// Add just add
func (ul *UserList) Add(user *User) error {
	if user == nil {
		return err.NilUser
	}
	name := user.Name
	if ul.Contains(name) {
		return err.ExistedUser
	}
	ul.Users[name] = user
	return nil
}

// Remove just remove
func (ul *UserList) Remove(user *User) error {
	if user == nil {
		return err.NilUser
	}
	name := user.Name
	if ul.Contains(name) {
		delete(ul.Users, name) // NOTE: never error, according to 'go-maps-in-action'
		return nil
	}
	return err.UserNotFound
}

// PickOut =~= Ref and then Remove
func (ul *UserList) PickOut(name Username) (*User, error) {
	if name.Empty() {
		return nil, err.EmptyUsername
	}
	u := ul.Ref(name)
	if u == nil {
		return u, err.UserNotFound
	}
	defer ul.Remove(u)
	return u, nil
}

// Slice returns a UserListRaw based on UserList ul
// NOTE: for the need of this simple agenda system, this seems somewhat needless
func (ul *UserList) Slice() UserListRaw {
	users := make(UserListRaw, 0, ul.Size())
	for _, u := range ul.Users {
		users = append(users, u) // CHECK: maybe better to use index in golang ?
	}
	return users
}

// ForEach used for all extension/concrete logic for whole UserList
func (ul *UserList) ForEach(fn func(*User) error) error {
	for _, v := range ul.Users {
		if err := fn(v); err != nil {
			// CHECK: Or, lazy error ?
			return err
		}
	}
	return nil
}

// Filter used for all extension/concrete select for whole UserList
func (ul *UserList) Filter(pred func(User) bool) *UserList {
	ret := NewUserList()
	for _, u := range ul.Users {
		if pred(*u) {
			ret.Add(u)
		}
	}
	return ret
}
