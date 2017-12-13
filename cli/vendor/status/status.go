package agenda

import (
	"config"
	"io/ioutil"
	"os"
	log "util/logger"
)

var loginedUser = ""

func LoginedUser() string {
	
	return loginedUser

}

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
	// fout.WriteString(string(loginedUser))
	fout.WriteString(loginedUser)

}
