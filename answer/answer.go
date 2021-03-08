package answer

import (
	"fmt"
	"encoding/json"
)

//Answer is struct that use to represent answer of the student in that test.
type Answer struct{
	TestID string `json:"testID"`
	StudentID string `json:"studentID"`
	StudentAnswer map[string]Info `json:"studentAnswer"`
	TotalScore int `json:"totalScore"`
	CompletePercent int `json:"completePercent"`
}

//Info is struct that use to store answer of the student.
type Info struct{
	QuestionName string
	QuestionType string
	Answer string
	Score string
}


// GetAnswerDetail is a function that use to get choice detail in JSON form.
func (a Answer) GetAnswerDetail() []byte{
	b,err := json.Marshal(a)
	if err != nil{
		fmt.Println(err)
	}
	return b
}