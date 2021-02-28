package question

import (
	"fmt"
	"encoding/json"
)

//Question is struct that use to represent question in question group.
type Question struct{
	TestID string
	GroupID string
	Question string
	QuestionID string
	QuestionType string
	Data string
}

// GetQuestionDetail is a function that use to get question detail in JSON form.
func (q Question) GetQuestionDetail() []byte{
	b,err := json.Marshal(q)
	if err != nil{
		fmt.Println(err)
	}
	return b
}