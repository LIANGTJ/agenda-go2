package status

import (
	"config"
	"io/ioutil"
	"os"
	// "errors"
	log "util/logger"
)

var loginedUser = ""
var debugMode = true

func DeBugMode() bool { return debugMode }
func UserExisted() bool { return loginedUser != "" }

func LoginedUser() string { return loginedUser }
func ChangeLoginedUser(user string) { loginedUser = user }

func LoadLoginStatus()  {
	buf, err := ioutil.ReadFile(config.UserLoginedStatusPath())
	// fin, err := os.Open(config.UserLoginStatusPath())
	// defer fin.Close()
	if err != nil {
		log.Fatal(err)
	}
	loginedUser = string(buf)
	// fin.Read()
	// return errors.ErrNeedImplement
}

func SaveLoginStatus()  {
	fout, err := os.Create(config.UserLoginedStatusPath())
	defer fout.Close()
	if err != nil {
		log.Fatal(err)
	}
	fout.WriteString(loginedUser)

}
