package choice

import (
	"fmt"
	"encoding/json"
)

//Choice is struct that use to represent choice in question.
type Choice struct{
	ChoiceID string `json:"choiceID"`
	QuestionID string `json:"questionID"`
	Data string `json:"data"`
	ImageLink string `json:"imageLink"`
	Check string `json:"check"`
}


// GetChoiceDetail is a function that use to get choice detail in JSON form.
func (c Choice) GetChoiceDetail() []byte{
	b,err := json.Marshal(c)
	if err != nil{
		fmt.Println(err)
	}
	return b
}