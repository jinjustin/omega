package student

import (
	"fmt"
	"encoding/json"
)

//Student is struct that use to represent student in this platform.
type Student struct{
	UserID string
	StudentID string
	Firstname string
	Surname string
	Email string
}

// GetStudentDetail is function use to get student detail in form of JSON.
func (s Student) GetStudentDetail() []byte{
	b,err := json.Marshal(s)
	if err != nil{
		fmt.Println(err)
	}
	return b
}