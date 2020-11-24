package course

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func Test_GetCourseDetail(t *testing.T) {

		//Create course struct
		c := Course{
			CourseCode:"000001",
			CourseID: "01076002",
			CourseName: "Computer Programing",
			Year: "2563",
			Permission: "Private",
		}
		//Expected Result
		expected,err := json.Marshal(c)
		if err != nil{
			//
		}

	t.Run("Unit Test : get Course Detail", func(t *testing.T) {
		//Output
		output := c.GetCourseDetail()
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})
 }

func Test_SetAnnouncement(t *testing.T) {

	//Create course struct
	c := Course{
		CourseCode:"000001",
		CourseID: "01076002",
		CourseName: "Computer Programing",
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

	//Create course struct
	c := Course{
		CourseCode:"000001",
		CourseID: "01076002",
		CourseName: "Computer Programing",
		Year: "2563",
		Permission: "Private",
	}

t.Run("Unit Test : set description", func(t *testing.T) {

	pts := &c

	//input
	input := "This is a course"

	//Expected Result
	expected := "This is a course"

	//Execute function
	pts.SetDescription(input)

	//Compare output to expected output 
	assert.Equal(t,expected,c.Description)
})
}