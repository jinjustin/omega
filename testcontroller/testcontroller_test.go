package testcontroller

import (
	"database/sql"
	"encoding/json"
	"omega/test"
	"omega/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateTest(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role) values  ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `UE0099`,`testcreatetest`,`123456`,`teacher`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Create with all information", func(t *testing.T) {
		//Input
		testInCourse := test.Test{
			TestID: "",
			CourseID: "77777777",
			CourseCode: "CC5563",
			UserID: "UE0099",
			Status: "Publish",
			Name: "test create test",
			Duration: "3",
			Start: "13:00:00",
			Date: "2021-01-21",
			Description: "description",
		}

		//Output
		var output test.Test
		json.Unmarshal(createTest(testInCourse.CourseID,testInCourse.CourseCode,"testcreatetest",testInCourse.Status,testInCourse.Name,testInCourse.Duration,testInCourse.Start,testInCourse.Date,testInCourse.Description), &output)

		//Compare output to expected output
		assert.Equal(t, testInCourse.CourseID,output.CourseID)
		assert.Equal(t, testInCourse.CourseCode,output.CourseCode)
		assert.Equal(t, testInCourse.Date,output.Date)
	})

	t.Run("Unit Test 002: Create with some information", func(t *testing.T) {
		//Input
		testInCourse := test.Test{
			TestID: "",
			CourseID: "77777777",
			CourseCode: "CC5563",
			UserID: "UE0099",
			Status: "Publish",
			Name: "test create test 2",
			Duration: "",
			Start: "",
			Date: "",
			Description: "",
		}

		//Output
		var output test.Test
		json.Unmarshal(createTest(testInCourse.CourseID,testInCourse.CourseCode,"testcreatetest",testInCourse.Status,testInCourse.Name,testInCourse.Duration,testInCourse.Start,testInCourse.Date,testInCourse.Description), &output)

		//Compare output to expected output
		assert.Equal(t, testInCourse.CourseID,output.CourseID)
		assert.Equal(t, testInCourse.CourseCode,output.CourseCode)
		assert.Equal(t, testInCourse.Date,output.Date)
	})

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "UE0099")
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM test WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "UE0099")
	if err != nil {
		panic(err)
	}
}

/*func Test_checkTestTime(t *testing.T) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO test (testid,courseid,coursecode,userid,status,name,duration,starttime,startdate,description) values  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement, `TS9856`,`55555555`,`CC9999`,`000002`,`Publish`,`Test check test time`,`3`,`14:00:00`,`2021-05-30`,`Course Description Here.`)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Don't have same test time", func(t *testing.T) {
		//Input
		courseCode := `CC9999`
		startTime := `14:00:00`
		startDate := `2021-05-30`

		//Output
		output := checkTestTime(courseCode,startTime,startDate)

		//Compare output to expected output
		assert.Equal(t, description,output)
	})

	t.Run("Unit Test 002: Edit Course Description", func(t *testing.T) {
		//Input
		description := "Course description being Edited"
		courseCode := "ZZZZZZ"

		//Output
		output := EditDescription(courseCode,description)

		//Compare output to expected output
		assert.Equal(t, output, "success")
		assert.Equal(t, description,GetDescription(courseCode))
	})

	sqlStatement = `DELETE FROM course WHERE coursecode=$1;`
	_, err = db.Exec(sqlStatement, "ZZZZZZ")
	if err != nil {
		panic(err)
	}
}*/