package testcontroller

import (
	"fmt"
	"omega/test"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"omega/database"
)

//CreateTest is function that use to create test in course.
func CreateTest(courseID string, courseCode string, username string, status string, name string, duration string, start string, date string, description string) []byte {

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

	t := test.Test{
		TestID: testID,
		CourseID: courseID,
		CourseCode: courseCode,
		UserID: userID,
		Status: status,
		Name: name,
		Duration: duration,
		Start: start,
		Date: date,
		Description: description,
	}

	sqlStatement = `INSERT INTO test (testID, courseid, coursecode, userid, status, name, duration, start, date, description)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = db.Exec(sqlStatement, t.TestID, t.CourseID, t.CourseCode, t.UserID, t.Status, t.Name, t.Duration, t.Start, t.Date, t.Description)
	if err != nil {
		panic(err)
	}

	return t.GetTestDetail()
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

/*
func checkTestTime(courseCode string, duration string, startTime string, startDate string) bool{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var tests []test.Test

	var test test.Test

	sqlStatement := `SELECT duration, startTime, startDate FROM test WHERE coursecode=$1;`
	rows, err := db.Query(sqlStatement, courseCode)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&test.Duration,&test.StartTime,&test.StartDate)
		if err != nil {
			panic(err)
		}

		tests = append(tests, test)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	testStartInt ,err := strconv.ParseInt(startTime[:2],10,32)

	testDurationInt ,err := strconv.ParseInt(duration,10,32)

	testEnd := testStartInt + testDurationInt

	for _, a := range tests {
		if(startDate == a.StartDate){

			aStartInt , err := strconv.ParseInt(a.StartTime[:2],10,32)
			if(err != nil){
				panic(err)
			}

			aDurationInt , err := strconv.ParseInt(a.Duration,10,32)
		
			aEnd := aStartInt + aDurationInt

			if((testStartInt >= aStartInt && testStartInt <= aEnd) || (testEnd >= aStartInt && testEnd <= aEnd)){
				return false
			}
		}
	}
	return true
}*/