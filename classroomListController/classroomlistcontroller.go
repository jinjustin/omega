package classroomlistcontroller

import (
	"omega/classroom"
	"omega/testdb"
)


//GetClassroomList is use to get all classrooms that user is being member.
func GetClassroomList(userID string) []classroom.Classroom{

	var classIDs []string

	var classes []classroom.Classroom

	for _, a := range testdb.UserInClassDB {
		if a.UserID == userID {
			classIDs = append(classIDs, a.ClassID)  
		}
	}
	 
	for _, a := range classIDs {
		for _, b := range testdb.ClassroomDB {
			if b.ClassID == a {
				classes = append(classes, b)  
			}
		}
	}

	return classes
}