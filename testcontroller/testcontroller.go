package testcontroller

import (
	"fmt"
	"time"

	"github.com/jinjustin/omega/test"

	//"github.com/jinjustin/omega/course"
	"github.com/jinjustin/omega/coursecontroller"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"

	"github.com/jinjustin/omega/authentication"
	"github.com/jinjustin/omega/database"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/iancoleman/orderedmap"
	//"github.com/sqs/goreturns/returns"
)

func postDetailTest(testID string, courseCode string, topic string, description string, datestart string, duration string, timestart string) error {

	var t test.Test

	checkExist := checkTestExist(testID)

	if checkExist == sql.ErrNoRows {
		t = test.Test{
			TestID:      testID,
			CourseCode:  courseCode,
			Topic:       topic,
			Description: description,
			Datestart:   datestart,
			Duration:    duration,
			Timestart:   timestart,
			Situation:   "wait",
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			return err
		}
		defer db.Close()

		sqlStatement := `INSERT INTO test (testid, coursecode, topic, description, datestart, duration, timestart, status, situation)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = db.Exec(sqlStatement, t.TestID, t.CourseCode, t.Topic, t.Description, t.Datestart, t.Duration, t.Timestart, "Unset", t.Situation)
		if err != nil {
			return err
		}

	} else if checkExist == nil {
		t = test.Test{
			TestID:      testID,
			CourseCode:  courseCode,
			Topic:       topic,
			Description: description,
			Datestart:   datestart,
			Duration:    duration,
			Timestart:   timestart,
			Situation:   "wait",
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
	} else {
		return checkExist
	}

	fmt.Println("Check Exist", checkExist)

	return nil
}

func deleteTest(testID string) error {

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

func getDetailTest(testID string, courseCode string) ([]byte, error) {

	t := test.Test{
		TestID:      testID,
		CourseCode:  "",
		Topic:       "",
		Description: "",
		Datestart:   "",
		Duration:    "",
		Timestart:   "",
		Status:      "",
		Situation:   "",
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

func getAllTestInCourse(courseCode string, role string) ([]test.Test, error) {

	t := test.Test{
		TestID:      "",
		CourseCode:  courseCode,
		Topic:       "",
		Description: "",
		Datestart:   "",
		Duration:    "",
		Timestart:   "",
		Status:      "",
	}

	var allTest []test.Test

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if role == "teacher" {

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

	} else if role == "student" {

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

func getAllFinishTestInCourse(courseCode string) ([]test.FinishTest, error) {

	var allFinishTest []test.FinishTest
	var finishTest test.FinishTest
	var dummy string

	membercount := 0

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT testid, topic FROM test WHERE coursecode=$1;`
	rows, err := db.Query(sqlStatement, courseCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&finishTest.TestID,&finishTest.Topic)
		if err != nil {
			return nil, err
		}
		allFinishTest = append(allFinishTest, finishTest)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	sqlStatement = `SELECT userid FROM coursemember WHERE coursecode=$1 and role='student';`
	coursememberRows, err := db.Query(sqlStatement, courseCode)
	if err != nil {
		return nil, err
	}
	defer coursememberRows.Close()
	for coursememberRows.Next() {
		err = coursememberRows.Scan(&dummy)
		if err != nil {
			return nil, err
		}
		membercount += 1
	}
	err = coursememberRows.Err()
	if err != nil {
		return nil, err
	}

	for num, f := range allFinishTest{
		participantCount := 0
		finshScoring := 0

		sqlStatement = `SELECT studentid FROM answer WHERE testid=$1;`
		answerRows, err := db.Query(sqlStatement, f.TestID)
		if err != nil {
			return nil, err
		}
		defer answerRows.Close()
		for answerRows.Next() {
			err = answerRows.Scan(&dummy)
			if err != nil {
				return nil, err
			}
			participantCount += 1
		}
		err = answerRows.Err()
		if err != nil {
			return nil, err
		}

		sqlStatement = `SELECT studentid FROM answer WHERE testid=$1 and completepercent='100.00';`
		answerRows2, err := db.Query(sqlStatement, f.TestID)
		if err != nil {
			return nil, err
		}
		defer answerRows2.Close()
		for answerRows2.Next() {
			err = answerRows2.Scan(&dummy)
			if err != nil {
				return nil, err
			}
			finshScoring += 1
		}
		err = answerRows2.Err()
		if err != nil {
			return nil, err
		}

		if membercount == 0{
			allFinishTest[num].Paticipant = "0:0"
			allFinishTest[num].Process = "0.00"
		}else if participantCount == 0 {
			member := strconv.Itoa(membercount)
			allFinishTest[num].Paticipant = "0:" + member
			allFinishTest[num].Process = "0.00"
		}else if finshScoring == 0{
			member := strconv.Itoa(membercount)
			participant := strconv.Itoa(participantCount)
			allFinishTest[num].Paticipant = participant + ":" + member
			allFinishTest[num].Process = "0.00"
		}else if finshScoring != 0{
			member := strconv.Itoa(membercount)
			participant := strconv.Itoa(participantCount)
			process := (float64(finshScoring)/float64(participantCount))*100
			processString := fmt.Sprintf("%.2f", process)
			allFinishTest[num].Paticipant = participant + ":" + member
			allFinishTest[num].Process = processString
		}
	}

	if allFinishTest == nil{
		allFinishTest = make([]test.FinishTest,0)
	}

	return allFinishTest, nil
}

func changeDraftStatus(testID string, status string) error {
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

func getDraftStatus(testID string) (string, error) {

	var status string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return "", err
	}
	defer db.Close()

	sqlStatement := `SELECT status FROM test WHERE testid=$1;`
	rows, err := db.Query(sqlStatement, testID)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return "", err
		}
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	return status, nil
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

/*func studentGetTestList(studentID string) ([]byte, error) {

	var testList []test.Test

	var testData []test.Test

	var t test.Test

	var testdates []string

	o := orderedmap.New()

	courselist, err := coursecontroller.GetStudentCourseList(studentID)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for _, c := range courselist {
		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart FROM test WHERE coursecode=$1 and status='false' and situation != 'finish';`
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
			t.Status = "false"
			testList = append(testList, t)

			check := true

			for _, d := range testdates {
				if d == t.Datestart {
					check = false
				}
			}

			if check {
				testdates = append(testdates, t.Datestart)
			}

		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	sortedDate, err := sortDate(testdates)
	if err != nil {
		return nil, err
	}

	fmt.Println(sortedDate)
	fmt.Println(testList)

	for _, d := range sortedDate {
		for _, l := range testList {
			if d == l.Datestart {
				t = l
				testData = append(testData, t)
			}
		}

		sortedTestData, err := sortTime(testData)
		if err != nil {
			return nil, err
		}

		o.Set(d, sortedTestData)
		testData = nil
	}

	b, err := o.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}*/

func studentGetTestList(studentID string) ([]byte, error) {

	var testList []test.ForStudent

	var testData []test.ForStudent

	var t test.ForStudent

	var testdates []string

	o := orderedmap.New()

	courselist, err := coursecontroller.GetStudentCourseList(studentID)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for _, c := range courselist {
		sqlStatement := `SELECT testid, topic, description, datestart, duration, timestart FROM test WHERE coursecode=$1 and status='false' and situation != 'finish';`
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
			t.Status = "false"
			t.CourseID = c.CourseID
			testList = append(testList, t)

			check := true

			for _, d := range testdates {
				if d == t.Datestart {
					check = false
				}
			}

			if check {
				testdates = append(testdates, t.Datestart)
			}

		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	sortedDate, err := sortDate(testdates)
	if err != nil {
		return nil, err
	}

	for _, d := range sortedDate {
		for _, l := range testList {
			if d == l.Datestart {
				t = l
				testData = append(testData, t)
			}
		}

		sortedTestData, err := sortTime(testData)
		if err != nil {
			return nil, err
		}

		o.Set(d, sortedTestData)
		testData = nil
	}

	b, err := o.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}

//sortDate is a function that use to sort date
func sortDate(testdates []string) ([]string, error) {
	//var date string

	for num1, i := range testdates {
		for num2, j := range testdates {
			check := false

			yearI, err := strconv.Atoi(i[0:4])
			if err != nil {
				return testdates, err
			}

			yearJ, err := strconv.Atoi(j[0:4])
			if err != nil {
				return testdates, err
			}

			monthI, err := strconv.Atoi(i[5:7])
			if err != nil {
				return testdates, err
			}

			monthJ, err := strconv.Atoi(j[5:7])
			if err != nil {
				return testdates, err
			}

			dayI, err := strconv.Atoi(i[8:10])
			if err != nil {
				return testdates, err
			}

			dayJ, err := strconv.Atoi(j[8:10])
			if err != nil {
				return testdates, err
			}

			if yearI < yearJ {
				check = true
			} else if monthI < monthJ && yearI == yearJ {
				check = true
			} else if dayI < dayJ && monthI == monthJ && yearI == yearJ {
				check = true
			}

			if check {
				testdates[num1], testdates[num2] = testdates[num2], testdates[num1]
			}
		}
	}

	return testdates, nil
}

//sortTime is a function that use to sort date
func sortTime(testdata []test.ForStudent) ([]test.ForStudent, error) {
	//var date string

	for num1, i := range testdata {
		for num2, j := range testdata {
			check := false

			hourI, err := strconv.Atoi(i.Timestart[0:2])
			if err != nil {
				return testdata, err
			}

			hourJ, err := strconv.Atoi(j.Timestart[0:2])
			if err != nil {
				return testdata, err
			}

			minuteI, err := strconv.Atoi(i.Timestart[3:5])
			if err != nil {
				return testdata, err
			}

			minuteJ, err := strconv.Atoi(j.Timestart[3:5])
			if err != nil {
				return testdata, err
			}

			if hourI < hourJ {
				check = true
			} else if minuteI < minuteJ && hourI == hourJ {
				check = true
			}

			if check {
				testdata[num1], testdata[num2] = testdata[num2], testdata[num1]
			}
		}
	}

	return testdata, nil
}

func UpdateTestSituation() error{

	var t test.Test

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT testid, datestart, duration, timestart, situation FROM test WHERE situation != 'finish';`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.TestID, &t.Datestart, &t.Duration, &t.Timestart, &t.Situation)
		if err != nil {
			return err
		}

		currentTime := time.Now()
		year := strconv.Itoa(time.Now().Year())

		var month string
		var day string

		if time.Now().Day() < 10 {
			day = "0" + strconv.Itoa(time.Now().Day())
		} else {
			day = strconv.Itoa(time.Now().Day())
		}

		if int(time.Now().Month()) < 10 {
			month = "0" + strconv.Itoa(int(time.Now().Month()))
		} else {
			month = strconv.Itoa(int(time.Now().Month()))
		}

		dateNow := year + "-" + month + "-" + day
		timeStampString := currentTime.Format("2006-01-02 15:04:05")
		layOut := "2006-01-02 15:04:05"
		timeStamp, err := time.Parse(layOut, timeStampString)
		if err != nil {
			fmt.Println(err)
		}
		hr, min, _ := timeStamp.Clock()

		duration, _ := strconv.Atoi(t.Duration)

		durationDay :=0
		durationMonth := 0
		durationYear := 0
	
		for duration > 24 {
			duration -= 24
			durationDay += 1
			if durationDay > 30 {
				durationDay = 0
				durationMonth += 1
				if durationMonth > 12 {
					durationMonth = 0
					durationYear += 1
				}
			}
		}
	
		startHour, _ := strconv.Atoi(t.Timestart[0:2])

		startDay, _ := strconv.Atoi(t.Datestart[8:10])

		startMonth, _ := strconv.Atoi(t.Datestart[5:7])

		startYear, _ := strconv.Atoi(t.Datestart[0:4])

		finishHour := startHour + duration
		if finishHour > 24{
			finishHour = finishHour%24
			durationDay += 1
			if durationDay > 30 {
				durationDay = 0
				durationMonth += 1
				if durationMonth > 12 {
					durationMonth = 0
					durationYear += 1
				}
			}
		}

		dayFinish := startDay + durationDay
		if dayFinish > 30{
			dayFinish = dayFinish%30
			durationMonth += 1
			if durationMonth > 12 {
				durationMonth = 0
				durationYear += 1
			}
		}

		monthFinish := startMonth + durationMonth
		if monthFinish >= 12 {
			monthFinish = monthFinish%12
			durationYear += 1
		}

		yearFinishString := strconv.Itoa(startYear + durationYear)

		var monthFinishString string
		var dayFinishString string


		if dayFinish < 10 {
			dayFinishString = "0" + strconv.Itoa(dayFinish)
		} else {
			dayFinishString = strconv.Itoa(dayFinish)
		}

		if monthFinish < 10 {
			monthFinishString = "0" + strconv.Itoa(monthFinish)
		} else {
			monthFinishString = strconv.Itoa(monthFinish)
		}

		dateFinish := yearFinishString + "-" + monthFinishString + "-" + dayFinishString

		var finishHourString string

		var strHr string
		var strMin string

		if hr < 10 {
			strHr = "0" + strconv.Itoa(hr)
		} else {
			strHr = strconv.Itoa(hr)
		}
		if min < 10 {
			strMin = "0" + strconv.Itoa(min)
		} else {
			strMin = strconv.Itoa(min)
		}

		if finishHour < 10 {
			finishHourString = "0" + strconv.Itoa(finishHour)
		}else{
			finishHourString = strconv.Itoa(finishHour)
		}

		timeNow := strHr + ":" + strMin

		timeFinish := finishHourString + ":" + t.Timestart[3:5]

		if dateNow == t.Datestart && timeNow == t.Timestart  && t.Situation == "wait"{
			sqlStatement := `UPDATE test SET situation=$1 WHERE testid=$2`
			_, err = db.Exec(sqlStatement, "start", t.TestID)
			if err != nil {
				return err
			}
		}else if dateNow == dateFinish && timeNow == timeFinish  && t.Situation == "start"{
			sqlStatement := `UPDATE test SET situation=$1 WHERE testid=$2`
			_, err = db.Exec(sqlStatement, "finish", t.TestID)
			if err != nil {
				return err
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

//API

//PostDetailTest is a API that use to send create or update detail of the test to database.
var PostDetailTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Topic       string
		Description string
		Datestart   string
		Duration    string
		Timestart   string
	}

	courseCode := r.Header.Get("CourseCode")
	testID := r.Header.Get("TestId")

	fmt.Println("Updatetestdatail", testID)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body.", http.StatusBadRequest)
		return
	}
	var input Input
	json.Unmarshal(reqBody, &input)

	err = postDetailTest(testID, courseCode, input.Topic, input.Description, input.Datestart, input.Duration, input.Timestart)
	if err != nil {
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

	if err != nil {
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
	if err != nil {
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

	err := changeDraftStatus(testID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetDraftStatus is a API that use to get draft status of the test
var GetDraftStatus = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	testID := r.Header.Get("TestId")

	status, err := getDraftStatus(testID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
})

//GetAllTestInCourse is a API that use to get information of all the tests in course.
var GetAllTestInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseCode := r.Header.Get("CourseCode")
	role := authentication.GetUserRole(r)

	allTest, err := getAllTestInCourse(courseCode, role)
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//StudentGetTestListByDay is a API that student use to get information of all the tests in course.
var StudentGetTestListByDay = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	username := authentication.GetUsername(r)

	testlist, err := studentGetTestList(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(testlist))
})


var GetAllFinishTestInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseCode := r.Header.Get("CourseCode")

	finishTests, err := getAllFinishTestInCourse(courseCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finishTests)
})
