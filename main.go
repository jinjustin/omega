package main

import (
	//"omega/classroom"
	"fmt"
	//"omega/classroomcreatorcontroller"
	//"omega/database"
	//"omega/classroomdeletercontroller"
	"omega/classroomlistcontroller"
	//"encoding/json"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func main() {
	//fmt.Println(string(classroomcreatorcontroller.CreateNewClass("Test","01076002","2564","Private","a2yQbN")))
	//fmt.Println(string(classroomdeletercontroller.DeleteClassroom("132749","a2yQbN")))
	fmt.Println(classroomlistcontroller.GetClassroomList("a2yQbN"))
}