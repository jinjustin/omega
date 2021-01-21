package coursecontroller

import (
	"database/sql"
	"encoding/json"
	"omega/course"
	"omega/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createCourse(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `a2yQbN`, `testcreatecourse`, `123456`, `teacher`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Correct Information", func(t *testing.T) {
		//Input
		courseID := "99999999"
		courseName := "Computer Programing"
		year := "2563"
		permission := "Private"
		username := "testcreatecourse"
		announcement := "Coure announcement here."
		description := "Course description here."
		//Expected Output
		expected := course.Course{
			CourseCode: "000001",
			CourseID:   "99999999",
			CourseName: "Computer Programing",
			Year:       "2563",
			Permission: "Private",
			Announcement: "Coure announcement here.",
			Description: "Course description here.",
		}
		//Output
		var output course.Course
		json.Unmarshal(createCourse(courseName, courseID, year, permission,announcement,description, username), &output)
		//Compare output to expected output
		assert.Equal(t, expected.CourseName, output.CourseName)
		assert.Equal(t, expected.CourseID, output.CourseID)
		assert.Equal(t, expected.Year, output.Year)
		assert.Equal(t, expected.Permission, output.Permission)
		assert.Equal(t, expected.Announcement, output.Announcement)
		assert.Equal(t, expected.Description, output.Description)
	})

	t.Run("Unit Test 002: Wrong Course ID", func(t *testing.T) {

		courseName := "Computer Programing"
		courseID := "9999999"
		year := "2563"
		permission := "Private"
		username := "testcreatecourse"

		var output course.Course
		json.Unmarshal(createCourse(courseName, courseID, year, permission,"","", username), &output)

		assert.Equal(t, "", output.CourseName)
		assert.Equal(t, "Course ID Error", output.CourseID)
		assert.Equal(t, "", output.Year)
		assert.Equal(t, "", output.Permission)
	})

	t.Run("Unit Test 003: Already Existed Course in same year", func(t *testing.T) {

		courseID := "99999999"
		courseName := "Computer Programing"
		year := "2563"
		permission := "Private"
		username := "testcreatecourse"
		var output course.Course
		json.Unmarshal(createCourse(courseName, courseID, year, permission,"","", username), &output)

		assert.Equal(t, "", output.CourseName)
		assert.Equal(t, "Course ID Error", output.CourseID)
		assert.Equal(t, "", output.Year)
		assert.Equal(t, "", output.Permission)
	})

	t.Run("Unit Test 004: Already Existed Course in different year", func(t *testing.T) {

		//input
		courseName := "Computer Programing"
		courseID := "99999999"
		year := "2564"
		permission := "Private"
		username := "testcreatecourse"
		//Expected output
		expected := course.Course{
			CourseCode: "000001",
			CourseID:   "99999999",
			CourseName: "Computer Programing",
			Year:       "2564",
			Permission: "Private"}

		var output course.Course
		json.Unmarshal(createCourse(courseName, courseID, year, permission,"","",username), &output)

		assert.Equal(t, expected.CourseName, output.CourseName)
		assert.Equal(t, expected.CourseID, output.CourseID)
		assert.Equal(t, expected.Year, output.Year)
		assert.Equal(t, expected.Permission, output.Permission)
	})

	db, err = sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement = `DELETE FROM course WHERE courseid=$1;`
	_, err = db.Exec(sqlStatement, "99999999")
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "a2yQbN")
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "a2yQbN")
	if err != nil {
		panic(err)
	}
}

func Test_deleteCourse(t *testing.T) {

	c := course.Course{
		CourseCode: "vB2p9U",
		CourseID:   "99999999",
		CourseName: "Test",
		Year:       "2563",
		Permission: "Public",
	}

	username := "testdeletecourse"

	userid := "ZZZZZZ"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, c.CourseCode, c.CourseID, c.CourseName, c.Year, c.Permission)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, c.CourseCode, userid, `teacher`, `join`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid,username,password,role) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, userid, username, `123456`, `teacher`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 005: Delete Non-Existed Course", func(t *testing.T) {

		//Input
		courseCode := "OBLMv5"

		var output course.Course
		json.Unmarshal(deleteCourse(courseCode, username), &output)

		assert.Equal(t, "Can't find.", output.CourseCode)
	})

	t.Run("Unit Test 006: Delete Existed Course but user don't in that course", func(t *testing.T) {

		//Input
		courseCode := c.CourseCode

		var output course.Course
		json.Unmarshal(deleteCourse(courseCode, "deletefail"), &output)

		assert.Equal(t, "Can't find.", output.CourseCode)
	})

	t.Run("Unit Test 007: Delete Existed Course", func(t *testing.T) {

		var output course.Course
		json.Unmarshal(deleteCourse(c.CourseCode, username), &output)

		assert.Equal(t, c, output)
	})

	db, err = sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid)
	if err != nil {
		panic(err)
	}
}

func Test_getCourselist(t *testing.T) {

	c1 := course.Course{
		CourseCode: "vB2p9U",
		CourseID:   "99999999",
		CourseName: "Test",
		Year:       "2563",
		Permission: "Public",
	}

	c2 := course.Course{
		CourseCode: "OBLMv5",
		CourseID:   "12309854",
		CourseName: "Test2",
		Year:       "2564",
		Permission: "Private",
	}

	username := "getcourselisttest"

	userid := "test01"

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, userid, username, `123456`, `teacher`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c1.CourseCode, c1.CourseID, c1.CourseName, c1.Year, c1.Permission)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, c1.CourseCode, userid, `teacher`, `join`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c2.CourseCode, c2.CourseID, c2.CourseName, c2.Year, c2.Permission)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, c2.CourseCode, userid, `teacher`, `join`)
	if err != nil {
		panic(err)
	}

	t.Run("get Course List from existed UserID", func(t *testing.T) {

		var expected []course.Course
		expected = append(expected, c1)
		expected = append(expected, c2)
		var output []course.Course
		output = getCourseList(username)

		assert.Equal(t, expected, output)
	})

	t.Run("get Course List from Abcdef UserID", func(t *testing.T) {

		var expected []course.Course
		var output []course.Course
		output = getCourseList("Abcdef")

		assert.Equal(t, expected, output)
	})

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, c1.CourseCode)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE coursecode=$1 and userid=$2;`
	_, err = db.Exec(sqlStatement, c1.CourseCode, userid)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, c2.CourseCode)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE coursecode=$1 and userid=$2;`
	_, err = db.Exec(sqlStatement, c2.CourseCode, userid)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, userid)
	if err != nil {
		panic(err)
	}
}

func Test_CourseDescriptionFunction(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename,year,permission,announcement,description)VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(sqlStatement, `ZZZZZZ`,`00000000`,`coursedescriptioncourse`,`2563`,`Private`,`Course announcement here.`,`Course description here.`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Get Course Description", func(t *testing.T) {
		//Input
		description := "Course description here."
		courseCode := "ZZZZZZ"

		//Output
		output := getDescription(courseCode)

		//Compare output to expected output
		assert.Equal(t, description,output)
	})

	t.Run("Unit Test 002: Edit Course Description", func(t *testing.T) {
		//Input
		description := "Course description being Edited"
		courseCode := "ZZZZZZ"

		//Output
		output := editDescription(courseCode,description)

		//Compare output to expected output
		assert.Equal(t, output, "success")
		assert.Equal(t, description,getDescription(courseCode))
	})

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "ZZZZZZ")
	if err != nil {
		panic(err)
	}
}

func Test_CourseAnnouncementFunction(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename,year,permission,announcement,description)VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(sqlStatement, `ZZZZZZ`,`00000000`,`coursedescriptioncourse`,`2563`,`Private`,`Course announcement here.`,`Course description here.`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Get Course Announcement", func(t *testing.T) {
		//Input
		announcement := "Course announcement here."
		courseCode := "ZZZZZZ"

		//Output
		output := getAnnouncement(courseCode)

		//Compare output to expected output
		assert.Equal(t, announcement, output)
	})

	t.Run("Unit Test 002: Edit Course Announcement", func(t *testing.T) {
		//Input
		description := "Course announcement being Edited"
		courseCode := "ZZZZZZ"

		//Output
		output := editAnnouncement(courseCode,description)

		//Compare output to expected output
		assert.Equal(t, output, "success")
		assert.Equal(t, description,getAnnouncement(courseCode))
	})

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "ZZZZZZ")
	if err != nil {
		panic(err)
	}
}
