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

func Test_GetTestList(t *testing.T) {

	testInCourse := test.Test{
		TestID: "TEST99",
		CourseID: "77777777",
		CourseCode: "CC5563",
		UserID: "UE0099",
		Status: "publish",
		Name: "test create test",
		Duration: "3",
		Start: "13:00:00",
		Date: "2021-01-21",
		Description: "description",
	}

	testInCourse2 := test.Test{
		TestID: "TEST88",
		CourseID: "77777777",
		CourseCode: "CC5563",
		UserID: "UE0099",
		Status: "private",
		Name: "test create test 2",
		Duration: "3",
		Start: "09:00:00",
		Date: "2021-01-21",
		Description: "description",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users (userid, username, password, role) values  ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `STUDEN`,`studenttestlist`,`123456`,`student`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO users (userid, username, password, role) values  ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, `TEACHE`,`teachertestlist`,`123456`,`teacher`)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO test (testid, courseid, coursecode, userid, status, name, duration, start, date, description) values  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement, testInCourse.TestID, testInCourse.CourseID, testInCourse.CourseCode, testInCourse.UserID,  testInCourse.Status, testInCourse.Name, testInCourse.Duration, testInCourse.Start, testInCourse.Date, testInCourse.Description)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO test (testid, courseid, coursecode, userid, status, name, duration, start, date, description) values  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement, testInCourse2.TestID, testInCourse2.CourseID, testInCourse2.CourseCode, testInCourse2.UserID,  testInCourse2.Status, testInCourse2.Name, testInCourse2.Duration, testInCourse2.Start, testInCourse2.Date, testInCourse2.Description)
	if err != nil {
		panic(err)
	}

	var expect []test.Test

	expect = append(expect, testInCourse)

	t.Run("Unit Test 001: Get Test List Student", func(t *testing.T) {

		//Output
		output := getTestList(testInCourse.CourseID,"studenttestlist")

		//Compare output to expected output
		assert.Equal(t, expect,output)
	})

	expect = append(expect, testInCourse2)

	t.Run("Unit Test 002: Get Test List Teacher", func(t *testing.T) {

		//Output
		output := getTestList(testInCourse.CourseID,"teachertestlist")

		//Compare output to expected output
		assert.Equal(t, expect,output)
	})

	sqlStatement = `DELETE FROM test WHERE testid=$1;`
	_, err = db.Exec(sqlStatement, testInCourse.TestID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM test WHERE testid=$1;`
	_, err = db.Exec(sqlStatement, testInCourse2.TestID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "STUDEN")
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM users WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "TEACHE")
	if err != nil {
		panic(err)
	}

}

func Test_GetTestInfo(t *testing.T) {

	testInCourse := test.Test{
		TestID: "TEST99",
		CourseID: "77777777",
		CourseCode: "CC5563",
		UserID: "UE0099",
		Status: "publish",
		Name: "test create test",
		Duration: "3",
		Start: "13:00:00",
		Date: "2021-01-21",
		Description: "description",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO test (testid, courseid, coursecode, userid, status, name, duration, start, date, description) values  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement, testInCourse.TestID, testInCourse.CourseID, testInCourse.CourseCode, testInCourse.UserID,  testInCourse.Status, testInCourse.Name, testInCourse.Duration, testInCourse.Start, testInCourse.Date, testInCourse.Description)
	if err != nil {
		panic(err)
	}

	t.Run("Unit Test 001: Get Test Information", func(t *testing.T) {

		//Output
		var output test.Test
		json.Unmarshal(getTestInfo(testInCourse.CourseID,testInCourse.Name), &output)


		//Compare output to expected output
		assert.Equal(t, testInCourse, output)
	})

	sqlStatement = `DELETE FROM test WHERE testid=$1;`
	_, err = db.Exec(sqlStatement, testInCourse.TestID)
	if err != nil {
		panic(err)
	}

}