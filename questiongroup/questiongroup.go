package questiongroup

import (
	"fmt"
	"encoding/json"
	"github.com/jinjustin/omega/question"
)

//QuestionGroup is struct that use to represent question group in test bank.
type QuestionGroup struct{
	Name string `json:"name"`
	ID string `json:"id"`
	GroupName string `json:"groupName"`
	NumQuestion string `json:"numQuestion"`
	MaxQuestion string `json:"maxQuestion"`
	Score string `json:"score"`
	CourseID string `json:"courseID"`
	TestID string `json:"testID"`
	UUID string `json:"uuid"`
	HeaderOrder int `json:"headerOrder"`
	GroupOrder int `json:"groupOrder"`
}

//GrouptestList is struct that use to return grouptestlist.
type GrouptestList struct{
	ID string `json:"id"` 
	GroupName string `json:"groupName"`
	QuestionList []question.AllQuestionInGroup `json:"questionList"`
}

// GetQuestionGroupDetail is a function that use to get question group Detail in JSON form.
func (g QuestionGroup) GetQuestionGroupDetail() []byte{
	b,err := json.Marshal(g)
	if err != nil{
		fmt.Println(err)
	}
	return b
}