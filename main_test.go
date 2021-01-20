package main

import (
	//"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"database/sql"
	"omega/database"
	"omega/student"
	"omega/teacher"
	"omega/login"

	"encoding/json"
	"omega/course"
)

func Test_createCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `bNMv75`, `createcourseintegrationtest`,`123456` ,`teacher`)
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Create Course", func(t *testing.T) {

		var jsonStr = `{"CourseName": "Test Test","CourseID": "55555555","Year": "2563","Permission": "Public", "Announcement": "Course announcement here.","Description": "Course description here.","Username": "createcourseintegrationtest"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/createcourse", createCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/createcourse","application/json",body)
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)
	
		//เก็บค่าตอบกลับไว้ใน course
		var course course.Course
		json.Unmarshal(output, &course)

		fmt.Print(course)
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"Test Test",course.CourseName)
		assert.Equal(t,"55555555",course.CourseID)
	})

	sqlStatement = `DELETE FROM course WHERE courseid=$1;`
	_, err = db.Exec(sqlStatement, "55555555")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "bNMv75")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM users WHERE userid=$1;`
   _, err = db.Exec(sqlStatement, "bNMv75")
   if err != nil {
	   panic(err)
  }
}

func Test_deleteCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid,username,password,role)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "userID", "deletecourseintegrationtest", "123456", "teacher");
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "123456", "userID", "teacher", "join");
	if err != nil {
		panic(err)
	}


	t.Run("Integration Test: Delete Course", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "123456", "Username": "deletecourseintegrationtest"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/deletecourse",deleteCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/deletecourse","application/json",body)
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)
	
		//เก็บค่าตอบกลับไว้ใน course
		var course course.Course
		json.Unmarshal(output, &course)
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"Test Test",course.CourseName)
		assert.Equal(t,"55555555",course.CourseID)
	})

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "userID")
	if err != nil {
		panic(err)
   }
}

func Test_getCourseList(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`

		_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public")
		if err != nil {
		panic(err)
		}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "123456", "userID", "teacher", "join");
	if err != nil {
	panic(err)
	}

	sqlStatement = `INSERT INTO users (userid,username,password,role)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "userID", "getcourselistintegrationtest", "123456", "teacher");
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Get Course List", func(t *testing.T) {

		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/login", login.Login).Methods("GET")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Get(ts.URL + "/login?username=getcourselistintegrationtest&password=123456")
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		loginRes, err := ioutil.ReadAll(resp.Body)

		var token login.Token

		json.Unmarshal(loginRes, &token)

		var jsonStr = `{"UserName": "getcourselistintegrationtest"}`

		body := strings.NewReader(jsonStr)

		r.Handle("/getcourselist", getCourseList).Methods("POST")
		ts = httptest.NewServer(r)
		defer ts.Close()
	
		resp, err = http.Post(ts.URL + "/getcourselist","application/json",body)
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)
		
		//fmt.Println(string(output))

		//เก็บค่าตอบกลับไว้ใน course
		var courses []course.Course
		json.Unmarshal(output, &courses)

		var expected []course.Course 

		expected = append(expected,course.Course{
			CourseCode: "123456",
			CourseID: "55555555",
			CourseName: "Test Test",
			Year: "2563",
			Permission: "Public",})
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,courses)
	})

	sqlStatement = `DELETE FROM course WHERE courseid=$1;`
	_, err = db.Exec(sqlStatement, "55555555")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "userID")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM users WHERE userid=$1;`
   _, err = db.Exec(sqlStatement, "userID")
   if err != nil {
	   panic(err)
  }
}

func Test_getTeacherInfo(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "AAAA99", "testgetteacherinfo", "123456", "teacher")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "AAAA99", "Myname", "Mysurname", "teacher@kmitl.ac.th");
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Get Teacher Info", func(t *testing.T) {
		var jsonStr = `{"Username": "testgetteacherinfo"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getteacherinfo", getTeacherInfo).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getteacherinfo","application/json",body)
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		//เก็บค่าตอบกลับไว้ใน course
		var teacherInfo teacher.Teacher
		json.Unmarshal(output, &teacherInfo)

		expected := teacher.Teacher{
			UserID: "",
			Firstname: "Myname",
			Surname: "Mysurname",
			Email: "teacher@kmitl.ac.th",
		}
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,teacherInfo)
	})

	sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "AAAA99")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "AAAA99")
	if err != nil {
		panic(err)
   }
}

func Test_editTeacherInfo(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "AAAA99", "testeditteacherinfo", "123456", "teacher")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "AAAA99", "Myname", "Mysurname", "teacher@kmitl.ac.th");
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Get Teacher Info", func(t *testing.T) {
		var jsonStr = `{"Firstname": "Editedname","Surname": "Editedsurname","Email": "Editedteacher@kmitl.ac.th","Username": "testeditteacherinfo"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/editteacherinfo", editTeacherInfo).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/editteacherinfo","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		//เก็บค่าตอบกลับไว้ใน course
		var teacherInfo teacher.Teacher
		json.Unmarshal(output, &teacherInfo)

		expected := teacher.Teacher{
			UserID: "",
			Firstname: "Editedname",
			Surname: "Editedsurname",
			Email: "Editedteacher@kmitl.ac.th",
		}
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,teacherInfo)
	})

	sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "AAAA99")
	if err != nil {
		panic(err)
   }

   sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "AAAA99")
	if err != nil {
		panic(err)
   }
}

func Test_getDescription(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission, announcement, description)VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public", "This is course announcement.", "This is course description.")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Get Description", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "123456"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getdescription", getDescription).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getdescription","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var description string

		json.Unmarshal(output, &description)

		expected := "This is course description."
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,description)
	})

   sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "123456")
	if err != nil {
		panic(err)
   }
}

func Test_editDescription(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission, announcement, description)VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public", "This is course announcement.", "This is course description.")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Edit Description", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "123456", "Description": "Edited Description." }`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/editdescription", editDescription).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/editdescription","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response string

		var description string

		json.Unmarshal(output, &response)

		sqlStatement := `SELECT description FROM course WHERE coursecode=$1;`
		rows, err := db.Query(sqlStatement, "123456")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&description)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"success",response)
		assert.Equal(t,"Edited Description.",description)
	})

    sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "123456")
	if err != nil {
		panic(err)
   }
}

func Test_getAnnouncement(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission, announcement, description)VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public", "This is course announcement.", "This is course description.")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Get Announcement", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "123456"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getannouncement", getAnnouncement).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getannouncement","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var announcement string

		json.Unmarshal(output, &announcement)

		expected := "This is course announcement."
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,announcement)
	})

   sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "123456")
	if err != nil {
		panic(err)
   }
}

func Test_editAnnouncement(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission, announcement, description)VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = db.Exec(sqlStatement, "123456", "55555555", "Test Test", "2563", "Public", "This is course announcement.", "This is course description.")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Edit Announcement", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "123456", "Announcement": "Edited course announcement."}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/editannouncement", editAnnouncement).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/editannouncement","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response string

		var announcement string

		sqlStatement := `SELECT announcement FROM course WHERE coursecode=$1;`
		rows, err := db.Query(sqlStatement, "123456")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&announcement)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		json.Unmarshal(output, &response)
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"success",response)
		assert.Equal(t,"Edited course announcement.",announcement)
	})

   sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "123456")
	if err != nil {
		panic(err)
   }
}

func Test_addStudentToCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "999XXX", "addstudenttocourseintegrationtest", "123456", "student")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO student (userid, studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, "999XXX", "88010135", "integrationtest", "addstudent", "student@kmitl.ac.th")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Add student to course", func(t *testing.T) {
		var jsonStr = `{"StudentID": "88010135", "CourseCode": "WWWWWW"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/addstudent", addStudentToCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/addstudent","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response student.Student

		var status string

		sqlStatement := `SELECT status FROM coursemember WHERE coursecode=$1 and userid=$2;`
		rows, err := db.Query(sqlStatement, "WWWWWW", "999XXX")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&status)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		json.Unmarshal(output, &response)
	
		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"88010135",response.StudentID)
		assert.Equal(t,"integrationtest",response.Firstname)
		assert.Equal(t,"addstudent",response.Surname)
		assert.Equal(t,"pending",status)
	})

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "999XXX")
	if err != nil {
		panic(err)
	}

   sqlStatement = `DELETE FROM student WHERE userid=$1;`
   _, err = db.Exec(sqlStatement, "999XXX")
   if err != nil {
	   panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "999XXX")
	if err != nil {
		panic(err)
	}
}

func Test_addTeacherToCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "999XXX", "addteachertocourseintegrationtest", "123456", "teacher")
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, "999XXX", "integrationtest", "addteacher", "teacher@kmitl.ac.th")
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Add teacher to course", func(t *testing.T) {
		var jsonStr = `{"Username": "addteachertocourseintegrationtest", "CourseCode": "WWWWWW"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/addteacher", addTeacherToCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/addteacher","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response teacher.Teacher

		json.Unmarshal(output, &response)

		var status string

		sqlStatement := `SELECT status FROM coursemember WHERE coursecode=$1 and userid=$2;`
		rows, err := db.Query(sqlStatement, "WWWWWW", "999XXX")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&status)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"integrationtest",response.Firstname)
		assert.Equal(t,"addteacher",response.Surname)
		assert.Equal(t,"pending",status)
	})

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "999XXX")
	if err != nil {
		panic(err)
	}

   sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
   _, err = db.Exec(sqlStatement, "999XXX")
   if err != nil {
	   panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "999XXX")
	if err != nil {
		panic(err)
	}
}

func Test_getStudentInCourse(t *testing.T) {

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

	t.Run("Integration Test: Get student in course", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "OOOOOO"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getstudentincourse", getStudentInCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getstudentincourse","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response []studentInCourse

		json.Unmarshal(output, &response)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,response)
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

func Test_getTeacherInCourse(t *testing.T) {

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

	t.Run("Integration Test: Get teacher in course", func(t *testing.T) {
		var jsonStr = `{"CourseCode": "OOOOOO"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getteacherincourse", getTeacherInCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getteacherincourse","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var response []teacherInCourse

		json.Unmarshal(output, &response)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,expected,response)
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

func Test_getUserRole(t *testing.T) {

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

	t.Run("Integration Test: Get student role", func(t *testing.T) {
		var jsonStr = `{"Username" : "getstudentrole"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.Handle("/getuserrole",getUserRole).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getuserrole","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"student",string(output))
	})

	t.Run("Integration Test: Get teacher role", func(t *testing.T) {
		var jsonStr = `{"Username" : "getteacherrole"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.Handle("/getuserrole", getUserRole).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getuserrole","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,"teacher",string(output))
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

func Test_deleteTeacherInCourse(t *testing.T) {

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

	t.Run("Integration Test: Delete Teacher", func(t *testing.T) {
		var jsonStr = `{"CourseCode" : "OOOOOO", "Username": "deleteteacher"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.Handle("/deleteteacherincourse",deleteTeacherInCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/deleteteacherincourse","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var teacherInfo teacher.Teacher
		json.Unmarshal(output, &teacherInfo)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,teacher1.Firstname,teacherInfo.Firstname)
		assert.Equal(t,teacher1.Surname,teacherInfo.Surname)
		assert.Equal(t,teacher1.Email,teacherInfo.Email)
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

func Test_deleteStudentInCourse(t *testing.T) {

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

	t.Run("Integration Test: Delete Student", func(t *testing.T) {
		var jsonStr = `{"CourseCode" : "OOOOOO", "Username": "deletestudent"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.Handle("/deletestudentincourse",deleteStudentInCourse).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/deletestudentincourse","application/json",body)
		if err != nil {
			t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)

		var studentInfo student.Student
		json.Unmarshal(output, &studentInfo)

		assert.Equal(t,http.StatusOK,resp.StatusCode)
		assert.Equal(t,student1.Firstname,studentInfo.Firstname)
		assert.Equal(t,student1.Surname,studentInfo.Surname)
		assert.Equal(t,student1.Email,studentInfo.Email)
	})

	sqlStatement = `DELETE FROM teacher WHERE userid=$1;`
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