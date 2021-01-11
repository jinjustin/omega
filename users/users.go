package users

import (
	"fmt"
	"encoding/json"
)

//Users is struct that use to represent user in this platform both student and teacher.
type Users struct{
	UserID string
	Username string
	Password string
	Role string
}

// GetStudentDetail is function use to get student detail in form of JSON.
func (u Users) GetStudentDetail() []byte{
	b,err := json.Marshal(u)
	if err != nil{
		fmt.Println(err)
	}
	return b
}