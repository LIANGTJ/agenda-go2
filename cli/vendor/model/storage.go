package model
import (
	"status"
	
)


func SaveAll() {
	status.SaveLoginStatus()
}

func LoadAll() {
	status.LoadLoginStatus()
}