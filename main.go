package main

import (
	"fmt"
	"net/http"
	"omega/coursecontroller"
	"omega/coursemembercontroller"
	"omega/teachercontroller"
	"omega/login"

	//"omega/database"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func createCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseName string
		CourseID string
		Year string
		Permission string
		Announcement string
		Description string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	fmt.Println(input)
	w.Write(coursecontroller.CreateCourse(input.CourseName,input.CourseID,input.Year,input.Permission,input.Announcement,input.Description,input.Username))
}

func getCourseList(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Username string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetCourseList(input.Username))
}

func deleteCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursecontroller.DeleteCourse(input.CourseCode,input.Username))
}

func getTeacherInfo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(teachercontroller.GetTeacherInfo(input.Username))
}

func editTeacherInfo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Firstname string
		Surname string
		Email string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(teachercontroller.EditTeacherInfo(input.Firstname,input.Surname,input.Email,input.Username))
}

func getDescription(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetDescription(input.CourseCode))
}

func editDescription(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Description string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.EditDescription(input.CourseCode,input.Description))
}

func getAnnouncement(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetAnnouncement(input.CourseCode))
}

func editAnnouncement(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Announcement string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.EditAnnouncement(input.CourseCode,input.Announcement))
}

func addStudentToCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		StudentID string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.AddStudentToCourse(input.StudentID,input.CourseCode))
}

func addTeacherToCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Username string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.AddTeacherToCourse(input.Username,input.CourseCode))
}

func getStudentInCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.GetStudentInCourse(input.CourseCode))
}

func getTeacherInCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.GetTeacherInCourse(input.CourseCode))
}

func test(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
    fmt.Fprintf(w, "Welcome to the HomePage!")
}

func test2(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		test string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	w.Write(reqBody)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET","POST","DELETE"},
	})

	myRouter.HandleFunc("/",test)
	myRouter.HandleFunc("/test",test2).Methods("POST")
	myRouter.HandleFunc("/createcourse",createCourse).Methods("POST")
	myRouter.HandleFunc("/getcourselist",getCourseList).Methods("POST")
	myRouter.HandleFunc("/deletecourse",deleteCourse).Methods("POST")
	myRouter.HandleFunc("/getteacherinfo",getTeacherInfo).Methods("POST")
	myRouter.HandleFunc("/editteacherinfo",editTeacherInfo).Methods("POST")
	myRouter.HandleFunc("/getdescription",getDescription).Methods("POST")
	myRouter.HandleFunc("/editdescription",editDescription).Methods("POST")
	myRouter.HandleFunc("/getannouncement",getAnnouncement).Methods("POST")
	myRouter.HandleFunc("/editannouncement",editAnnouncement).Methods("POST")
	myRouter.HandleFunc("/addstudent",addStudentToCourse).Methods("POST")
	myRouter.HandleFunc("/addteacher",addTeacherToCourse).Methods("POST")
	myRouter.HandleFunc("/getstudentincourse",getStudentInCourse).Methods("POST")
	myRouter.HandleFunc("/getteacherincourse",getTeacherInCourse).Methods("POST")
	myRouter.HandleFunc("/login", login.Login).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000",c.Handler(myRouter)))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	handleRequests()
}