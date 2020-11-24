package course

import (
	"fmt"
	"encoding/json"
)

//Course is struct that use to represent class in Omega.
type Course struct{
	CourseCode string
	CourseID string
	CourseName string
	Year string
	Permission string
	Announcement string
	Description string
}

//SetAnnouncement is function that use to set Announcement.
func (c *Course) SetAnnouncement(announcement string){
	c.Announcement = announcement
}

//SetDescription is function that use to set Description.
func (c *Course) SetDescription(description string){
	c.Description = description
}

// GetCourseDetail is function that use to get course Detail in JSON form.
func (c Course) GetCourseDetail() []byte{
	b,err := json.Marshal(c)
	if err != nil{
		fmt.Println(err)
	}
	return b
}

