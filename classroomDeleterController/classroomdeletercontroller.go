package classroomdeletercontroller

import (
	"omega/classroom"
	"omega/testdb"
)

//DeleteClassroom is use to delete classroom
func DeleteClassroom(classID string) []byte{
	c := classroom.Classroom{
		ClassID: "Can't find.",
		ClassCode: "",
		ClassName: "",
		Year: "",
		Permission: "",
	}

	index := 0
	for _, a := range testdb.ClassroomDB {
		if a.ClassID == classID {
			c = classroom.Classroom{
				ClassID: a.ClassID,
				ClassName: a.ClassName,
				ClassCode: a.ClassCode,
				Year: a.Year,
				Permission: a.Permission,
			}
			testdb.ClassroomDB = append(testdb.ClassroomDB[:index], testdb.ClassroomDB[index+1:]...)  
		}
		index++
	 }
	 return c.GetClassroomDetail()
}