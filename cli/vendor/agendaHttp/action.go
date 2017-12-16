package agendaHttp

import (
	"encoding/json"
	"net/http"
	log "util/logger"
	"fmt"
	"bytes"
	"entity"
	"errors"
	"io/ioutil"
	"status"
	"time"
	"config"
	// "io"
	// "os"
	
)
// ---------------------------- Cmd Function ---------------------------------------


func Register(username, password, email, phone string) (*RegisterResBody, error) {
	defer func(){
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	
	// user := entity.NewUser(username, password, email, phone)
	user := NewRegisterReqBody(username, password, email, phone)
	if user.Invalid() {
		err := errors.New("user regiestered invalid")
		return nil, err
	}
	jsonData := ToJson(*user)

	resp, err := http.Post(RegisterURL(), "application/json", jsonData)
	if err != nil {
		fmt.Println("before")
		panic(err)
	}
	if resp.Status[0] != '2' {
		return nil,ErrorHandle(resp)
	}
	// whether in debug mode
	
	defer resp.Body.Close() //一定要关闭resp.Body

	// var u 
	var u RegisterResBody
	err = json.NewDecoder(resp.Body).Decode(&u)
	if err != nil {
		fmt.Println("after")
		panic(err)
	}

	actWhenInDebugMode(resp, "[Register]")
	return &u, nil

}

func Login(username, password string)  error {
	defer func(){
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	if status.UserExisted() {
		return  errors.New("another curUser " + username + " has logined")
	}

	loginReqBody := NewLoginReqBody(username, password)
	if loginReqBody.Invalid() {
		err := errors.New("user regiestered invalid")
		return err
	}
	jsonData := ToJson(&loginReqBody)
	resp, err := http.Post(LoginURL(), "application/json", jsonData)
	if err != nil {
		panic(err)
	}
	if resp.Status[0] != '2' {
		return ErrorHandle(resp)
	}
	defer resp.Body.Close()
	var body LoginResBody
	json.NewDecoder(resp.Body).Decode(&body)
	status.ChangeLoginedToken(body.Token)

	actWhenInDebugMode(resp, "[Login]")
	
	return nil

}


func Logout() error {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	
	if !status.UserExisted() {
		return errors.New("User Existed")
	}
	client := &http.Client{}
	data := NewLogoutReqBody(status.LoginedToken())
	jsonData := ToJson(*data)
	
	req, err := http.NewRequest(http.MethodDelete, LogoutURL(), jsonData)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	
	if resp.Status[0] != '2' {
		ErrorHandle(resp)
	} else {
		status.ChangeLoginedToken("")
	}

	// actWhenInDebugMode(resp, "[Logout]")
	return err
	
}

func QueryAccountAll() (*QueryAccountAllResBody, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	
	client := &http.Client{}
	data := NewQueryAccountAllReqBody(status.LoginedToken())
	jsonData := ToJson(data)
	req, err := http.NewRequest(http.MethodGet, QueryAccountAllURL(), jsonData)
	// resp, err := http.Get(QueryAccountAllURL())
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.Status[0] != '2' {
		return nil,ErrorHandle(resp)
	}
	defer resp.Body.Close()

	var userlist QueryAccountAllResBody

	err = json.NewDecoder(resp.Body).Decode(&userlist)
	if err != nil {
		panic(err)
	}

	actWhenInDebugMode(resp, "[QueryAccountAll]")
	return &userlist, nil

	
}

func CreateMeeting(title string, participators []string, startTime, endTime time.Time) (*entity.Meeting, error) {
	
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	meeting, err := entity.NewMeeting(title, participators, startTime, endTime)
	if err != nil {
		return nil, err
	}

	jsonData := ToJson(*(meeting.Serialized()))
	resp, err := http.Post(CreateMeetingURL(), "application/json", jsonData)
	if err != nil {
		panic(err)
	}
	if resp.Status[0] != '2' {
		return nil, ErrorHandle(resp)
	}
	defer resp.Body.Close()
	var meetingInfo entity.Meeting
	err = json.NewDecoder(resp.Body).Decode(&meetingInfo)
	if err != nil {
		panic(err)
	}
	fmt.Print(meetingInfo)
	actWhenInDebugMode(resp, "[CreateMeeting]")
	return &meetingInfo, err
	
}

// ---------------------------- Util Function ---------------------------------------

func ToJson(v interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(v)
	return buf
}

func ErrorHandle(resp *http.Response) error {
	if(resp.Status[0] != 2) {

		// var body ErrorResponseBody
		// err := json.NewDecoder(res.Body).Decode(&body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// return errors.New(body.msg)
		// var msg = make([]byte,0,0)
		// io.ReadFull(res.Body,msg)
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(msg))

	}
	return nil

}

func actWhenInDebugMode(resp *http.Response, cmd string) {
	if config.DeBugMode() {
		data, _ := ioutil.ReadAll(resp.Body)
		// ...已经没多大作用了，毕竟前面resp都被decode了
		fmt.Println(cmd + " Response: ", string(data))
	}
}