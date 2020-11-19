package classroomdeletercontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"omega/database"
	"encoding/json"
	"database/sql"
)

func Test_deleteClassroom(t *testing.T) {

	c := classroom.Classroom{
		ClassID: "vB2p9U",
		ClassCode: "99999999",
		ClassName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u := "ZZZZZZ"

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `INSERT INTO class (classid,classname,classcode, year, permission)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, c.ClassID,c.ClassName, c.ClassCode, c.Year, c.Permission)
		if err != nil {
		panic(err)
		}

	sqlStatement = `INSERT INTO userinclass (classid,userid)VALUES ($1, $2)`

	_, err = db.Exec(sqlStatement, c.ClassID, u)
		if err != nil {
		panic(err)
		}

	t.Run("Unit Test 005: Delete Non-Existed Classroom", func(t *testing.T) {

		//Input
		classID := "OBLMv5"
		userID := u

		var output classroom.Classroom
		json.Unmarshal(DeleteClassroom(classID,userID),&output)

		assert.Equal(t,"Can't find.",output.ClassID)
	})

	t.Run("Unit Test 006: Delete Existed Classroom but user don't in that classroom", func(t *testing.T) {

		//Input
		classID := c.ClassID
		userID := "AAAAAA"

		var output classroom.Classroom
		json.Unmarshal(DeleteClassroom(classID,userID),&output)

		assert.Equal(t,"Can't find.",output.ClassID)
	})

	t.Run("Unit Test 007: Delete Existed Classroom", func(t *testing.T) {

		//Input
		classID := c.ClassID
		userID := u

		var output classroom.Classroom
		json.Unmarshal(DeleteClassroom(classID,userID),&output)

		assert.Equal(t,c,output)
	})
 }