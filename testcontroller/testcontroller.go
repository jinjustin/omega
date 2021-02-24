package testcontroller

import (
	"fmt"
	"github.com/jinjustin/omega/test"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"github.com/jinjustin/omega/database"
	//"github.com/jinjustin/omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/sqs/goreturns/returns"
)

func postDetailTest(testID string, courseID string, topic string ,description string , datestart string, duration string, timestart string) []byte{
	
	var t test.Test

	if testID == ""{
		t = test.Test{
			TestID : generateTestID(),
			CourseID : courseID,
			Topic: topic,
			Description: description,
			Datestart: datestart,
			Duration: duration,
			Timestart: timestart,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `INSERT INTO test (testid, courseid, topic, description, datestart, duration, timestart, status)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(sqlStatement, t.TestID, t.CourseID, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart,"Unset")
		if err != nil {
			panic(err)
		}

	}else{
		t = test.Test{
			TestID : testID,
			CourseID : courseID,
			Topic: topic,
			Description: description,
			Datestart: datestart,
			Duration: duration,
			Timestart: timestart,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `UPDATE test SET topic=$1, description=$2, datestart=$3, duration=$4, timestart=$5 WHERE testid=$6`

		_, err = db.Exec(sqlStatement, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart, t.TestID)
		if err != nil {
			panic(err)
		}
	}

	return t.GetTestDetail()
}

func deleteTest(testID string) []byte{

	t := test.Test{
		TestID : testID,
		CourseID : "",
		Topic: "",
		Description: "",
		Datestart: "",
		Duration: "",
		Timestart: "",
		Status: "",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT courseid, topic, description, datestart, duration, timestart, status FROM test WHERE testid=$1;`
	rows, err := db.Query(sqlStatement, testID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
		if err != nil {
			panic(err)
		}
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

func getDetailTest(testID string, courseID string) []byte{
	t := test.Test{
		TestID : testID,
		CourseID : "",
		Topic: "",
		Description: "",
		Datestart: "",
		Duration: "",
		Timestart: "",
		Status: "",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT courseid, topic, description, datestart, duration, timestart, status FROM test WHERE testid=$1 and courseid=$2;`
	rows, err := db.Query(sqlStatement, testID,courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
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

func changeDraftStatus(testid string, status string) string{
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE test SET status=$1 WHERE testid=$2`

	_, err = db.Exec(sqlStatement, status, testID)
	if err != nil {
		panic(err)
	}

	return status
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

//PostDetailTest is a API that use to send create or update detail of the test to database.
var PostDetailTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Topic string
		Description string
		Datestart string
		Duration string
		Timestart string
	}

	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestId")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(postDetailTest(testID, courseID, input.Topic, input.Description, input.Datestart, input.Duration, input.Timestart))
})

//GetDetailTest is a API that use to get detail of the test in database.
var GetDetailTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestId")
	w.Write(getDetailTest(testID, courseID))
})

//DeleteTest is a API that use to delete test
var DeleteTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	testID := r.Header.Get("TestId")
	w.Write(deleteTest(testID))
})

//ChangeDraftStatus is a API that use to change draft status of the test
var ChangeDraftStatus = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	testID := r.Header.Get("TestId")
	status := r.Header.Get("Status")
	w.Write(changeDraftStatus(testID,status))
})
