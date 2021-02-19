package test

import (
	"fmt"
	"encoding/json"
)

//Test is struct that use to represent class in Omega.
type Test struct{
	TestID string
	CourseID string
	//CourseCode string
	//UserID string
	Status string
	Topic string
	Duration string
	Timestart string
	Datestart string
	Description string
}

// GetTestDetail is function that use to get test Detail in JSON form.
func (t Test) GetTestDetail() []byte{
	b,err := json.Marshal(t)
	if err != nil{
		fmt.Println(err)
	}
	return b
}