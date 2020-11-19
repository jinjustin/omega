package classroomlistcontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"omega/database"
	"database/sql"
)

func Test_deleteClassroom(t *testing.T) {

	c1 := classroom.Classroom{
		ClassID: "vB2p9U",
		ClassCode: "99999999",
		ClassName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u1 := "vBUp15"

	c2 := classroom.Classroom{
		ClassID: "OBLMv5",
		ClassCode: "12309854",
		ClassName: "Test2",
		Year: "2564",
		Permission: "Private",
	}

	u2 := "vBUp15"

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO class (classid,classname,classcode, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c1.ClassID,c1.ClassName, c1.ClassCode, c1.Year, c1.Permission)
		if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO userinclass (classid,userid)VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, c1.ClassID, u1)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO class (classid,classname,classcode, year, permission)VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, c2.ClassID,c2.ClassName, c2.ClassCode, c2.Year, c2.Permission)
		if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO userinclass (classid,userid)VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, c2.ClassID, u2)
	if err != nil {
		panic(err)
	}

	t.Run("get Class List from existed UserID", func(t *testing.T) {

		var expected []classroom.Classroom
		expected = append(expected,c1)
		expected = append(expected,c2)
		var output []classroom.Classroom
		output = GetClassroomList("vBUp15")

		assert.Equal(t,expected,output)
	})

	 t.Run("get Class List from Abcdef UserID", func(t *testing.T) {

		var expected []classroom.Classroom
		var output []classroom.Classroom
		output = GetClassroomList("Abcdef")

		assert.Equal(t,expected,output)
	})

	sqlStatement = `DELETE FROM class WHERE classid=$1;`
	 _, err = db.Exec(sqlStatement, c1.ClassID)
	 if err != nil {
	 	panic(err)
	}

	sqlStatement = `DELETE FROM userinclass WHERE classid=$1 and userid=$2;`
	 _, err = db.Exec(sqlStatement, c1.ClassID,u1)
	 if err != nil {
	 	panic(err)
	}

	sqlStatement = `DELETE FROM class WHERE classid=$1;`
	_, err = db.Exec(sqlStatement, c2.ClassID)
	if err != nil {
		panic(err)
   }

	sqlStatement = `DELETE FROM userinclass WHERE classid=$1 and userid=$2;`
	_, err = db.Exec(sqlStatement, c2.ClassID,u2)
	if err != nil {
		panic(err)
   }
 }