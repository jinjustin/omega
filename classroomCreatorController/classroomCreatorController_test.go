package classroomcreatorcontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"encoding/json"
	"omega/database"
	"database/sql"
)

func Test_createNewClassroom(t *testing.T) {
	t.Run("Unit Test 001: Correct Information", func(t *testing.T) {
		//Input
		classCode := "99999999"
		className := "Computer Programing"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"
		//Expected Output
		expected := classroom.Classroom{
			ClassID: "000001",
			ClassCode: "99999999",
			ClassName: "Computer Programing",
			Year: "2563",
			Permission: "Private",}
		//Output
		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)
		//Compare output to expected output 
		assert.Equal(t,expected.ClassName,output.ClassName)
		assert.Equal(t,expected.ClassCode,output.ClassCode)
		assert.Equal(t,expected.Year,output.Year)
		assert.Equal(t,expected.Permission,output.Permission)
	 })

	 t.Run("Unit Test 002: Wrong Class Code", func(t *testing.T) {

		className := "Computer Programing"
		classCode := "9999999"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,"",output.ClassName)
		assert.Equal(t,"ClassCode Error",output.ClassCode)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 003: Already Existed Class in same year", func(t *testing.T) {

		className := "Computer Programing"
		classCode := "99999999"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,"",output.ClassName)
		assert.Equal(t,"ClassCode Error",output.ClassCode)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 004: Already Existed Class in different year", func(t *testing.T) {

		//input
		className := "Computer Programing"
		classCode := "99999999"
		year := "2564"
		permission := "Private"
		userID := "a2yQbN"
		//Expected output
		expected := classroom.Classroom{
			ClassID: "000001",
			ClassCode: "99999999",
			ClassName: "Computer Programing",
			Year: "2564",
			Permission: "Private",}

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,expected.ClassName,output.ClassName)
		assert.Equal(t,expected.ClassCode,output.ClassCode)
		assert.Equal(t,expected.Year,output.Year)
		assert.Equal(t,expected.Permission,output.Permission)
	 })

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()

	sqlStatement := `DELETE FROM class WHERE classcode=$1;`
	_, err = db.Exec(sqlStatement, "99999999")
	if err != nil {
	panic(err)
	}

	sqlStatement = `DELETE FROM userinclass WHERE userid=$1;`
	_, err = db.Exec(sqlStatement, "a2yQbN")
	if err != nil {
	panic(err)
	}
 }