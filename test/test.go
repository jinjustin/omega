package test

import (
	"fmt"
	"encoding/json"
)

//Test is struct that use to represent class in Omega.
type Test struct{
	TestID string
	CourseCode string
	Status string
	Topic string
	Duration string
	Timestart string
	Datestart string
	Description string
	Situation string
}

type ForStudent struct{
	TestID string
	CourseCode string
	CourseID string
	Status string
	Topic string
	Duration string
	Timestart string
	Datestart string
	Description string
	Situation string
}

type FinishTest struct{
	TestID string `json:"testID"`
	Topic string `json:"topic"`
	Paticipant string `json:"paticipant"`
	Process string `json:"process"`
}

// GetTestDetail is function that use to get test Detail in JSON form.
func (t Test) GetTestDetail() []byte{
	b,err := json.Marshal(t)
	if err != nil{
		fmt.Println(err)
	}
	return b
}