package agenda

import (
	"agenda/mux"
	"agenda/server"
	"bytes"
	errors "convention/agendaerror"
	"encoding/json"
	"entity"
	"fmt"
	"math/rand"
	"model"
	"net/http"
	"strings"
	"time"
	"util"
	log "util/logger"
)

type Username = entity.Username
type Auth = entity.Auth

// type UserInfo = entity.UserInfo // @@binly:
type UserInfo struct {
	Name  string `json:"username"`
	Auth  string `json:"password"`
	Mail  string `json:"mail"`
	Phone string `json:"phone"`
}

type UserInfoPublic = entity.UserInfoPublic
type User = entity.User
type MeetingInfo = entity.MeetingInfo
type Meeting = entity.Meeting
type MeetingTitle = entity.MeetingTitle

func MakeUserInfo(username Username, password Auth, email, phone string) entity.UserInfo {
	info := entity.UserInfo{}

	info.Name = username
	info.Auth = password
	info.Mail = email
	info.Phone = phone

	return info
}
func MakeMeetingInfo(title MeetingTitle, sponsor Username, participators []Username, startTime, endTime time.Time) entity.MeetingInfo {
	info := entity.MeetingInfo{}

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

// Server ...

const (
	DefaultPort = "8080"
)

var (
	agenda struct {
		*server.Server
	}
)

var logInHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "POST"
}
var logOutHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var getUserKeyHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var getUserByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var deleteUserByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var getMeetingsForUserHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var deleteMeetingsForUserHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var getUsersHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}

type ResponseJSON struct {
	Error   string      `json:"error"`
	Content interface{} `json:"content"`
}

var registerUserHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "POST"
	util.PanicIf(r.Method != "POST")

	var uInfoRaw UserInfo
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := ResponseJSON{Error: "wrong format for elements POST-ed"}
		json.NewEncoder(w).Encode(res)
		return
	}

	uInfo := MakeUserInfo(
		Username(uInfoRaw.Name),
		Auth(uInfoRaw.Auth),
		uInfoRaw.Mail,
		uInfoRaw.Phone,
	)
	switch err := RegisterUser(uInfo); err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		res := ResponseJSON{Content: uInfo}
		json.NewEncoder(w).Encode(res)
		model.UserInfoService.Create(&uInfo)
		return
	case errors.ErrInvalidUsername:
		w.WriteHeader(http.StatusBadRequest)
		res := ResponseJSON{Error: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	case errors.ErrExistedUser:
		w.WriteHeader(http.StatusConflict)
		res := ResponseJSON{Error: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
var getMeetingByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var deleteMeetingByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var modifyMeetingByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "PATCH"
}
var getMeetingByIntervalHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var sponsorMeetingHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "POST"
}

type HTTPMethod = string
type HandlerMap = map[HTTPMethod]http.HandlerFunc

func HandlerMapper(mapping HandlerMap) http.HandlerFunc {
	wantedMethods := ""
	for m := range mapping {
		wantedMethods += "/" + m
	}

	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := mapping[r.Method]
		if ok {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			res := ResponseJSON{
				Error: fmt.Sprintf("wrong request method: %v, however, wanted: %v", r.Method, wantedMethods),
			}
			json.NewEncoder(w).Encode(res)
			return
		}
	}
}

func init() {
	mux := mux.NewServeMux()
	api := "/v1"

	// Group Session

	// mux.HandleFunc(api+"/session", logInHandler) // Method: "POST"
	// mux.HandleFunc(api+"/session", logOutHandler) // Method: "DELETE"
	mux.HandleFunc(api+"/session", HandlerMapper(HandlerMap{
		"POST":   logInHandler,
		"DELETE": logOutHandler,
	}))

	// Group User

	mux.HandleFunc(api+"/user/getkey", getUserKeyHandler) // Method: "GET"

	// mux.HandleFunc(api+"/user/{identifier}", getUserByIDHandler)    // Method: "GET"
	// mux.HandleFunc(api+"/user/{identifier}", deleteUserByIDHandler) // Method: "DELETE"
	mux.HandleFunc(api+"/user/{identifier}", HandlerMapper(HandlerMap{
		"GET":    getUserByIDHandler,
		"DELETE": deleteUserByIDHandler,
	}))

	// mux.HandleFunc(api+"/user/{identifier}/meetings", getMeetingsForUserHandler) // Method: "GET"
	// mux.HandleFunc(api+"/user/{identifier}/meetings", deleteMeetingsForUserHandler) // Method: "DELETE"
	mux.HandleFunc(api+"/user/{identifier}/meetings", HandlerMapper(HandlerMap{
		"GET":    getMeetingsForUserHandler,
		"DELETE": deleteMeetingsForUserHandler,
	}))

	// Group Users
	// mux.HandleFunc(api+"/users", getUsersHandler)     // Method: "GET"
	// mux.HandleFunc(api+"/users", registerUserHandler) // Method: "POST"
	mux.HandleFunc(api+"/users", HandlerMapper(HandlerMap{
		"GET":  getUsersHandler,
		"POST": registerUserHandler,
	}))

	// Group Meeting
	// mux.HandleFunc(api+"/meetings/{identifier}", getMeetingByIDHandler) // Method: "GET"
	// mux.HandleFunc(api+"/meetings/{identifier}", deleteMeetingByIDHandler) // Method: "DELETE"
	// mux.HandleFunc(api+"/meetings/{identifier}", modifyMeetingByIDHandler) // Method: "PATCH"
	mux.HandleFunc(api+"/meetings/{identifier}", HandlerMapper(HandlerMap{
		"GET":    getMeetingByIDHandler,
		"DELETE": deleteMeetingByIDHandler,
		"PATCH":  modifyMeetingByIDHandler,
	}))

	// Group Meetings
	// mux.HandleFunc(api+"/meetings", getMeetingByIntervalHandler) // Method: "GET"
	// mux.HandleFunc(api+"/meetings", sponsorMeetingHandler)       // Method: "POST"
	mux.HandleFunc(api+"/meetings", HandlerMapper(HandlerMap{
		"GET":  getMeetingByIntervalHandler,
		"POST": sponsorMeetingHandler,
	}))

	// ...
	mux.HandleFunc("/api/test", apiTestHandler())
	mux.HandleFunc("/unknown/", sayDeveloping)
	mux.HandleFunc("/say/", sayhelloName)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./asset/"))))

	srv := server.NewServer()
	srv.SetHandler(mux)

	agenda.Server = srv
}

func Listen(addr string) error {
	if addr == "" {
		addr = DefaultPort
	}
	return agenda.Listen(addr)
}

// NOTE: Now, assume the operations' actor are always the `Current User`

// RegisterUser ...
func RegisterUser(uInfo entity.UserInfo) error {
	if !uInfo.Name.Valid() {
		return errors.ErrInvalidUsername
	}

	u := entity.NewUser(uInfo)
	err := entity.GetAllUsersRegistered().Add(u)
	return err
}

/*
func RegisterUser(uInfo UserInfo) error {
	if !uInfo.Name.Valid() {
		return errors.ErrInvalidUsername
	}

	u := entity.NewUser(uInfo)
	err := entity.GetAllUsersRegistered().Add(u)
	return err
}
*/

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
	if info.Sponsor == nil {
		return nil, errors.ErrNilSponsor
	} else if info.Sponsor.Name != LoginedUser().Name {
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

// ...

// detail handlers, etc ... ----------------------------------------------------------------

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	segments := strings.Split(r.URL.Path, "/")
	name := segments[len(segments)-1]
	fmt.Fprintf(w, "Hello %v!\n", name)
}

func sayDeveloping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)

	fmt.Fprintf(w, "Developing!\n")
	fmt.Fprintf(w, "Now NotImplemented!\n")
}

func apiTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "9527", Content: "Hello from Go!\n"}

		// json.NewEncoder(w).Encode(res)
		j, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rand.Seed(time.Now().UnixNano())
		prettyPrint := rand.Float32() < 0.5
		if prettyPrint {
			var out bytes.Buffer
			json.Indent(&out, j, "", "\t")
			j = out.Bytes()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}
