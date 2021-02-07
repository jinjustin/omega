package teacher

import(
	"fmt"
	"encoding/json"
)

//Teacher is struct that use to represent student in this platform.
type Teacher struct{
	UserID string
	Firstname string
	Surname string
	Email string
}

// GetTeacherDetail is function use to get student detail in form of JSON.
func (t Teacher) GetTeacherDetail() []byte{
	b,err := json.Marshal(t)
	if err != nil{
		fmt.Println(err)
	}
	return b
}
