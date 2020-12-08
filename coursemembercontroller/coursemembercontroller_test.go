package coursemembercontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"omega/database"
	"omega/student"
	"omega/teacher"
	"database/sql"
)

func Test_AddStudentToCourse(t *testing.T){
	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Input
		studentID := "60010135"
		courseCode := "8BE0E6"

		//Expected Output
		expected := student.Student{
			UserID: "",
			StudentID: "60010135",
			Username: "",
			Password: "",
			Firstname: "Jirakit",
			Surname: "Jitpenthom",
			Email: "",
		}

		//Output
		var output student.Student
		json.Unmarshal(AddStudentToCourse(studentID,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	t.Run("Unit Test 002: Repetitive Student", func(t *testing.T) {

		//Input
		studentID := "60010135"
		courseCode := "8BE0E6"

		//Expected Output
		expected := student.Student{
			UserID: "Can't Invite this student",
			StudentID: "",
			Username: "",
			Password: "",
			Firstname: "",
			Surname: "",
			Email: "",
		}

		//Output
		var output student.Student
		json.Unmarshal(AddStudentToCourse(studentID,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TOMSON")
	if err != nil {
	panic(err)
	}
}

func Test_AddTeacherToCourse(t *testing.T){
	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Input
		userID := "TECH01"
		courseCode := "8BE0E6"

		//Expected Output
		expected := teacher.Teacher{
			UserID: "",
			Username: "",
			Password: "",
			Firstname: "Test",
			Surname: "Teacher",
			Email: "",
		}

		//Output
		var output teacher.Teacher
		json.Unmarshal(AddTeacherToCourse(userID,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	t.Run("Unit Test 002: Repetitive Student", func(t *testing.T) {

		//Input
		userID := "TECH01"
		courseCode := "8BE0E6"

		//Expected Output
		expected := teacher.Teacher{
			UserID: "Can't Join this course",
			Username: "",
			Password: "",
			Firstname: "",
			Surname: "",
			Email: "",
		}

		//Output
		var output teacher.Teacher
		json.Unmarshal(AddTeacherToCourse(userID,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TECH01")
	if err != nil {
	panic(err)
	}
}