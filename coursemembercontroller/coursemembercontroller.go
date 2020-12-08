package coursemembercontroller

import(
	"omega/student"
	"omega/teacher"
	"omega/database"
	"database/sql"
)

//AddStudentToCourse is ฟังก์ชันใช้สำหรับให้ผู้สอนเชิญผู้เรียนเข้าไปใน Course หรือผู้เรียนต้องการเข้าร่วม Course
func AddStudentToCourse(studentID string,courseCode string) []byte{

	var userID string
	var firstName string
	var surName string

	s := student.Student{
		UserID: "Can't Invite this student",
		StudentID: "",
		Username: "",
		Password: "",
		Firstname: "",
		Surname: "",
		Email: "",
	}

   db, err := sql.Open("postgres", database.PsqlInfo())
   if err != nil {
	   panic(err)
   }
   defer db.Close()

   sqlStatement := `SELECT userid,firstname,surname FROM student WHERE studentid=$1;`
   rows,err := db.Query(sqlStatement, studentID)
   if err != nil {
	   panic(err)
   }
   defer rows.Close()

   for rows.Next() {
	err = rows.Scan(&userID,&firstName,&surName)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
	panic(err)
	}

	if(checkMemberInCourse(userID,courseCode)){
		s = student.Student{
			UserID: "",
			StudentID: studentID,
			Username: "",
			Password: "",
			Firstname: firstName,
			Surname: surName,
			Email: "",
		}
		
		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

			_, err = db.Exec(sqlStatement, courseCode, userID, "student", "pending")
			if err != nil {
			panic(err)
			}
	}

	return s.GetStudentDetail()
}

//AddTeacherToCourse is ฟังก์ชันที่ใช้เมื่อ ผู้สอน ต้องการเข้าร่วม Course
func AddTeacherToCourse(userID string,courseCode string) []byte{

	var firstName string
	var surName string

	t := teacher.Teacher{
		UserID: "Can't Join this course",
		Username: "",
		Password: "",
		Firstname: "",
		Surname: "",
		Email: "",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT firstname,surname FROM teacher WHERE userid=$1;`
	rows,err := db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
 
	for rows.Next() {
	 err = rows.Scan(&firstName,&surName)
		 if err != nil {
			 panic(err)
		 }
	 }
	 err = rows.Err()
	 if err != nil {
	 panic(err)
	 }

	if(checkMemberInCourse(userID,courseCode)){
		t = teacher.Teacher{
			UserID: "",
			Username: "",
			Password: "",
			Firstname: firstName,
			Surname: surName,
			Email: "",
		}

		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

			_, err = db.Exec(sqlStatement, courseCode, userID, "teacher", "pending")
			if err != nil {
			panic(err)
			}
	}

	return t.GetTeacherDetail()
}

//ApproveJoinCourse is ฟังก์ชันสำหรับให้ผู้สอนรองรับการเข้าร่วม Course ของผู้เรียน
func ApproveJoinCourse(){

}

func checkMemberInCourse(userID string,courseCode string) bool{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var status string

	sqlStatement := `SELECT status FROM coursemember WHERE userid=$1 and coursecode=$2;`
	row := db.QueryRow(sqlStatement, userID,courseCode)
	err = row.Scan(&status)
	switch err {
	case sql.ErrNoRows: return true
	case nil: return false
	default: panic(err)
	}
}