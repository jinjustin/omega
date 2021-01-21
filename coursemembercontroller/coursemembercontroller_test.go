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
		json.Unmarshal(addStudentToCourse(studentID,courseCode),&output)
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
		json.Unmarshal(addStudentToCourse(studentID,courseCode),&output)
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
		json.Unmarshal(addTeacherToCourse(username,courseCode),&output)
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
		json.Unmarshal(addTeacherToCourse(username,courseCode),&output)
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

func Test_ApproveJoinCourse(t *testing.T){

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

func Test_DeclineJoinCourse(t *testing.T){

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

func Test_GetStudentInCourse(t *testing.T){

	type studentInCourse struct{
		StudentID string
		Firstname string
		Surname string
		Status string
	}

	studentInCourse1 := studentInCourse{
		StudentID: "99010135",
		Firstname: "testgetstudentincourse",
		Surname: "one",
		Status: "join",
	}

	studentInCourse2 := studentInCourse{
		StudentID: "99010136",
		Firstname: "testgetstudentincourse",
		Surname: "two",
		Status: "pending",
	}

	userid1 := "AAAAAA"

	userid2 := "AAAAAB"

	courseCode := "OOOOOO"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO student (userid, studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, userid1, studentInCourse1.StudentID, studentInCourse1.Firstname, studentInCourse1.Surname, `student1@email.com`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO student (userid, studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, userid2, studentInCourse2.StudentID, studentInCourse2.Firstname, studentInCourse2.Surname, `student2@email.com`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, userid1, `student`, studentInCourse1.Status)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, userid2, `student`, studentInCourse2.Status)
	if err != nil {
		panic(err)
	}

	var expected []studentInCourse

	expected = append(expected, studentInCourse1)
	expected = append(expected, studentInCourse2)

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Output
		var output []studentInCourse
		json.Unmarshal(getStudentInCourse(courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	sqlStatement = `DELETE FROM student WHERE firstname=$1;`
	_, err = db.Exec(sqlStatement, "testgetstudentincourse")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid1)
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid2)
	if err != nil {
	panic(err)
	}
}

func Test_GetTeacherInCourse(t *testing.T){

	type teacherInCourse struct{
		Firstname string
		Surname string
		Status string
	}

	teacherInCourse1 := teacherInCourse{
		Firstname: "testgetteacherincourse",
		Surname: "one",
		Status: "join",
	}

	teacherInCourse2 := teacherInCourse{
		Firstname: "testgetteacherincourse",
		Surname: "two",
		Status: "pending",
	}

	userid1 := "AAAAAA"

	userid2 := "AAAAAB"

	courseCode := "OOOOOO"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, userid1, teacherInCourse1.Firstname, teacherInCourse1.Surname, `teacher1@email.com`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, userid2, teacherInCourse2.Firstname, teacherInCourse2.Surname, `teacher2@email.com`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, userid1, `teacher`, teacherInCourse1.Status)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, userid2, `teacher`, teacherInCourse2.Status)
	if err != nil {
		panic(err)
	}

	var expected []teacherInCourse

	expected = append(expected, teacherInCourse1)
	expected = append(expected, teacherInCourse2)

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Output
		var output []teacherInCourse
		json.Unmarshal(getTeacherInCourse(courseCode),&output)
		//Compare output to expected output 
		assert.Equal(t,expected,output)
	})

	sqlStatement = `DELETE FROM teacher WHERE firstname=$1;`
	_, err = db.Exec(sqlStatement, "testgetteacherincourse")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid1)
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid2)
	if err != nil {
	panic(err)
	}
}

func Test_DeleteTeacherInCourse(t *testing.T){

	teacher1 := teacher.Teacher{
		UserID: "TE0001",
		Firstname: "testdeleteteacherincourse",
		Surname: "one",
		Email: "teacher1@kmitl.ac.th",
	}

	courseCode := "OOOOOO"
	username := "deleteteacher"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, teacher1.UserID, teacher1.Firstname, teacher1.Surname, teacher1.Email)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, teacher1.UserID, username, "123456", "teacher")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, teacher1.UserID, `teacher`, "Join")
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Output
		var output teacher.Teacher
		json.Unmarshal(deleteTeacherInCourse(courseCode,username),&output)
		//Compare output to expected output 
		assert.Equal(t,teacher1.Firstname,output.Firstname)
		assert.Equal(t,teacher1.Surname,output.Surname)
		assert.Equal(t,teacher1.Email,output.Email)
	})

	sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, teacher1.UserID)
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, teacher1.UserID)
	if err != nil {
	panic(err)
	}
}

func Test_DeleteStudentInCourse(t *testing.T){

	student1 := student.Student{
		UserID: "ST0001",
		StudentID: "99010139",
		Firstname: "testdeletestudentincourse",
		Surname: "two",
		Email: "student1@kmitl.ac.th",
	}

	courseCode := "OOOOOO"
	username := "deletestudent"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO student (userid, studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, student1.UserID, student1.StudentID, student1.Firstname, student1.Surname, student1.Email)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, student1.UserID, username, "123456", "student")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, courseCode, student1.UserID, "student", "Join")
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Work Right", func(t *testing.T) {

		//Output
		var output student.Student
		json.Unmarshal(deleteStudentInCourse(courseCode,username),&output)
		//Compare output to expected output
		assert.Equal(t,student1.StudentID,output.StudentID) 
		assert.Equal(t,student1.Firstname,output.Firstname)
		assert.Equal(t,student1.Surname,output.Surname)
		assert.Equal(t,student1.Email,output.Email)
	})

	sqlStatement = `DELETE FROM student WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, student1.UserID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, student1.UserID)
	if err != nil {
		panic(err)
	}
}

func Test_GetUserRole(t *testing.T){

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, "ST0001", "getstudentrole", "123456", "student")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, "TE0001", "getteacherrole", "123456", "teacher")
	if err != nil {
		panic(err)
	}


	t.Run("Unit Test 001: Get Student Role", func(t *testing.T) {

		//Output
		output := getUserRole("getstudentrole")

		//Compare output to expected output
		assert.Equal(t,"student",output) 
		
	})

	t.Run("Unit Test 002: Get Teacher Role", func(t *testing.T) {

		//Output
		output := getUserRole("getteacherrole")

		//Compare output to expected output
		assert.Equal(t,"teacher",output) 
		
	})

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "ST0001")
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TE0001")
	if err != nil {
		panic(err)
	}
}