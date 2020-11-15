package classroomcreatorcontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"encoding/json"
)

func Test_createNewClassroom(t *testing.T) {
	t.Run("Correct Information", func(t *testing.T) {

		className := "Computer Programing"
		classCode := "01076002"
		year := "2563"
		permission := "Private"
		userID := "a2yQbN"

		var output classroom.Classroom
		json.Unmarshal(CreateNewClass(className,classCode,year,permission,userID),&output)

		assert.Equal(t,className,output.ClassName)
		assert.Equal(t,classCode,output.ClassCode)
		assert.Equal(t,year,output.Year)
		assert.Equal(t,permission,output.Permission)
	 })

	 t.Run("Wrong Class Code", func(t *testing.T) {

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

	 t.Run("Already Existed Class", func(t *testing.T) {

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

	 t.Run("Already Existed Class", func(t *testing.T) {

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