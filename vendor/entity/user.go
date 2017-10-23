package entity

import (
	"auth"
	"fmt"
	"log"
	"util"
)

var LOG = util.Log
var LOGF = util.Logf

type Username = util.Identifier

type UserInfo struct {
	Name  Username
	Auth  auth.Auth
	Mail  string
	Phone string
}

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
	return u
}

func DeserializeUser(decoder Decoder) (*User, error) {
	uInfo := new(UserInfo)
	err := decoder.Decode(uInfo)
	if err != nil {
		log.Fatal(err) // FIXME:
		return nil, err
	}
	user := NewUser(*uInfo)
	return user, nil
}

func (u User) Serialize(encoder Encoder) error {
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

const (
	defaultUserListLength = 5
)

/* func NewUserList() *UserList {
	users := make(UserList, defaultUserListLength)

	return &UserList{Users: users}
} */
func NewUserList() *UserList {
	ul := new(UserList)
	ul.Users = make(map[Username](*User))
	return ul
}

func DeserializeUserList(decoder Decoder) (*UserList, error) {
	ul := NewUserList()
	for decoder.More() {
		uInfo := new(UserInfo)
		err := decoder.Decode(uInfo)
		if err != nil {
			log.Fatal(err) // FIXME:
			return nil, err
		}
		user := NewUser(*uInfo)
		if err := ul.Add(user); err != nil {
			log.Fatal(err)
			return ul, err // FIXME:
		}
	}
	return ul, nil
}

func (ul *UserList) Serialize(encoder Encoder) error {
	sl := ul.toSerializable()
	// LOGF("sl: %+v\n", sl)
	return encoder.Encode(sl)
}

// TODO:  Need in-place deserialize ?
func (ul *UserList) Deserialize(decoder Decoder) error {
	// users, err := DeserializeUserList(decoder)
	// if err == nil {
	// 	// ul <- users
	// }
	return ErrNeedImplement
}

func (ul UserList) Size() int {
	return len(ul.Users)
}

func (ul *UserList) Get(name Username) *User {
	return ul.Users[name] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ul UserList) Contains(name Username) bool {
	u := ul.Get(name)
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
	u := ul.Get(name)
	if u == nil {
		return u, ErrUserNotFound
	}
	defer ul.Remove(u)
	return u, nil
}

func (ul UserList) Slice() UserListRaw {
	// NOTE: when make a `Slice`, the `len` value decide the position `append` to TODEL:
	users := make(UserListRaw, 0, ul.Size())
	for _, u := range ul.Users {
		users = append(users, u) // CHECK: maybe better to use index in golang ?
	}
	return users
}

func (ul UserList) toSerializable() []UserInfo {
	users := ul.Slice()
	ret := make([]UserInfo, 0, ul.Size())

	// LOG("ul.Size(): ", ul.Size())
	LOGF("toSerializable: %+v \n", users)
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

// func (ul *UserList) ForEach(f func(*User)) error {
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
