package classroomdeletercontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"omega/testdb"
	"encoding/json"
)

func Test_deleteClassroom(t *testing.T) {

	c := classroom.Classroom{
		ClassID: "vB2p9U",
		ClassCode: "99999999",
		ClassName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	w := classroom.Classroom{
		ClassID: "OBLMv5",
		ClassCode: "12309854",
		ClassName: "Test2",
		Year: "2564",
		Permission: "Private",
	}

	testdb.ClassroomDB = append(testdb.ClassroomDB,c)

	t.Run("Unit Test 005: Delete Existed Classroom", func(t *testing.T) {

		var output classroom.Classroom
		json.Unmarshal(DeleteClassroom(c.ClassID),&output)

		assert.Equal(t,c,output)
	 })

	 t.Run("Unit Test 006: Delete Non-Existed Classroom", func(t *testing.T) {

		var output classroom.Classroom
		json.Unmarshal(DeleteClassroom(w.ClassID),&output)

		assert.Equal(t,"Can't find.",output.ClassID)
	 })
 }