package testcontroller

import (
	"fmt"
	"github.com/jinjustin/omega/test"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/sqs/goreturns/returns"
)

func createTest(courseID string, courseCode string, username string, status string, name string, duration string, start string, date string, description string) []byte {

	var t test.Test

	if checkTestName(courseID, name) == false {

		t = test.Test{
			TestID:      "",
			CourseID:    "",
			CourseCode:  "",
			UserID:      "",
			Status:      "",
			Name:        "Duplicate test name.",
			Duration:    "",
			Start:       "",
			Date:        "",
			Description: "",
		}

		return t.GetTestDetail()
	}

	testID := generateTestID()

	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if duration == "" {
		duration = "0"
	}

	if start == "" {
		start = "00:00:00"
	}

	if date == "" {
		date = "1-1-1970"
	}

	if description == "" {
		description = "this test don't have description"
	}

	t = test.Test{
		TestID:      testID,
		CourseID:    courseID,
		CourseCode:  courseCode,
		UserID:      userID,
		Status:      status,
		Name:        name,
		Duration:    duration,
		Start:       start,
		Date:        date,
		Description: description,
	}

	sqlStatement = `INSERT INTO test (testID, courseid, coursecode, userid, status, name, duration, start, date, description)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = db.Exec(sqlStatement, t.TestID, t.CourseID, t.CourseCode, t.UserID, t.Status, t.Name, t.Duration, t.Start, t.Date, t.Description)
	if err != nil {
		panic(err)
	}

	return t.GetTestDetail()
}

func getTestList(courseID string, username string) []test.Test {
	var testList []test.Test
	var t test.Test
	var role string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT role FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&role)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if (role == "student"){

		sqlStatement = `SELECT testid,coursecode,userid,status,name,duration,start,date,description FROM test WHERE courseid=$1 and status=$2;`
		rows, err = db.Query(sqlStatement, courseID,"publish")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t.TestID, &t.CourseCode, &t.UserID, &t.Status, &t.Name, &t.Duration, &t.Start, &t.Date, &t.Description)
			if err != nil {
				panic(err)
			}
			t.CourseID = courseID
			testList = append(testList, t)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		return testList
	}

	sqlStatement = `SELECT testid,coursecode,userid,status,name,duration,start,date,description FROM test WHERE courseid=$1;`
	rows, err = db.Query(sqlStatement, courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.TestID, &t.CourseCode, &t.UserID, &t.Status, &t.Name, &t.Duration, &t.Start, &t.Date, &t.Description)
		if err != nil {
			panic(err)
		}
		t.CourseID = courseID
		testList = append(testList, t)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return testList
}

func getTestInfo(testID string) []byte{
	var t test.Test

	t.TestID = testID

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT courseid,coursecode,userid,status,name,duration,start,date,description FROM test WHERE testid=$1;`
	rows, err := db.Query(sqlStatement, t.TestID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseID,&t.CourseCode, &t.UserID, &t.Status, &t.Name, &t.Duration, &t.Start, &t.Date, &t.Description)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return t.GetTestDetail()
}

func editTestInfo(testID string, name string, duration string, start string, date string, description string) string{

	if duration == "" {
		duration = "0"
	}

	if start == "" {
		start = "00:00:00"
	}

	if date == "" {
		date = "1-1-1970"
	}

	if description == "" {
		description = "this test don't have description"
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE test SET name=$1, duration=$2, start=$3, date=$4, description=$5 WHERE testid=$6`

	_, err = db.Exec(sqlStatement, name, duration, start, date, description, testID)
	if err != nil {
		panic(err)
	}

	return "success"
}

func deleteTest(testID string) []byte{

	t := test.Test{
		TestID:      "",
		CourseID:    "",
		CourseCode:  "",
		UserID:      "",
		Status:      "",
		Name:        "",
		Duration:    "",
		Start:       "",
		Date:        "",
		Description: "",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT courseid,coursecode,userid,status,name,duration,start,date,description FROM test WHERE testid=$1;`
	rows, err := db.Query(sqlStatement, testID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseID,&t.CourseCode,&t.UserID,&t.Status,&t.Name,&t.Duration,&t.Start,&t.Date,&t.Description)
		if err != nil {
			panic(err)
		}
		t.TestID = testID
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM test WHERE testID=$1;`
	_, err = db.Exec(sqlStatement, testID)
	if err != nil {
		panic(err)
	}

	return t.GetTestDetail()
}

func checkTestName(courseID string, name string) bool {
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var testID string

	sqlStatement := `SELECT testid FROM test WHERE courseid=$1 and name=$2;`
	row := db.QueryRow(sqlStatement, courseID, name)
	err = row.Scan(&testID)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
}

func generateTestID() string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

//API

//CreateTest is a API that use to create test in the course.
var CreateTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseID    string
		CourseCode  string
		Status      string
		Name        string
		Duration    string
		Start       string
		Date        string
		Description string
	}

	username := authentication.GetUsername(r)
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(createTest(input.CourseID, input.CourseCode, username, input.Status, input.Name, input.Duration, input.Start, input.Date, input.Description))
})

//GetTestList is a API that use to get all test in the course.
var GetTestList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseID string
		Username string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(getTestList(input.CourseID,input.Username))
})

//GetTestInfo is a API that use to get test information
var GetTestInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		TestID string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(getTestInfo(input.TestID))
})

//EditTestInfo is a API that use to edit test information
var EditTestInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		TestID string
		Name string
		Duration string
		Start string
		Date string
		Description string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(editTestInfo(input.TestID,input.Name,input.Duration,input.Start,input.Date,input.Description))
})

//DeleteTest is a API that use to delete test
var DeleteTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		TestID string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(deleteTest(input.TestID))
})
