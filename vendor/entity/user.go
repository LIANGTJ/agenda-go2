package entity

import (
	"auth"
	"fmt"
	"log"
	"util"
)

var logln = util.Log
var logf = util.Logf

type Username = util.Identifier

// TODO: Not sure where to place ...
var UsersAllRegistered = *(NewUserList())

// TODO: if Username keeps its nonlocal-type, this func never okay ...
//      However, if Username be a stand-alone type, the func:Empty would screw up
// func (name Username) RefInAllUsers() *User {
// 	return UsersAllRegistered.Ref(name)
// }
func GetAllUsersRegistered() *UserList {
	return &UsersAllRegistered
}

type UserInfo struct {
	Name  Username
	Auth  auth.Auth
	Mail  string
	Phone string
}

// NOTE: Maybe not to export ?
// NOTE: Maybe to stand alone type ?
type UserInfoSerializable = UserInfo

type User struct {
	UserInfo
}

func NewUser(info UserInfo) *User {
	if info.Name.Empty() {
		// FIXME: more elegant ?
		// TODO: provide a ready-user that allowe to be empty to be loaded info into ?
		log.Printf("An empty UserInfo is passed to new a User. Just return `nil`.")
		return nil
	}
	u := new(User)
	u.UserInfo = info

	// NOTE: support `UsersAllRegistered`
	// TODO: Do something when re-newing a existed-user.
	UsersAllRegistered.Add(u)

	return u
}

func LoadUsersAllRegistered(decoder Decoder) {
	users := &(UsersAllRegistered)
	LoadUserList(decoder, users)
}

func LoadUser(decoder Decoder, u *User) {
	uInfo := new(UserInfo)
	err := decoder.Decode(uInfo)
	if err != nil {
		log.Fatal(err)
	}
	u.UserInfo = *uInfo
}
func LoadedUser(decoder Decoder) *User {
	u := new(User)
	LoadUser(decoder, u)
	return u
}

func (u *User) Save(encoder Encoder) error {
	return encoder.Encode(u.UserInfo)
}

func (u *User) CancelAccount() error {
	// check if under login status
	// notification.

	// del this account
	// del all meeting that this user is sponsor
	// remove this user from participators of all meeting that this user participate
	//      if removing cause people count < 0, del the meeting
	return ErrNeedImplement
}

func (u *User) QueryAccount() error {
	return ErrNeedImplement
}

func (u *User) CreateMeeting() (*Meeting, error) {
	return nil, ErrNeedImplement
}

func (u *User) LogOff() error {
	return ErrNeedImplement
}

func (u *User) meetingsSpoonsored() ([]*Meeting, error) {
	return nil, ErrNeedImplement
}

func (u *User) CancelMeeting() error { return ErrNeedImplement }

func (u *User) QuitMeeting() error { return ErrNeedImplement }

// ................................................................

type UserList struct {
	Users map[Username](*User)
}

type UserListRaw = []*User

type UserInfoList []UserInfo
type UserInfoListSerializable []UserInfoSerializable

func NewUserList() *UserList {
	ul := new(UserList)
	ul.Users = make(map[Username](*User))
	return ul
}

// CHECK: Need in-place load method ?
func LoadUserList(decoder Decoder, ul *UserList) {
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
func (ul *UserList) LoadFrom(decoder Decoder) {
	LoadUserList(decoder, ul)
}

func LoadedUserList(decoder Decoder) *UserList {
	ul := NewUserList()
	LoadUserList(decoder, ul)
	return ul
}

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

func (ul UserInfoListSerializable) Size() int {
	return len(ul)
}

func (ulSerial UserInfoListSerializable) Deserialize() *UserList {
	ret := NewUserList()

	for _, uInfo := range ulSerial {

		// FIXME: these are introduced since up to now, it is possible that UserList contains nil User
		// FIXME: Not use `== nil` because `uInfo` is a struct
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

func (ul *UserList) Save(encoder Encoder) error {
	sl := ul.Serialize()
	// logf("sl: %+v\n", sl)
	return encoder.Encode(sl)
}

func (ul *UserList) Size() int {
	return len(ul.Users)
}

func (ul *UserList) Ref(name Username) *User {
	return ul.Users[name] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ul *UserList) Contains(name Username) bool {
	u := ul.Ref(name)
	return u != nil
}

func (ul *UserList) Add(user *User) error {
	if user == nil {
		return ErrNilUser
	}
	name := user.Name
	if ul.Contains(name) {
		return ErrExistedUser
	}
	ul.Users[name] = user
	return nil
}
func (ul *UserList) Remove(user *User) error {
	if user == nil {
		return ErrNilUser
	}
	name := user.Name
	if ul.Contains(name) {
		delete(ul.Users, name) // NOTE: never error, according to 'go-maps-in-action'
		return nil
	}
	return ErrUserNotFound
}
func (ul *UserList) PickOut(name Username) (*User, error) {
	if name.Empty() {
		return nil, ErrEmptyUsername
	}
	u := ul.Ref(name)
	if u == nil {
		return u, ErrUserNotFound
	}
	defer ul.Remove(u)
	return u, nil
}

func (ul *UserList) Slice() UserListRaw {
	users := make(UserListRaw, 0, ul.Size())
	for _, u := range ul.Users {
		users = append(users, u) // CHECK: maybe better to use index in golang ?
	}
	return users
}

// CHECK: should limit func type ?
func (ul *UserList) ForEach(fn interface{}) error {
	switch f := fn.(type) {
	case func(Username) error:
		for k := range ul.Users {
			if err := f(k); err != nil {
				return err
			}
		}
	case func(*User) error:
		for _, v := range ul.Users {
			if err := f(v); err != nil {
				return err
			}
		}
	case func(Username, *User) error:
		for k, v := range ul.Users {
			if err := f(k, v); err != nil {
				return err
			}
		}
	default:
		return NewAgendaError(fmt.Sprintf("Given function has unmatched signature: %T", f))
	}
	return nil
}
