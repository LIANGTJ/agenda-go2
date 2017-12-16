package status

import (
	"config"
	"io/ioutil"
	"os"
	// "errors"
	log "util/logger"
)

// var loginedUser = ""
var loginedToken = ""


func UserExisted() bool { 
	// return loginedUser != "" 
	return loginedToken != ""
	
}



// func LoginedUser() string { 
	// return loginedUser 

// }

func LoginedToken() string{
	return loginedToken
}
// func ChangeLoginedUser(user string) { 
// 	loginedUser = user 
// }

func ChangeLoginedToken(token string) {
	loginedToken = token
}

func LoadLoginStatus()  {
	buf, err := ioutil.ReadFile(config.UserLoginedStatusPath())
	// fin, err := os.Open(config.UserLoginStatusPath())
	// defer fin.Close()
	if err != nil {
		log.Fatal(err)
	}
	// loginedUser = string(buf)
	loginedToken = string(buf)
	// fin.Read()
	// return errors.ErrNeedImplement
}

func SaveLoginStatus()  {
	fout, err := os.Create(config.UserLoginedStatusPath())
	defer fout.Close()
	if err != nil {
		log.Fatal(err)
	}
	// fout.WriteString(loginedUser)
	fout.WriteString(loginedToken)

}
