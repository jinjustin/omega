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

	"encoding/json"
	"omega/course"
)

func Test_createCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `bNMv75`, `createcourseintegrationtest`,`123456` ,`teacher`)
	if err != nil {
		panic(err)
	}

	t.Run("Integration Test: Create Course", func(t *testing.T) {

		var jsonStr = `{"CourseName": "Test Test","CourseID": "55555555","Year": "2563","Permission": "Public","UserID": "createcourseintegrationtest"}`

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
		r.HandleFunc("/deletecourse", deleteCourse).Methods("POST")
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
		var jsonStr = `{"UserName": "getcourselistintegrationtest"}`

		body := strings.NewReader(jsonStr)
	
		r := mux.NewRouter().StrictSlash(true)
		r.HandleFunc("/getcourselist", getCourseList).Methods("POST")
		ts := httptest.NewServer(r)
		defer ts.Close()
	
		resp, err := http.Post(ts.URL + "/getcourselist","application/json",body)
		if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
		}
		defer resp.Body.Close()
		output, err := ioutil.ReadAll(resp.Body)
		
		fmt.Println(string(output))

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