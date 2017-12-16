package entity

import (
	"testing"
	// "time"
	// "strings"
)

func Test_UserValid(t *testing.T) {
	tests := make([]User,0,0)
	testCaseNum := 5
	for i := 0; i < testCaseNum; i++ {
		name := "test" + string(i)
		var password string = "psw"
		if i % 2 == 0 {
			password = ""
		}
		email := ""
		phone := ""
		tests = append(tests, *NewUser(name, password, email, phone))
	}
	for index, user := range tests {
		t.Run(user.Username, func(subt * testing.T) {
			if index % 2 == 0 {
				if user.Password != "" {
					subt.Error("test failed, the field password should be null")
				} else {
					subt.Log("test sucessfully, graduation")
				}
			}
		})
	}
}