package coursecontroller

import (
	"fmt"
	"omega/course"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"omega/database"
)

//CreateCourse is function that use to create classroom
func CreateCourse(courseName string, courseID string, year string, permission string,announcement string, description string,  username string) []byte {

	courseCode := generateCourseCode()

	c := course.Course{
		CourseCode: "",
		CourseID:   "Course ID Error",
		CourseName: "",
		Year:       "",
		Permission: "",
		Announcement: "",
		Description: "",
	}

	if checkInputValue(courseID, year) == true {

		var userID string

		c = course.Course{
			CourseCode: courseCode,
			CourseID:   courseID,
			CourseName: courseName,
			Year:       year,
			Permission: permission,
			Announcement: announcement,
			Description: description,
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

		sqlStatement = `INSERT INTO course (coursecode,courseid,coursename, year, permission,announcement,description)VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = db.Exec(sqlStatement, c.CourseCode, c.CourseID, c.CourseName, c.Year, c.Permission,c.Announcement,c.Description)
		if err != nil {
			panic(err)
		}

		sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, c.CourseCode, userID, `teacher`, `join`)
		if err != nil {
			panic(err)
		}

		return c.GetCourseDetail()
	}

	return c.GetCourseDetail()
}

//DeleteCourse is function that use to delete course by use CourseCode and userID
func DeleteCourse(courseCode string, username string) []byte {
	c := course.Course{
		CourseCode: "Can't find.",
		CourseID:   "",
		CourseName: "",
		Year:       "",
		Permission: "",
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

	if checkCourse(courseCode) && checkUser(courseCode, userID) {
		var courseName string
		var courseID string
		var year string
		var permission string
		var announcement string
		var description string

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT courseName,courseID,year,permission,announcement,description FROM course WHERE coursecode=$1;`
		row := db.QueryRow(sqlStatement, courseCode)
		err = row.Scan(&courseName, &courseID, &year, &permission, &announcement, &description)
		if err != nil {
			panic(err)
		}

		c := course.Course{
			CourseCode: courseCode,
			CourseID:   courseID,
			CourseName: courseName,
			Year:       year,
			Permission: permission,
			Announcement: announcement,
			Description: description,
		}

		sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
		_, err = db.Exec(sqlStatement, courseCode)
		if err != nil {
			panic(err)
		}

		sqlStatement = `DELETE FROM coursemember WHERE coursecode=$1 and userid=$2;`
		_, err = db.Exec(sqlStatement, courseCode, userID)
		if err != nil {
			panic(err)
		}
		return c.GetCourseDetail()
	}
	return c.GetCourseDetail()
}

//GetCourseList is use to get all classrooms that user is being member.
func GetCourseList(username string) []course.Course {
	var courseCodes []string
	var courses []course.Course
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

	sqlStatement = `SELECT coursecode FROM coursemember WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var courseCode string
		err = rows.Scan(&courseCode)
		if err != nil {
			panic(err)
		}
		courseCodes = append(courseCodes, courseCode)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range courseCodes {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT courseid,coursename,year,permission FROM course WHERE coursecode=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var courseID string
			var courseName string
			var year string
			var permission string

			err = rows.Scan(&courseID, &courseName, &year, &permission)
			if err != nil {
				panic(err)
			}

			c := course.Course{
				CourseCode: a,
				CourseID:   courseID,
				CourseName: courseName,
				Year:       year,
				Permission: permission,
			}

			courses = append(courses, c)
		}

		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	return courses
}

//GetDescription is a function that use to get description of the course
func GetDescription(courseCode string) string{

	var description string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT description FROM course WHERE coursecode=$1;`
	rows, err := db.Query(sqlStatement, courseCode)
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

	return description
}

//EditDescription is a function that use to edit course description.
func EditDescription(courseCode string, description string) string{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE course SET description=$1 WHERE coursecode=$2;`

	_, err = db.Exec(sqlStatement, description, courseCode)
	if err != nil {
	panic(err)
	}

	return "success"
}

//GetAnnouncement is a function that use to get course info
func GetAnnouncement(courseID string) string{
	var announcement string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT announcement FROM course WHERE coursecode=$1;`
	rows, err := db.Query(sqlStatement, courseID)
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

	return announcement
}

//EditAnnouncement is a function that use to edit course description.
func EditAnnouncement(courseCode string, announcement string) string{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE course SET announcement=$1 WHERE coursecode=$2;`

	_, err = db.Exec(sqlStatement, announcement, courseCode)
	if err != nil {
	panic(err)
	}

	return "success"
}

func checkInputValue(courseID string, year string) bool {
	var courseCode string

	if len(courseID) != 8 {
		return false
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT coursecode FROM course WHERE courseid=$1 and year=$2;`
	row := db.QueryRow(sqlStatement, courseID, year)
	err = row.Scan(&courseCode)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
}

func generateCourseCode() string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func checkCourse(courseCode string) bool {
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var courseName string

	sqlStatement := `SELECT coursename FROM course WHERE coursecode=$1;`
	row := db.QueryRow(sqlStatement, courseCode)
	err = row.Scan(&courseName)
	switch err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}

func checkUser(courseCode string, userID string) bool {
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var check string

	sqlStatement := `SELECT coursecode FROM coursemember WHERE userid=$1 and coursecode=$2;`
	row := db.QueryRow(sqlStatement, userID, courseCode)
	err = row.Scan(&check)
	switch err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}
