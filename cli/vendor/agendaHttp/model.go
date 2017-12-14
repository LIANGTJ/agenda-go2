package agendaHttp

import (
	// "net/http"
	// "encoding/json"
	// "fmt"
	// log "util/logger"
	"entity"
)

type ErrorResponseBody struct {
	msg string
}

type LoginReqBody struct {

	Username string
	Password string
}

func NewLoginReqBody(username, password string) *LoginReqBody {
	data := LoginReqBody {
		username,
		password,
	}
	return &data
}

func (L * LoginReqBody) Invalid() bool { 
	return L.Username == "" || L.Password == ""
}

type QueryAccountAllResBody struct {
	users queryUserList
	
}

type queryUserList []struct {
	Id string
	entity.User
}

