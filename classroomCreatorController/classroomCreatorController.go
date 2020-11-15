package classroomcreatorcontroller

import (
	"omega/classroom"
	"fmt"
	//"encoding/json"
	"crypto/rand"
	"omega/testdb"
)

//CreateNewClass is function that teacher use to create new class
func CreateNewClass(className string,classCode string,year string,permission string,userID string) []byte{

	check := checkClassCode(classCode,year)
	classID := generateClassID()

	if check == true {
		c := classroom.Classroom{
			ClassID: classID,
			ClassName: className,
			ClassCode: classCode,
			Year: year,
			Permission: permission,
		}

		addCreator := testdb.UserInClass{
			UserID: userID,
			ClassID: classID,
		}

		testdb.ClassroomDB = append(testdb.ClassroomDB,c)
		testdb.UserInClassDB = append(testdb.UserInClassDB,addCreator)
		return c.GetClassroomDetail()
	} else {
		c := classroom.Classroom{
			ClassID: "",
			ClassName: "",
			ClassCode: "Error",
			Year: "",
			Permission: "",
		}
		return c.GetClassroomDetail()
	}
}

//GenerateClassID is function that use to Generate Class ID
func generateClassID() string{
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
        panic(err)
    }
	s := fmt.Sprintf("%X", b)
	return s
}

func checkClassCode(classCode string,year string) bool{
	correct := true
	if len(classCode) != 8{
		correct = false
	}

	for _, a := range classCode {
		if a > '9' && a < '0'{
		   correct = false
		}
	 }

	for _, a := range testdb.ClassroomDB {
		if a.ClassCode == classCode && a.Year == year{
		   correct = false
		}
	 }

	return correct
}