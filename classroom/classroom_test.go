package classroom

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func Test_getClassroomDetail(t *testing.T) {

		//Create classroom struct
		c := Classroom{
		ClassID: "000001",
		ClassCode: "01076002",
		ClassName: "Computer Programing",
		Year: "2563",
		Permission: "Private",
		}
		//Expected Result
		expected,err := json.Marshal(c)
		if err != nil{
			//
		}

	t.Run("Unit Test : get Classroom Detail", func(t *testing.T) {
		//Output
		var output Classroom
		json.Unmarshal(c.GetClassroomDetail(),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})
 }