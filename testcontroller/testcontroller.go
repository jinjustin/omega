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
	"github.com/jinjustin/omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/iancoleman/orderedmap"
	//"github.com/sqs/goreturns/returns"
)

func postDetailTest(testID string, courseCode string, topic string ,description string , datestart string, duration string, timestart string) error{
	
	var t test.Test

	checkExist := checkTestExist(testID)

	if checkExist == sql.ErrNoRows{
		t = test.Test{
			TestID : testID,
			CourseCode : courseCode,
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
	
		sqlStatement := `INSERT INTO test (testid, coursecode, topic, description, datestart, duration, timestart, status)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(sqlStatement, t.TestID, t.CourseCode, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart,"Unset")
		if err != nil {
			return err
		}

	}else if checkExist == nil{
		t = test.Test{
			TestID : testID,
			CourseCode : courseCode,
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
	}else{
		return checkExist
	}

	fmt.Println("Check Exist", checkExist)

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

func getDetailTest(testID string, courseCode string) ([]byte, error){

	t := test.Test{
		TestID : testID,
		CourseCode : "",
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

	sqlStatement := `SELECT courseCode, topic, description, datestart, duration, timestart, status FROM test WHERE testid=$1 and coursecode=$2;`
	rows, err := db.Query(sqlStatement, testID, courseCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.CourseCode, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
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

func getAllTestInCourse(courseCode string, role string) ([]test.Test, error){

	t := test.Test{
		TestID : "",
		CourseCode : courseCode,
		Topic: "",
		Description: "",
		Datestart: "",
		Duration: "",
		Timestart: "",
		Status: "",
	}

	var allTest []test.Test

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if role=="teacher"{

		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart, status FROM test WHERE coursecode=$1;`
		rows, err := db.Query(sqlStatement, courseCode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t.TestID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
			if err != nil {
				return nil, err
			}
	
			allTest = append(allTest, t)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}

	}else if role=="student"{

		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart, status FROM test WHERE coursecode=$1 and status='false';`
		rows, err := db.Query(sqlStatement, courseCode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t.TestID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart, &t.Status)
			if err != nil {
				return nil, err
			}
	
			allTest = append(allTest, t)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}

	}

	return allTest, nil
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

func checkTestExist(testID string) error {

	var status string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT status FROM test WHERE testid=$1;`
	row := db.QueryRow(sqlStatement, testID)
	err = row.Scan(&status)
	return err
}

func studentGetTestList(studentID string) ([]byte, error){

	var testList []test.Test

	var testData []test.Test

	var t test.Test

	var testdates []string

	o := orderedmap.New()

	courselist , err := coursecontroller.GetStudentCourseList(studentID)
	if err != nil{
		return nil, err
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for _, c := range courselist {
		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart FROM student WHERE coursecode=$1, status='true';`
		rows, err := db.Query(sqlStatement, c.CourseCode)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t.TestID, &t.Topic, &t.Description, &t.Datestart, &t.Duration, &t.Timestart)
			if err != nil {
				return nil, err
			}
			t.CourseCode = c.CourseCode
			t.Status = "true"
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
			return nil, err
		}
	}

	sortedDate, err := sortDate(testdates)
	if err != nil{
		return nil, err
	}

	for _, d := range sortedDate{
		for _, l := range testList{
			if d == l.Datestart{
				t = l
			}
			testData = append(testData, t)
		}

		sortedTestData, err := sortTime(testData)
		if err != nil{
			return nil, err
		}

		o.Set(d,sortedTestData)
		testData = nil
	}

	b,err := o.MarshalJSON()
	if err != nil{
		return nil, err
	}

	return b, nil
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

	courseCode := r.Header.Get("CourseCode")
	testID := r.Header.Get("TestId")

	fmt.Println("Updatetestdatail", testID)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "Can't read body.", http.StatusBadRequest)
            return
	}
	var input Input
	json.Unmarshal(reqBody, &input)

	err = postDetailTest(testID, courseCode, input.Topic, input.Description, input.Datestart, input.Duration, input.Timestart)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetDetailTest is a API that use to get detail of the test in database.
var GetDetailTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	courseCode := r.Header.Get("CourseCode")
	testID := r.Header.Get("TestId")

	test, err := getDetailTest(testID, courseCode)

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

//GetAllTestInCourse is a API that use to get information of all the tests in course.
var GetAllTestInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseCode := r.Header.Get("CourseCode")
	role := authentication.GetUserRole(r)

	allTest, err := getAllTestInCourse(courseCode, role)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allTest)
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

//StudentGetTestListByDay is a API that student use to get information of all the tests in course.
var StudentGetTestListByDay = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	studentID := r.Header.Get("StudentID")

	testlist, err := studentGetTestList(studentID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testlist)
})