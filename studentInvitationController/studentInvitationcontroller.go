package studentinvitationcontroller

import (
	//"encoding/json"
	"omega/testdb"
	"omega/student"
	//"fmt"
)

//InviteStudentMock is Use to invite student to classroom
func InviteStudentMock(studentID string, classID string) []byte {

	var userID string

	var userInClass testdb.UserInClass

	if(checkStudentID(studentID,classID) == true){

		for _, a := range testdb.StudentDB {
			if a.StudentID == studentID {
				userID = a.UserID
			}
		}

		if(checkUserID(userID,classID) == true){
			//
		}

	}

}

//CheckStudentID is Class that use to check correctness of studentID that user input
func checkStudentID(studentID string, classID string) bool {
	return len(studentID) == 8
}

func checkUserID(userID string,classID string) bool{
	correct := true

	for _, a := range testdb.UserInClassDB {
		if a.UserID == userID && a.ClassID == classID {
			correct = false
		}
	}

	return correct
}