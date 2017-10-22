package main

// import "errors"
import "fmt"

import "errors"

type Meeting = string

var (
    ErrNeedImplement = errors.New("This function need to be implemented.")

	ErrNilUser      = errors.New("A nil user/*user is to be used.")
	ErrExistedUser  = errors.New("The user has been existed.")
	ErrUserNotFound = errors.New("Cannot find the user.")

	ErrEmptyUsername = errors.New("Given username cannot be empty.")
	// Err
	// Err
)

type AgendaError struct {
	msg string
}

func NewAgendaError(msg string) *AgendaError {
	return &AgendaError{
		msg: msg,
	}
}

func (e *AgendaError) Error() string {
	return e.msg
}

//
//
//

type Username string

var emptyUsername = *new(Username)

func (n Username) Empty() bool {
	return n == emptyUsername
}

type UserInfo struct {
	Name  Username
	// Auth  auth.Auth
	Mail  string
	Phone string
}

type User struct {
	UserInfo
}

func NewUser(info UserInfo) *User {
	if info.Name.Empty() {
		// FIXME: more elegant ?
		return nil
	}
	u := new(User)
	u.UserInfo = info
	return u
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

func (ul UserList) Size() int {
	return len(ul.Users)
}

func (ul *UserList) Get(name Username) *User {
	return ul.Users[name] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ul UserList) Exist(name Username) bool {
	u := ul.Get(name)
	return u != nil
}

func (ul *UserList) Add(user *User) error {
	if user == nil {
		return ErrNilUser
	}
	name := user.Name
	if ul.Exist(name) {
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
	if ul.Exist(name) {
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
	users := make(UserListRaw, ul.Size()) // CHECK: diff between `len` and `cap` in golang
	for _, u := range ul.Users {
		users = append(users, u) // CHECK: maybe better to use index in golang ?
	}
	return users
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


var (
    counter  = 0
)

func count() {
    counter += 1
}

func main() {
    ul := NewUserList()
    // ul.Users["a"] = NewUser(UserInfo{"a", "a", "a"})
    ul.Users["a"] = NewUser(UserInfo{})
    ul.Users["aa"] = NewUser(UserInfo{"aa", "aa", "aa"})
    ul.Users["b"] = NewUser(UserInfo{"b", "b", "b"})
    ul.Users["bb"] = NewUser(UserInfo{"bb", "bb", "bb"})

    if err := ul.ForEach(func(key Username)(error) {
        println(key, ul.Users[key])
        defer count()
        return nil
    }); err != nil {
        panic(err)
    }

    println(ul.Size())
{
    u := ul.Get("bb")
    println(u)

    println(ul.Size())
}
    u, _ := ul.PickOut("bb")
    println(u)

    println(ul.Size())
}
