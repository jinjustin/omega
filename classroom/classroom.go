package classroom

import (
	"fmt"
	"encoding/json"
)

//Classroom is struct that use to represent class in Omega.
type Classroom struct{
	ClassID string
	ClassCode string
	ClassName string
	Year string
	Permission string
	Announcement string
	Description string
}

//SetAnnouncement is function that use to set Announcement.
func (c *Classroom) SetAnnouncement(announcement string){
	c.Announcement = announcement
}

//SetDescription is function that use to set Description.
func (c *Classroom) SetDescription(description string){
	c.Description = description
}

// GetClassroomDetail ...
func (c Classroom) GetClassroomDetail() []byte{
	b,err := json.Marshal(c)
	if err != nil{
		fmt.Println(err)
	}
	return b
}

