package answer

import (
	"fmt"
	"encoding/json"
)

//Answer is struct that use to represent answer of the student in that test.
type Answer struct{
	TestID string `json:"testID"`
	StudentID string `json:"studentID"`
	StudentAnswer []Info `json:"studentAnswer"`
	TotalScore string `json:"totalScore"`
	CheckedAnswer string `json:"checkedAnswer"`
	CompletePercent string `json:"completePercent"`
}

//Info is struct that use to store answer of the student.
type Info struct{
	QuestionID string `json:"questionID"`
	Data string `json:"data"`
	QuestionType string `json:"type"`
	Answer []string `json:"answer"`
	MaxScore string `json:"maxscore"`
	Score string `json:"score"`
}

type StudentAnswerInformation struct{
	Firstname string `json:"firstname"`
	Surname string `json:"surname"`
	StudentID string `json:"studentID"`
	CompletePercent string `json:"completePercent"`
}

type StatisticValue struct{
	Max string `json:"max"`
	Min string `json:"min"`
	Mean string `json:"mean"`
	SD string `json:"sd"`
}

type StudentScore struct{
	TotalScore string `json:"totalScore"`
	TestID string `json:"testID"`
	Topic string `json:"topic"`
	Max string `json:"max"`
	Min string `json:"min"`
	Mean string `json:"mean"`
	SD string `json:"sd"`
	CourseID string `json:"courseID"`
	CourseName string `json:"courseName"`
}

// GetAnswerDetail is a function that use to get choice detail in JSON form.
func (a Answer) GetAnswerDetail() []byte{
	b,err := json.Marshal(a)
	if err != nil{
		fmt.Println(err)
	}
	return b
}