package coursecontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/course"
	"encoding/json"
	"omega/database"
	"database/sql"
)

func Test_createCourse(t *testing.T) {
	t.Run("Unit Test 001: Correct Information", func(t *testing.T) {
		//Input
		courseID := "99999999"
		courseName := "Computer Programing"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"
		//Expected Output
		expected := course.Course{
			CourseCode: "000001",
			CourseID: "99999999",
			CourseName: "Computer Programing",
			Year: "2563",
			Permission: "Private",}
		//Output
		var output course.Course
		json.Unmarshal(CreateCourse(courseName,courseID,year,permission,userID),&output)
		//Compare output to expected output 
		assert.Equal(t,expected.CourseName,output.CourseName)
		assert.Equal(t,expected.CourseID,output.CourseID)
		assert.Equal(t,expected.Year,output.Year)
		assert.Equal(t,expected.Permission,output.Permission)
	 })

	 t.Run("Unit Test 002: Wrong Course ID", func(t *testing.T) {

		courseName := "Computer Programing"
		courseID := "9999999"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output course.Course
		json.Unmarshal(CreateCourse(courseName,courseID,year,permission,userID),&output)

		assert.Equal(t,"",output.CourseName)
		assert.Equal(t,"Course ID Error",output.CourseID)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 003: Already Existed Course in same year", func(t *testing.T) {

		courseID := "99999999"
		courseName := "Computer Programing"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"
		var output course.Course
		json.Unmarshal(CreateCourse(courseName,courseID,year,permission,userID),&output)

		assert.Equal(t,"",output.CourseName)
		assert.Equal(t,"Course ID Error",output.CourseID)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 004: Already Existed Course in different year", func(t *testing.T) {

		//input
		courseName := "Computer Programing"
		courseID := "99999999"
		year := "2564"
		permission := "Private"
		userID := "a2yQbN"
		//Expected output
		expected := course.Course{
			CourseCode: "000001",
			CourseID: "99999999",
			CourseName: "Computer Programing",
			Year: "2564",
			Permission: "Private",}

		var output course.Course
		json.Unmarshal(CreateCourse(courseName,courseID,year,permission,userID),&output)

		assert.Equal(t,expected.CourseName,output.CourseName)
		assert.Equal(t,expected.CourseID,output.CourseID)
		assert.Equal(t,expected.Year,output.Year)
		assert.Equal(t,expected.Permission,output.Permission)
	 })

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `DELETE FROM course WHERE courseid=$1;`
	_, err = db.Exec(sqlStatement, "99999999")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "a2yQbN")
	if err != nil {
	panic(err)
	}
 }

 func Test_deleteCourse(t *testing.T) {

	c := course.Course{
		CourseCode: "vB2p9U",
		CourseID: "99999999",
		CourseName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u := "ZZZZZZ"

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, c.CourseCode,c.CourseID, c.CourseName, c.Year, c.Permission)
		if err != nil {
		panic(err)
		}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, c.CourseCode, u,`teacher`,`join`);
		if err != nil {
		panic(err)
		}

	t.Run("Unit Test 005: Delete Non-Existed Course", func(t *testing.T) {

		//Input
		courseCode := "OBLMv5"
		userID := u

		var output course.Course
		json.Unmarshal(DeleteCourse(courseCode,userID),&output)

		assert.Equal(t,"Can't find.",output.CourseCode)
	})

	t.Run("Unit Test 006: Delete Existed Course but user don't in that course", func(t *testing.T) {

		//Input
		courseCode := c.CourseCode
		userID := "AAAAAA"

		var output course.Course
		json.Unmarshal(DeleteCourse(courseCode,userID),&output)

		assert.Equal(t,"Can't find.",output.CourseCode)
	})

	t.Run("Unit Test 007: Delete Existed Course", func(t *testing.T) {

		//Input
		courseCode := c.CourseCode
		userID := u

		var output course.Course
		json.Unmarshal(DeleteCourse(courseCode,userID),&output)

		assert.Equal(t,c,output)
	})
 }

 func Test_getCourselist(t *testing.T) {

	c1 := course.Course{
		CourseCode: "vB2p9U",
		CourseID: "99999999",
		CourseName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u1 := "Test01"

	c2 := course.Course{
		CourseCode: "OBLMv5",
		CourseID: "12309854",
		CourseName: "Test2",
		Year: "2564",
		Permission: "Private",
	}

	u2 := "Test01"

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c1.CourseCode,c1.CourseID, c1.CourseName, c1.Year, c1.Permission)
		if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, c1.CourseCode, u1, `teacher`, `join`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO course (coursecode,courseid,coursename, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c2.CourseCode,c2.CourseID, c2.CourseName, c2.Year, c2.Permission)
		if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO coursemember (coursecode,userid,role,status)VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, c2.CourseCode, u2, `teacher`, `join`)
	if err != nil {
		panic(err)
	}

	t.Run("get Course List from existed UserID", func(t *testing.T) {

		var expected []course.Course
		expected = append(expected,c1)
		expected = append(expected,c2)
		var output []course.Course
		output = GetCourseList("Test01")

		assert.Equal(t,expected,output)
	})

	 t.Run("get Course List from Abcdef UserID", func(t *testing.T) {

		var expected []course.Course
		var output []course.Course
		output = GetCourseList("Abcdef")

		assert.Equal(t,expected,output)
	})

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	 _, err = db.Exec(sqlStatement, c1.CourseCode)
	 if err != nil {
	 	panic(err)
	}

	sqlStatement = `DELETE FROM coursemember WHERE coursecode=$1 and userid=$2;`
	 _, err = db.Exec(sqlStatement, c1.CourseCode,u1)
	 if err != nil {
	 	panic(err)
	}

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, c2.CourseCode)
	if err != nil {
		panic(err)
   }

	sqlStatement = `DELETE FROM coursemember WHERE coursecode=$1 and userid=$2;`
	_, err = db.Exec(sqlStatement, c2.CourseCode,u2)
	if err != nil {
		panic(err)
   }
 }