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

	userid := "TOMSON"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO student (userid, studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, userid, `60010135`,`testaddstudent` ,`tocourse`, `student@email.com`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Input
		studentID := "60010135"
		courseCode := "8BE0E6"

		//Expected Output
		expected := student.Student{
			UserID: "",
			StudentID: "60010135",
			Firstname: "testaddstudent",
			Surname: "tocourse",
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

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid)
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM student WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid)
	if err != nil {
	panic(err)
	}
}

func Test_AddTeacherToCourse(t *testing.T){

	username := "testaddteacher"
	userID := "TECH01"
	courseCode := "8BE0E6"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, userID, username,`123456` ,`teacher`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, userID, `testadd`,`teacher` ,`teacher@email.com`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Expected Output
		expected := teacher.Teacher{
			UserID: "",
			Firstname: "testadd",
			Surname: "teacher",
			Email: "",
		}

		//Output
		var output teacher.Teacher
		json.Unmarshal(AddTeacherToCourse(username,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	t.Run("Unit Test 002: Repetitive Teacher", func(t *testing.T) {

		//Expected Output
		expected := teacher.Teacher{
			UserID: "Can't Join this course",
			Firstname: "",
			Surname: "",
			Email: "",
		}

		//Output
		var output teacher.Teacher
		json.Unmarshal(AddTeacherToCourse(username,courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TECH01")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TECH01")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TECH01")
	if err != nil {
	panic(err)
	}
}

func Test_AddApproveJoinCourse(t *testing.T){

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1,$2,$3,$4);`

	_, err = db.Exec(sqlStatement,"COURSE", "TOMSON", "student","pending")
	if err != nil {
	panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1,$2,$3,$4);`

	_, err = db.Exec(sqlStatement,"COURSE", "TECH01", "teacher","pending")
	if err != nil {
	panic(err)
	}

	t.Run("Unit Test 001: Approve Student Join Course", func(t *testing.T) {

		//Input
		userID := "TOMSON"
		courseCode := "COURSE"

		//Expected Output
		expected := "join"

		//Output
		var output string

		//Execute Function
		ApproveJoinCourse(userID,courseCode)

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `SELECT status FROM coursemember WHERE coursecode=$1 and userid=$2;`
		rows,err := db.Query(sqlStatement, courseCode,userID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&output)
				if err != nil {
					panic(err)
				}
			}

		err = rows.Err()
		if err != nil {
		panic(err)
		}
		
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	t.Run("Unit Test 002: Approve Teacher Join Course", func(t *testing.T) {

		//Input
		userID := "TECH01"
		courseCode := "COURSE"

		//Expected Output
		expected := "join"

		//Output
		var output string

		//Execute Function
		ApproveJoinCourse(userID,courseCode)

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `SELECT status FROM coursemember WHERE coursecode=$1 and userid=$2;`
		rows,err := db.Query(sqlStatement, courseCode,userID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&output)
				if err != nil {
					panic(err)
				}
			}

		err = rows.Err()
		if err != nil {
		panic(err)
		}
		
		//Compare output to expected output 
		assert.Equal(t,expected,output)

	})

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TECH01")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TOMSON")
	if err != nil {
	panic(err)
	}
}

func Test_AddDeclineJoinCourse(t *testing.T){

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1,$2,$3,$4);`

	_, err = db.Exec(sqlStatement,"COURSE", "TOMSON", "student","pending")
	if err != nil {
	panic(err)
	}

	t.Run("Unit Test 001: Decline Join Course", func(t *testing.T) {

		//Input
		userID := "TOMSON"
		courseCode := "COURSE"

		//Expected Output
		expected := "success"

		//Output
		var output string

		//Execute Function
		output = DeclineJoinCourse(userID,courseCode)
		
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})
}