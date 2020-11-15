package main

import (
	"omega/classroom"
	"fmt"
	"omega/classroomcreatorcontroller"
	"omega/testdb"
	//"omega/classroomdeletercontroller"
	"omega/classroomlistcontroller"
	//"encoding/json"
)

func main() {

	/*var output classroom.Classroom
	json.Unmarshal(classroomcreatorcontroller.CreateNewClass("Computer Programing","01076002","2563","Private","a2yQbN"),&output)
	println(output.ClassName)*/

	//fmt.Println(string(classroomcreatorcontroller.CreateNewClass("Computer Programing","01076002","2563","Private","a2yQbN")))
	fmt.Println(string(classroomcreatorcontroller.CreateNewClass("Computer Programing","01076002","2564","Private","a2yQbN")))

	c := classroom.Classroom{
		ClassID: "vB2p9U",
		ClassCode: "999999999",
		ClassName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u := testdb.UserInClass{
		UserID: "vBUp15",
		ClassID: "vB2p9U",
	}

	testdb.ClassroomDB = append(testdb.ClassroomDB,c)
	testdb.UserInClassDB = append(testdb.UserInClassDB,u)
	/*fmt.Println(testdb.ClassroomDB)

	fmt.Println(string(classroomdeletercontroller.DeleteClassroom("vB2p9U")))*/
	
	fmt.Println(testdb.ClassroomDB)
	fmt.Println(testdb.UserInClassDB)

	fmt.Println(classroomlistcontroller.GetClassroomList("a2yQbN"))
}