package question

import (
	"fmt"
	"encoding/json"
	"github.com/jinjustin/omega/choice"
)

//Question is struct that use to represent question in question group.
type Question struct{
	TestID string `json:"testID"`
	GroupID string `json:"groupID"`
	QuestionName string `json:"question"`
	QuestionID string `json:"questionID"`
	QuestionType string `json:"type"`
	Data string `json:"data"`
}

//AllQuestionInGroup is struct that use to return all question in group.
type AllQuestionInGroup struct{
	QuestionID string `json:"questionID"`
	QuestionName string `json:"question"`
}

//WithChoice is struct that use to get and return question and choice in test.
type WithChoice struct{
	TestID string `json:"testID"`
	GroupID string `json:"groupID"`
	QuestionName string `json:"question"`
	QuestionType string `json:"type"`
	Data string `json:"data"`
	ChoiceDetail []choice.Choice `json:"choice"`
}

// GetQuestionDetail is a function that use to get question detail in JSON form.
func (q Question) GetQuestionDetail() []byte{
	b,err := json.Marshal(q)
	if err != nil{
		fmt.Println(err)
	}
	return b
}