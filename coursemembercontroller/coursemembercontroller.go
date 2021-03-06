package coursemembercontroller

import (
	"database/sql"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/student"
	"github.com/jinjustin/omega/teacher"
	"github.com/jinjustin/omega/mail"

	"github.com/jinjustin/omega/authentication"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/rand"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	
	"os"
	"io"

	"strconv"

	"errors"
)

func addStudentToCourse(studentID string, courseCode string) []byte {

	var userID string
	var firstName string
	var surName string
	var email string

	s := student.Student{
		UserID:    "Can't Invite this student",
		StudentID: "",
		Firstname: "",
		Surname:   "",
		Email:     "",
	}

	err := addStudentToSystem(studentID)
	if err != nil{
		panic(err)
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid,firstname,surname,email FROM student WHERE studentid=$1;`
	rows, err := db.Query(sqlStatement, studentID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userID, &firstName, &surName, &email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if checkMemberInCourse(userID, courseCode) {
		s = student.Student{
			UserID:    "",
			StudentID: studentID,
			Firstname: firstName,
			Surname:   surName,
			Email:     email,
		}

		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, courseCode, userID, "student", "pending")
		if err != nil {
			panic(err)
		}

		var courseName string
		var courseID string

		sqlStatement = `SELECT coursename,courseid FROM course WHERE coursecode=$1;`
		rows, err := db.Query(sqlStatement, courseCode)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&courseName,&courseID)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		message := "<br>คุณได้รับคำเชิญให้เข้าร่วมการสอบในวิชา " + courseID + " - " + courseName + " <br> <br>กดที่ลิงค์ต่อไปนี้เพื่อเข้าร่วมการสอบ http://142.93.177.152:3000/Accept?studentid=" + studentID + "&coursecode=" + courseCode + "<br>"
		mail.Send(email,"Course Invitation",message)
	}

	return s.GetStudentDetail()
}

func addMultipleStudentsToCourse(studentIDs []string, courseCode string) []student.Student {
	var students []student.Student

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, studentID := range studentIDs {

		err := addStudentToSystem(studentID)
		if err != nil{
			panic(err)
		}

		sqlStatement := `SELECT userid,firstname,surname,email FROM student WHERE studentid=$1;`
		rows, err := db.Query(sqlStatement, studentID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
	
		for rows.Next() {
			s := student.Student{
				UserID:    "Can't Invite this student",
				StudentID: studentID,
				Firstname: "",
				Surname:   "",
				Email:     "",
			}
			err = rows.Scan(&s.UserID, &s.Firstname, &s.Surname, &s.Email)
			if err != nil {
				panic(err)
			}
			students = append(students, s)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	for _, s := range students {
		if checkMemberInCourse(s.UserID, courseCode) && s.UserID != "Can't Invite this student" {
	
			sqlStatement := `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`
	
			_, err = db.Exec(sqlStatement, courseCode, s.UserID, "student", "pending")
			if err != nil {
				panic(err)
			}
	
			var courseName string
			var courseID string
	
			sqlStatement = `SELECT coursename,courseid FROM course WHERE coursecode=$1;`
			rows, err := db.Query(sqlStatement, courseCode)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
		
			for rows.Next() {
				err = rows.Scan(&courseName,&courseID)
				if err != nil {
					panic(err)
				}
			}
			err = rows.Err()
			if err != nil {
				panic(err)
			}
	
			message := "<br>คุณได้รับคำเชิญให้เข้าร่วมการสอบในวิชา " + courseID + " - " + courseName + " <br> <br>กดที่ลิงค์ต่อไปนี้เพื่อเข้าร่วมการสอบ http://142.93.177.152:3000/Accept?studentid=" + s.StudentID + "&coursecode=" + courseCode + "<br>"
			mail.Send(s.Email,"Course Invitation",message)
		}
	}

	return students
}

func addTeacherToCourse(username string, courseCode string) []byte {

	var userID string
	var firstName string
	var surName string

	t := teacher.Teacher{
		UserID:    "Can't Join this course",
		Firstname: "",
		Surname:   "",
		Email:     "",
	}

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

	sqlStatement = `SELECT firstname,surname FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&firstName, &surName)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if checkMemberInCourse(userID, courseCode) {
		t = teacher.Teacher{
			UserID:    "",
			Firstname: firstName,
			Surname:   surName,
			Email:     "",
		}

		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, courseCode, userID, "teacher", "join")
		if err != nil {
			panic(err)
		}
		return t.GetTeacherDetail()
	}

	return t.GetTeacherDetail()
}

func teacherAddCourse(username string, courseCode string) error {

	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	if checkMemberInCourse(userID, courseCode) {

		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, courseCode, userID, "teacher", "join")
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("omega: User already course member.")
}

func getStudentInCourse(courseCode string) []byte {
	var userIDs []string

	type studentInCourse struct {
		StudentID string
		Firstname string
		Surname   string
		Status    string
	}

	var studentInCourses []studentInCourse

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err := db.Query(sqlStatement, courseCode, "student", "join")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT studentid,firstname,surname FROM student WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID string
			var firstname string
			var surname string

			err = rows.Scan(&studentID, &firstname, &surname)
			if err != nil {
				panic(err)
			}

			s := studentInCourse{
				StudentID: studentID,
				Firstname: firstname,
				Surname:   surname,
				Status:    "join",
			}

			studentInCourses = append(studentInCourses, s)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	userIDs = nil

	sqlStatement = `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err = db.Query(sqlStatement, courseCode, "student", "pending")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT studentid,firstname,surname FROM student WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID string
			var firstname string
			var surname string

			err = rows.Scan(&studentID, &firstname, &surname)
			if err != nil {
				panic(err)
			}

			s := studentInCourse{
				StudentID: studentID,
				Firstname: firstname,
				Surname:   surname,
				Status:    "pending",
			}

			studentInCourses = append(studentInCourses, s)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	for num1, i := range studentInCourses{
		for num2, j := range studentInCourses{
			idI, _ := strconv.Atoi(i.StudentID)
			idJ, _ := strconv.Atoi(j.StudentID)
			if idI < idJ{
				studentInCourses[num1], studentInCourses[num2] = studentInCourses[num2], studentInCourses[num1]
			}
		}
	}

	b, err := json.Marshal(studentInCourses)
	if err != nil {
		panic(err)
	}

	return b
}

func getTeacherInCourse(courseCode string) []byte {
	var userIDs []string

	type teacherInCourse struct {
		Firstname string
		Surname   string
		Status    string
		Username string
	}

	var teacherInCourses []teacherInCourse

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err := db.Query(sqlStatement, courseCode, "teacher", "join")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		sqlStatement := `SELECT firstname,surname FROM teacher WHERE userid=$1;`
		teacherRows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer teacherRows.Close()

		for teacherRows.Next() {
			var firstname string
			var surname string

			err = teacherRows.Scan(&firstname, &surname)
			if err != nil {
				panic(err)
			}

			var username string

			sqlStatement := `SELECT username FROM users WHERE userid=$1;`
			usersRows, err := db.Query(sqlStatement, a)
			if err != nil {
				panic(err)
			}
			defer usersRows.Close()
			for usersRows.Next() {
				err = usersRows.Scan(&username)
				if err != nil {
					panic(err)
				}
			}
			err = usersRows.Err()
			if err != nil {
				panic(err)
			}

			t := teacherInCourse{
				Firstname: firstname,
				Surname:   surname,
				Status:    "join",
				Username: username,
			}

			teacherInCourses = append(teacherInCourses, t)
		}
		err = teacherRows.Err()
		if err != nil {
			panic(err)
		}
	}

	userIDs = nil

	sqlStatement = `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err = db.Query(sqlStatement, courseCode, "teacher", "pending")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {

		sqlStatement := `SELECT firstname,surname FROM teacher WHERE userid=$1;`
		teacherRows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer teacherRows.Close()

		for teacherRows.Next() {
			var firstname string
			var surname string

			err = teacherRows.Scan(&firstname, &surname)
			if err != nil {
				panic(err)
			}

			var username string

			sqlStatement := `SELECT username FROM users WHERE userid=$1;`
			usersRows, err := db.Query(sqlStatement, a)
			if err != nil {
				panic(err)
			}
			defer usersRows.Close()
			for usersRows.Next() {
				err = usersRows.Scan(&username)
				if err != nil {
					panic(err)
				}
			}
			err = usersRows.Err()
			if err != nil {
				panic(err)
			}

			t := teacherInCourse{
				Firstname: firstname,
				Surname:   surname,
				Status:    "pending",
				Username: username,
			}

			teacherInCourses = append(teacherInCourses, t)
		}

		err = teacherRows.Err()
		if err != nil {
			panic(err)
		}
	}

	b, err := json.Marshal(teacherInCourses)
	if err != nil {
		panic(err)
	}

	return b
}

func approveStudentJoinCourse(studentID string, courseCode string) string {

	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM student WHERE studentid=$1;`
	rows, err := db.Query(sqlStatement, studentID)
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

	sqlStatement = `UPDATE coursemember SET status=$1 WHERE coursecode=$2 and userid=$3;`

	_, err = db.Exec(sqlStatement, "join", courseCode, userID)
	if err != nil {
		panic(err)
	}

	return "success"
}

func approveTeacherJoinCourse(username string, courseCode string) string {

	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM student WHERE username=$1;`
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

	sqlStatement = `UPDATE coursemember SET status=$1 WHERE coursecode=$2 and userid=$3;`

	_, err = db.Exec(sqlStatement, "join", courseCode, userID)
	if err != nil {
		panic(err)
	}

	return "success"
}

//DeclineJoinCourse is ฟังก์ชันสำหรับให้ผู้สอนรองรับการเข้าร่วม Course ของผู้เรียน
func DeclineJoinCourse(userID string, courseCode string) string {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `DELETE from coursemember WHERE coursecode=$1 and userid=$2;`

	_, err = db.Exec(sqlStatement, courseCode, userID)
	if err != nil {
		panic(err)
	}

	return "success"
}

func getUserRole(username string) string {
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

	return role
}

func deleteTeacherInCourse(courseCode string, username string) []byte {
	t := teacher.Teacher{
		UserID:    "Can't find.",
		Firstname: "",
		Surname:   "",
		Email:     "",
	}

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

	var firstname string
	var surname string
	var email string

	sqlStatement = `SELECT firstname, surname, email FROM teacher WHERE userid=$1;`
	row := db.QueryRow(sqlStatement, userID)
	err = row.Scan(&firstname, &surname, &email)
	if err != nil {
		panic(err)
	}

	t = teacher.Teacher{
		Firstname: firstname,
		Surname:   surname,
		Email:     email,
	}

	sqlStatement = `DELETE FROM coursemember WHERE userID=$1 and coursecode=$2;`
	_, err = db.Exec(sqlStatement, userID, courseCode)
	if err != nil {
		panic(err)
	}

	return t.GetTeacherDetail()
}

func deleteStudentInCourse(courseCode string, username string) []byte {
	s := student.Student{
		UserID:    "Can't find.",
		StudentID: "",
		Firstname: "",
		Surname:   "",
		Email:     "",
	}

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

	var studentID string
	var firstname string
	var surname string
	var email string

	sqlStatement = `SELECT studentid, firstname, surname, email FROM student WHERE userid=$1;`
	row := db.QueryRow(sqlStatement, userID)
	err = row.Scan(&studentID, &firstname, &surname, &email)
	if err != nil {
		panic(err)
	}

	s = student.Student{
		StudentID: studentID,
		Firstname: firstname,
		Surname:   surname,
		Email:     email,
	}

	sqlStatement = `DELETE FROM coursemember WHERE userID=$1 and coursecode=$2;`
	_, err = db.Exec(sqlStatement, userID, courseCode)
	if err != nil {
		panic(err)
	}

	return s.GetStudentDetail()

}

func checkMemberInCourse(userID string, courseCode string) bool {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var status string

	sqlStatement := `SELECT status FROM coursemember WHERE userid=$1 and coursecode=$2;`
	row := db.QueryRow(sqlStatement, userID, courseCode)
	err = row.Scan(&status)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
}

func readStudentIDFromExcel(name string) []string{
	var studentIDs []string

	f, err := excelize.OpenFile(name)
    if err != nil {
        fmt.Println(err)
        return studentIDs
	}

    // Get all the rows in the Sheet1.
    rows, err := f.GetRows("ใบคะแนนสอบ")
    for _, row := range rows {
		if len(row) == 4 {
			if len(row[1]) == 8{
				studentIDs = append(studentIDs, row[1])
			}
		}
	}
	return studentIDs
}

func deleteFile(name string){
	err := os.Remove(name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addStudentToSystem(studentID string) error{

	userID := generateUserID()
	password := generatePassword()

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		 return err
	}
	defer db.Close()

	if checkUserExist(studentID){
		sqlStatement := `INSERT INTO users (userid , username, password, role)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, userID, studentID, password, "student")
		if err != nil {
			return err
		}
	
		sqlStatement = `INSERT INTO student (userid , studentid, firstname, surname, email)VALUES ($1, $2, $3, $4, $5)`
	
		_, err = db.Exec(sqlStatement, userID, studentID, "Annonymous", "Student", studentID + "@kmitl.ac.th")
		if err != nil {
			return err
		}

		message := "<br>คุณได้รับคำเชิญให้เข้าร่วมใช้งาน Omega: Online Exam Platform <br> <br>คุณสามารถใช้ Username และ Password นี้ในการ Login เข้าสู่ระบบ<br><br>Username: " + studentID + "<br><br>Password: " + password + " <br><br>คุณสามารถเข้าใช้งานระบบได้ที่ http://142.93.177.152:3000<br>"
		mail.Send(studentID + "@kmitl.ac.th","Course Invitation",message)
	}
	return nil
}

func addTeacherToSystem(username string, password string, email string) error{

	userID := generateUserID()

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	if checkUserExist(username){
		sqlStatement := `INSERT INTO users (userid , username, password, role)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, userID, username, password, "teacher")
		if err != nil {
			return err
		}
	
		sqlStatement = `INSERT INTO teacher (userid, firstname, surname, email)VALUES ($1, $2, $3, $4)`
	
		_, err = db.Exec(sqlStatement, userID, "Annonymous", "Teacher", email)
		if err != nil {
			return err
		}
	}

	return nil
}

func editStudentPassword(username string, oldPassword string, newPassword string) error{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `UPDATE users SET password = $1 WHERE username = $2 and password = $3`

	_, err = db.Exec(sqlStatement, newPassword, username, oldPassword)
	if err != nil {
		return err
	}

	return nil
}

func generateUserID() string{
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func generatePassword() string{
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func checkUserExist(username string) bool {

	var userid string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	row := db.QueryRow(sqlStatement, username)
	err = row.Scan(&userid)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
}

//API

//AddStudentToCourse is a API that use to add student to course.
var AddStudentToCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		StudentID  string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(addStudentToCourse(input.StudentID, input.CourseCode))
})

//AddTeacherToCourse is a API that use to add teacher to course.
var AddTeacherToCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Username   string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(addTeacherToCourse(input.Username, input.CourseCode))
})

//GetStudentInCourse is a API that use to get information of all student in course.
var GetStudentInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(getStudentInCourse(input.CourseCode))
})

//GetTeacherInCourse is a API that use to get information of all teacher in course.
var GetTeacherInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(getTeacherInCourse(input.CourseCode))
})

//DeleteTeacherInCourse is a function that use to delete teacher in course
var DeleteTeacherInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseCode string
		Username   string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(deleteTeacherInCourse(input.CourseCode, input.Username)))
})

//DeleteStudentInCourse is a function that use to delete student in course
var DeleteStudentInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseCode string
		Username   string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(deleteStudentInCourse(input.CourseCode, input.Username)))
})

//ApproveStudentJoinCourse is a function that teacher use to approve for student to join course.
var ApproveStudentJoinCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		StudentID string
		CourseCode string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(approveStudentJoinCourse(input.StudentID, input.CourseCode)))
})

//StudentAcceptJoinCourse is a function that student use to accept invitation to join course.
var StudentAcceptJoinCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	studentID := r.FormValue("studentid")
	courseCode := r.FormValue("coursecode")
	w.Write([]byte(approveStudentJoinCourse(studentID, courseCode)))
})

//AddStudentFromExcelFile is a function that teacher can add student to course by using excel file.
var AddStudentFromExcelFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseCode := r.Header.Get("coursecode")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer deleteFile(handler.Filename)
	defer file.Close()

	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(addMultipleStudentsToCourse(readStudentIDFromExcel(handler.Filename),courseCode))
})

//TestUploadExcelFile is a function that use to test upload file
var TestUploadExcelFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer deleteFile(handler.Filename)
	defer file.Close()

	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(readStudentIDFromExcel(handler.Filename))
})

//ChangeStudentPassword is a api that use to allow student to change password.
var ChangeStudentPassword = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
	oldPassword := r.Header.Get("OldPassword")
	newPassword := r.Header.Get("NewPassword")

	username := authentication.GetUsername(r)

	err := editStudentPassword(username, oldPassword, newPassword)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//TeacherAddCourse is a api that teacher use to add course.
var TeacherAddCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
	courseCode := r.Header.Get("CourseCode")
	username := authentication.GetUsername(r)

	err := teacherAddCourse(username, courseCode)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})


var AddTeacherToSystem = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	

	username := r.Header.Get("username")
	password := r.Header.Get("password")
	email := r.Header.Get("email")

	err := addTeacherToSystem(username, password, email)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})