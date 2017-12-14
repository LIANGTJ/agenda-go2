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
	// "os"
	
)
// ---------------------------- Cmd Function ---------------------------------------

// func Get(url string) *json.Decoder{
// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	defer res.Body.Close()

// 	return json.NewDecoder(res.Body)

// }



func Register(username, password, email, phone string) (*json.Decoder, error) {
	defer func(){
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	
	user := entity.NewUser(username, password, email, phone)
	if user.Invalid() {
		err := errors.New("user regiestered invalid")
		return nil, err
	}
	jsonData := ToJson(*user)
	registerURL := RegisterURL()

	resp, err := http.Post(registerURL, "application/json", jsonData)
	if err != nil {
		panic(err)
	}
	if resp.Status[0] != '2' {
		return nil,ErrorHandle(resp)
	}
	// whether in debug mode
	actWhenInDebugMode(resp, "[Register]")
	// defer resp.Body.Close() //一定要关闭resp.Body
	return json.NewDecoder(resp.Body), nil

}

func Login(username, password string) (*json.Decoder, error) {
	defer func(){
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	if status.AnotherUserExisted() {
		return nil, errors.New("another curUser has logined")
	}

	loginReqBody := NewLoginReqBody(username, password)
	if loginReqBody.Invalid() {
		err := errors.New("user regiestered invalid")
		return nil, err
	}
	jsonData := ToJson(&loginReqBody)
	loginURL := LoginURL()
	resp, err := http.Post(loginURL, "application/json", jsonData)
	if err != nil {
		panic(err)
	}
	if resp.Status[0] != '2' {
		return nil, ErrorHandle(resp)
	}

	status.AddLoginedUser(username)

	actWhenInDebugMode(resp, "[Login]")

	return json.NewDecoder(resp.Body), nil

}

// ---------------------------- Util Function ---------------------------------------

func ToJson(v interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(v)
	return buf
}

func ErrorHandle(res *http.Response) error {
	if(res.Status[0] != 2) {

		var body ErrorResponseBody
		err := json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			log.Fatal(err)
		}
		return errors.New(body.msg)
	}
	return nil

}

func actWhenInDebugMode(resp *http.Response, cmd string) {
	if status.DeBugMode() {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(cmd + " Response: ", string(data))
	}
}