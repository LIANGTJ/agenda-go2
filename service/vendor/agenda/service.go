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

type RequestJSON struct {
	Token entity.Token `json:"token"`
	UserInfoRaw
	// ...
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

var logInHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "POST")

	var uInfoRaw UserInfoRaw
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		RespondErrorDecoding(w, err)
		return
	}

	name := Username(uInfoRaw.Name)
	uInfo, err := QueryAccountByUsername(name)
	if err != nil {
		RespondError(w, err)
		return
	}

	// LogIn(name, authTrial)
	authTrial := Auth(uInfoRaw.Auth)
	if !uInfo.Auth.Verify(authTrial) {
		RespondError(w, errors.ErrFailedAuth)
	} else {
		expire := time.Now().Add(10 * time.Minute)
		// cookie := http.Cookie{"test", "tcookie", "/", "www.domain.com", expire, expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}}
		// http.SetCookie(w, &cookie)

		sInfo := entity.SessionInfo{
			ExpiredAt: expire,
			User:      uInfo,
		}
		if err := CreateSession(&sInfo); err != nil {
			RespondError(w, err)
			return
		}
		RespondJSON(w, http.StatusOK, ResponseToken{sInfo.Token})
	}
}
var logOutHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "DELETE")

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	sInfo, err := Authorize(rInfo.Token)
	if err != nil {
		RespondError(w, err)
		return
	}

	if err := DeleteSession(&sInfo); err != nil {
		RespondError(w, err)
		return
	}

	// RespondJSON(w, http.StatusNoContent)
	// RespondError(w, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}
var getUserKeyHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "GET"
}
var getUserByIDHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	if _, err := Authorize(rInfo.Token); err != nil {
		RespondError(w, err)
		return
	}

	if us := muxx.Vars(r)["identifier"]; len(us) > 0 { // FIXME: used muxx
		// if us := r.URL.Query()["username"]; len(us) > 0 {
		// name := Username(us[0])
		name := Username(us)
		uInfo, err := QueryAccountByUsername(name)
		if err != nil {
			RespondError(w, err)
			return
		}

		res := ResponseUserInfoPublic(uInfo.UserInfoPublic)
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

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	if _, err := Authorize(rInfo.Token); err != nil {
		RespondError(w, err)
		return
	}

	// uInfos := QueryAccountAll()
	if uInfos, err := model.UserInfoService.FindAll(); err != nil {
		RespondError(w, err)
	} else {
		res := make([]entity.UserInfoPublic, 0, len(uInfos))
		for _, u := range uInfos {
			res = append(res, u.UserInfoPublic)
		}
		RespondJSON(w, http.StatusOK, res)
	}
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

	res := ResponseUserInfoPublic(uInfo.UserInfoPublic)
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

func init() {
	// FIXME: when use `curl` and no-trail-slash url to test, fail to be redirected to with-trail-slash version like when using browser .... whatever mux or muxx
	// when using muxx, seems not redirect sub-tree (like '/users/a' --> '/users/') ...
	// mux := mux.NewServeMux()
	mux := muxx.NewRouter() // TODO: replace `mux` ?
	api := "/v1"

	// Group Session
	mux.HandleFunc(api+"/sessions/", HandlerMapper(HandlerMap{
		"POST": logInHandler,
	}))
	mux.HandleFunc(api+"/session", HandlerMapper(HandlerMap{
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
