package questiongroup

import (
	"fmt"
	"encoding/json"
)

//QuestionGroup is struct that use to represent question group in test bank.
type QuestionGroup struct{
	Name string
	ID string
	GroupName string
	NumQuestion string
	Score string
	CourseID string
	TestID string
	UUID string
	HeaderOrder int
	GroupOrder int
}

// GetQuestionGroupDetail is a function that use to get question group Detail in JSON form.
func (g QuestionGroup) GetQuestionGroupDetail() []byte{
	b,err := json.Marshal(g)
	if err != nil{
		fmt.Println(err)
	}
	return b
}