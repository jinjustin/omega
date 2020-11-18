package classroom

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func Test_GetClassroomDetail(t *testing.T) {

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
		output := c.GetClassroomDetail()
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})
 }

func Test_SetAnnouncement(t *testing.T) {

	//Create classroom struct
	c := Classroom{
	ClassID: "000001",
	ClassCode: "01076002",
	ClassName: "Computer Programing",
	Year: "2563",
	Permission: "Private",
	}

t.Run("Unit Test : set Announcement", func(t *testing.T) {

	pts := &c

	//input
	input := "Hello World,Hello People"

	//Expected Result
	expected := "Hello World,Hello People"

	//Execute function
	pts.SetAnnouncement(input)

	//Compare output to expected output 
	assert.Equal(t,expected,c.Announcement)
})
}

func Test_SetDescription(t *testing.T) {

	//Create classroom struct
	c := Classroom{
	ClassID: "000001",
	ClassCode: "01076002",
	ClassName: "Computer Programing",
	Year: "2563",
	Permission: "Private",
	}

t.Run("Unit Test : set Description", func(t *testing.T) {

	pts := &c

	//input
	input := "This is a class"

	//Expected Result
	expected := "This is a class"

	//Execute function
	pts.SetDescription(input)

	//Compare output to expected output 
	assert.Equal(t,expected,c.Description)
})
}