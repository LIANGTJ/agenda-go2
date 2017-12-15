package agenda

import (
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

	muxx "github.com/gorilla/mux"
)

type Username = entity.Username
type Auth = entity.Auth

// type UserInfo = entity.UserInfoRaw
type UserInfoRaw struct {
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
var getUserByIDHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	// print(r.URL.Path, "\n\n")
	// print(r.URL.Query(), "\n\n")

	// id := muxx.Vars(r)["identifier"]
	// print("id: "+id, "\n")

	if us := muxx.Vars(r)["identifier"]; len(us) > 0 { // FIXME: used muxx
		// if us := r.URL.Query()["username"]; len(us) > 0 {
		// name := Username(us[0])
		name := Username(us)
		uInfo, err := QueryAccountByUsername(name)
		if err != nil {
			RespondError(w, err)
			return
		}
		// log.Printf("err: %+v; name: %+v; uInfo: %+v \n", err, name, uInfo)
		res := ResponseJSON{Content: uInfo.UserInfoPublic}
		RespondJSON(w, http.StatusOK, res)
	}
}
var deleteUserByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var getMeetingsForUserHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var deleteMeetingsForUserHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}
var getUsersHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	uInfos := QueryAccountAll()
	res := ResponseJSON{Content: uInfos}
	RespondJSON(w, http.StatusOK, res)
}

var registerUserHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "POST")

	var uInfoRaw UserInfoRaw
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		RespondError(w, http.StatusBadRequest, err.Error(), "decode error for elements POST-ed")
		return
	}

	uInfo := MakeUserInfo(
		Username(uInfoRaw.Name),
		Auth(uInfoRaw.Auth),
		uInfoRaw.Mail,
		uInfoRaw.Phone,
	)
	if err := RegisterUser(uInfo); err != nil {
		RespondError(w, err)
		return
	}

	res := ResponseJSON{Content: uInfo.UserInfoPublic}
	RespondJSON(w, http.StatusCreated, res)
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

type HTTPStatusCode = int

// ErrorOrCode can only hold `error` or `HTTPStatusCode` type
type ErrorOrCode = interface{}

var ErrInvalidMethod = errors.New("invalid request method")

var StatusCodeCorrespondingToAgendaError = map[error]HTTPStatusCode{
	errors.ErrInvalidUsername: http.StatusBadRequest,
	errors.ErrExistedUser:     http.StatusConflict,
	ErrInvalidMethod:          http.StatusBadRequest,
}

func RespondError(w http.ResponseWriter, err ErrorOrCode, msg ...string) {
	errString := strings.Join(msg, "\n")
	errCode := http.StatusInternalServerError

	switch e := err.(type) {
	case error:
		errString = e.Error() + "\n\n" + errString
		code, ok := StatusCodeCorrespondingToAgendaError[e]
		if ok {
			errCode = code
		}
	case HTTPStatusCode:
		errCode = e
	default:
		log.Panicf("type `ErrorOrCode` expects `error` or `HTTPStatusCode`, but not %T", e)
	}

	// NOTE: seems that only using `http.Error` to handle simple error is enough ...
	// w.WriteHeader(code)
	// res := ResponseJSON{Error: errString}
	// json.NewEncoder(w).Encode(res)

	http.Error(w, errString, errCode)
}

type ResponseJSON struct {
	Error   string      `json:"error"`
	Content interface{} `json:"content"`
}

func RespondJSON(w http.ResponseWriter, code HTTPStatusCode, res ResponseJSON) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

type HTTPMethod = string
type HandlerMap = map[HTTPMethod]http.HandlerFunc

func HandlerMapper(mapping HandlerMap) http.HandlerFunc {
	methods := make([]string, 0, len(mapping))
	for m := range mapping {
		methods = append(methods, m)
	}
	wantedMethods := strings.Join(methods, "/")

	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := mapping[r.Method]
		if ok {
			handler(w, r)
		} else {
			RespondError(w, ErrInvalidMethod, fmt.Sprintf("used method: %v, however, wanted: %v", r.Method, wantedMethods))
			return
		}
	}
}

func init() {
	// FIXME: when use `curl` and no-trail-slash url to test, fail to be redirected to with-trail-slash version like when using browser .... whatever mux or muxx
	// when using muxx, seems not redirect sub-tree (like '/users/a' --> '/users/') ...
	// mux := mux.NewServeMux()
	mux := muxx.NewRouter() // TODO: replace `mux` ?
	api := "/v1"

	// Group Session
	mux.HandleFunc(api+"/sessions/", HandlerMapper(HandlerMap{
		"POST":   logInHandler,
		"DELETE": logOutHandler,
	}))

	// Group User
	mux.HandleFunc(api+"/user/getkey", getUserKeyHandler) // Method: "GET" TODEL:

	mux.HandleFunc(api+"/user/{identifier}", HandlerMapper(HandlerMap{
		"GET":    getUserByIDHandler,
		"DELETE": deleteUserByIDHandler,
	}))
	mux.HandleFunc(api+"/user/{identifier}/meetings", HandlerMapper(HandlerMap{
		"GET":    getMeetingsForUserHandler,
		"DELETE": deleteMeetingsForUserHandler,
	}))

	// Group Users
	mux.HandleFunc(api+"/users/", HandlerMapper(HandlerMap{
		"GET":  getUsersHandler,
		"POST": registerUserHandler,
	}))

	// Group Meeting
	mux.HandleFunc(api+"/meetings/{identifier}", HandlerMapper(HandlerMap{
		"GET":    getMeetingByIDHandler,
		"DELETE": deleteMeetingByIDHandler,
		"PATCH":  modifyMeetingByIDHandler,
	}))

	// Group Meetings
	mux.HandleFunc(api+"/meetings/", HandlerMapper(HandlerMap{
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

// ----------------------------------------------------------------
// @@binly: new

func QueryAccountByUsername(name entity.Username) (entity.UserInfo, error) {
	if !name.Valid() {
		return entity.UserInfo{}, errors.ErrInvalidUsername
	}
	uInfo, err := model.UserInfoService.FindByUsername(name)
	return uInfo, err
}
