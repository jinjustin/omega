package classroomcreatorcontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"encoding/json"
)

func Test_createNewClassroom(t *testing.T) {
	t.Run("Unit Test 001: Correct Information", func(t *testing.T) {
		//Input
		classCode := "01076002"
		className := "Computer Programing"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"
		//Expected Output
		expected := classroom.Classroom{
			ClassID: "000001",
			ClassCode: "01076002",
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
		classCode := "0107602"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,"",output.ClassName)
		assert.Equal(t,"Error",output.ClassCode)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 003: Already Existed Class in same year", func(t *testing.T) {

		className := "Computer Programing"
		classCode := "01076002"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,"",output.ClassName)
		assert.Equal(t,"Error",output.ClassCode)
		assert.Equal(t,"",output.Year)
		assert.Equal(t,"",output.Permission)
	 })

	 t.Run("Unit Test 004: Already Existed Class in different year", func(t *testing.T) {

		className := "Computer Programing"
		classCode := "01076002"
		year := "2564"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,className,output.ClassName)
		assert.Equal(t,classCode,output.ClassCode)
		assert.Equal(t,year,output.Year)
		assert.Equal(t,permission,output.Permission)
	 })
 }