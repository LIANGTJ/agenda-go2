package agendaHttp

import (
	// "net/http"
	// "encoding/json"
	// "fmt"
	// log "util/logger"
	"entity"
)

type withTokenBody struct {
	Token string `json:"token"`
}


// type ErrorResponseBody struct {
// 	msg string
// }

type LoginReqBody struct {

	Username string	`json:"username"`
	Password string	`json:"password"`
}

type LoginResBody struct {
	withTokenBody
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

type LogoutReqBody struct {
	withTokenBody
}

func NewLogoutReqBody(token string) *LogoutReqBody{
	return &LogoutReqBody{
		withTokenBody{
			token,
		},
	}
}

// type QueryAccountAllResBody struct {
	// users queryUserList
	// users
	
// }

type QueryAccountAllReqBody struct {
	withTokenBody
}

func NewQueryAccountAllReqBody(token string) *QueryAccountAllReqBody{
	return &QueryAccountAllReqBody{
		withTokenBody{
			token,
		},
		
	}
}


type QueryAccountAllResBody []struct {
	entity.User
}

type RegisterReq struct {
	entity.User
}

func NewRegisterReqBody(username, password, email, phone string) *RegisterReq {
	return &RegisterReq{
		entity.User{
			username,
			password,
			email,
			phone,
		},
	}
}

func (r *RegisterReq) Invalid() bool{
	return r.User.Invalid()
}

type RegisterResBody struct {
	entity.User
}

