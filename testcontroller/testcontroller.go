package testcontroller

import (
	"fmt"
	"github.com/jinjustin/omega/test"
	//"github.com/jinjustin/omega/course"
	"github.com/jinjustin/omega/coursecontroller"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"github.com/jinjustin/omega/database"
	//"github.com/jinjustin/omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	//"github.com/sqs/goreturns/returns"
)

func postDetailTest(testID string, courseID string, topic string ,description string , datestart string, duration string, timestart string) error{
	
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
			return err
		}
		defer db.Close()
	
		sqlStatement := `INSERT INTO test (testid, courseid, topic, description, datestart, duration, timestart, status)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(sqlStatement, t.TestID, t.CourseID, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart,"Unset")
		if err != nil {
			return err
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
			return err
		}
		defer db.Close()

		sqlStatement := `UPDATE test SET topic=$1, description=$2, datestart=$3, duration=$4, timestart=$5 WHERE testid=$6`

		_, err = db.Exec(sqlStatement, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart, t.TestID)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteTest(testID string) error{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `DELETE FROM test WHERE testID=$1;`
	_, err = db.Exec(sqlStatement, testID)
	if err != nil {
		return err
	}

	return nil
}

func getDetailTest(testID string, courseID string) ([]byte, error){

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
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT courseid, topic, description, datestart, duration, timestart, status FROM test WHERE testid=$1 and courseid=$2;`
	rows, err := db.Query(sqlStatement, testID,courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return t.GetTestDetail(), nil
}

func changeDraftStatus(testID string, status string) error{
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `UPDATE test SET status=$1 WHERE testid=$2`

	_, err = db.Exec(sqlStatement, status, testID)
	if err != nil {
		return err
	}

	return nil
}

//GenerateTestID is a function that use to generate testID.
func GenerateTestID() string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func studentGetTestList(studentID string) ([]test.StudentCourseList, error){
	var studentCourseList  test.StudentCourseList
	var studentCourseLists  []test.StudentCourseList

	var testList []test.Test

	var testData []test.Test

	var t test.Test

	var testdates []string

	courselist , err := coursecontroller.GetStudentCourseList(studentID)
	if err != nil{
		return studentCourseLists, err
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return studentCourseLists, err
	}
	defer db.Close()

	for _, c := range courselist {
		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart FROM student WHERE courseid=$1, status='publish';`
		rows, err := db.Query(sqlStatement, c.CourseID)
		if err != nil {
			return studentCourseLists, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t.TestID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart)
			if err != nil {
				return studentCourseLists, err
			}
			t.CourseID = c.CourseID
			t.Status = "publish"
			testList = append(testList, t)

			check := true

			for _, d := range testdates{
				if d == t.Datestart{
					check = false
				}
			}

			if check{
				testdates = append(testdates, t.Datestart)
			}

		}
		err = rows.Err()
		if err != nil {
			return studentCourseLists, err
		}
	}

	sortedDate, err := sortDate(testdates)
	if err != nil{
		return studentCourseLists, err
	}

	for _, d := range sortedDate{
		studentCourseList.Datestart = d
		for _, l := range testList{
			if d == l.Datestart{
				t = l
			}
			testData = append(testData, t)
		}

		sortedTestData, err := sortTime(testData)
		if err != nil{
			return studentCourseLists, err
		}

		studentCourseList.TestData = sortedTestData
		testData = nil
		studentCourseLists = append(studentCourseLists, studentCourseList)
	}

	return studentCourseLists, nil
}

//sortDate is a function that use to sort date
func sortDate(testdates []string) ([]string, error){
	//var date string

	for num1, i := range testdates{
		for num2, j := range testdates{
			check := false

			yearI , err := strconv.Atoi(i[6:10])
			if err != nil{
				return testdates, err
			}

			yearJ , err := strconv.Atoi(j[6:10])
			if err != nil{
				return testdates, err
			}

			monthI , err := strconv.Atoi(i[3:5])
			if err != nil{
				return testdates, err
			}

			monthJ , err := strconv.Atoi(j[3:5])
			if err != nil{
				return testdates, err
			}

			dayI , err := strconv.Atoi(i[0:2])
			if err != nil{
				return testdates, err
			}

			dayJ , err := strconv.Atoi(j[0:2])
			if err != nil{
				return testdates, err
			}

			if yearI < yearJ{
				check = true
			}else if monthI < monthJ && yearI == yearJ{
				check = true
			}else if dayI < dayJ && monthI == monthJ && yearI == yearJ{
				check = true
			}

			if check{
			testdates[num1], testdates[num2] = testdates[num2], testdates[num1]
			}
		}
	}

	return testdates, nil
}

//sortTime is a function that use to sort date
func sortTime(testdata []test.Test) ([]test.Test, error){
	//var date string

	for num1, i := range testdata{
		for num2, j := range testdata{
			check := false

			hourI , err := strconv.Atoi(i.Timestart[0:2])
			if err != nil{
				return testdata, err
			}

			hourJ , err := strconv.Atoi(j.Timestart[0:2])
			if err != nil{
				return testdata, err
			}

			minuteI , err := strconv.Atoi(i.Timestart[3:5])
			if err != nil{
				return testdata, err
			}

			minuteJ , err := strconv.Atoi(j.Timestart[3:5])
			if err != nil{
				return testdata, err
			}

			secondI , err := strconv.Atoi(i.Timestart[6:8])
			if err != nil{
				return testdata, err
			}

			secondJ , err := strconv.Atoi(j.Timestart[6:8])
			if err != nil{
				return testdata, err
			}

			if hourI < hourJ{
				check = true
			}else if minuteI < minuteJ && hourI == hourJ{
				check = true
			}else if secondI < secondJ && minuteI == minuteJ && hourI == hourJ{
				check = true
			}

			if check{
				testdata[num1], testdata[num2] = testdata[num2], testdata[num1]
			}
		}
	}

	return testdata, nil
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

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "Can't read body.", http.StatusBadRequest)
            return
	}
	var input Input
	json.Unmarshal(reqBody, &input)

	err = postDetailTest(testID, courseID, input.Topic, input.Description, input.Datestart, input.Duration, input.Timestart)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetDetailTest is a API that use to get detail of the test in database.
var GetDetailTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestId")

	test, err := getDetailTest(testID, courseID)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(test)
})

//DeleteTest is a API that use to delete test
var DeleteTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	testID := r.Header.Get("TestId")

	err := deleteTest(testID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//ChangeDraftStatus is a API that use to change draft status of the test
var ChangeDraftStatus = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	testID := r.Header.Get("TestId")
	status := r.Header.Get("Status")

	err := changeDraftStatus(testID,status)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//TestSortDate is a API that use to change draft status of the test
var TestSortDate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	date1 := "01-01-2565"
	date2 := "16-06-2564"
	date3 := "16-05-2564"
	date4 := "15-05-2564"

	var dates []string

	dates = append(dates, date1)
	dates = append(dates, date2)
	dates = append(dates, date3)
	dates = append(dates, date4)

	data, err := sortDate(dates)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//TestSortTime is a API that use to change draft status of the test
var TestSortTime = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var t test.Test

	time1 := "01:01:00"
	time2 := "13:30:00"
	time3 := "13:00:00"
	time4 := "09:30:00"

	var testdata []test.Test

	t.Timestart = time1
	testdata = append(testdata, t)
	t.Timestart = time2
	testdata = append(testdata, t)
	t.Timestart = time3
	testdata = append(testdata, t)
	t.Timestart = time4
	testdata = append(testdata, t)

	data, err := sortTime(testdata)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})