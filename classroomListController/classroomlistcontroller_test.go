package classroomlistcontroller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"omega/classroom"
	"omega/testdb"
	"encoding/json"
)

func Test_deleteClassroom(t *testing.T) {

	c1 := classroom.Classroom{
		ClassID: "vB2p9U",
		ClassCode: "99999999",
		ClassName: "Test",
		Year: "2563",
		Permission: "Public",
	}

	u1 := testdb.UserInClass{
		UserID: "vBUp15",
		ClassID: "vB2p9U",
	}

	c2 := classroom.Classroom{
		ClassID: "OBLMv5",
		ClassCode: "12309854",
		ClassName: "Test2",
		Year: "2564",
		Permission: "Private",
	}

	u2 := testdb.UserInClass{
		UserID: "vBUp15",
		ClassID: "OBLMv5",
	}

	c3 := classroom.Classroom{
		ClassID: "QWERTY",
		ClassCode: "01234567",
		ClassName: "Test3",
		Year: "2564",
		Permission: "Public",
	}

	u3 := testdb.UserInClass{
		UserID: "cicd1",
		ClassID: "QWERTY",
	}

	testdb.ClassroomDB = append(testdb.ClassroomDB,c1)
	testdb.UserInClassDB = append(testdb.UserInClassDB,u1)
	testdb.ClassroomDB = append(testdb.ClassroomDB,c2)
	testdb.UserInClassDB = append(testdb.UserInClassDB,u2)
	testdb.ClassroomDB = append(testdb.ClassroomDB,c3)
	testdb.UserInClassDB = append(testdb.UserInClassDB,u3)

	t.Run("get Class List from existed UserID", func(t *testing.T) {

		var expected []classroom.Classroom
		expected = append(expected,c1)
		expected = append(expected,c2)
		var output []classroom.Classroom
		output = GetClassroomList("vBUp15")

		assert.Equal(t,expected,output)
	 })

	 t.Run("get Class List from AAAAAA UserID", func(t *testing.T) {

		var expected []classroom.Classroom
		var output []classroom.Classroom
		output = GetClassroomList("Abcdef")

		assert.Equal(t,expected,output)
	 })

 }