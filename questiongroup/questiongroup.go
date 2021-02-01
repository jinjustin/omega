package questiongroup

import (
	"fmt"
	"encoding/json"
)

//QuestionGroup is struct that use to represent question group in test bank.
type QuestionGroup struct{
	QuestionGroupID string
	Name string
	CourseID string
	Type string
}

// GetQuestionGroupDetail is a function that use to get question group Detail in JSON form.
func (g QuestionGroup) GetQuestionGroupDetail() []byte{
	b,err := json.Marshal(g)
	if err != nil{
		fmt.Println(err)
	}
	return b
}