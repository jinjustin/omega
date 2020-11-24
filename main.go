package main

import (
	"fmt"
	"net/http"
	"omega/coursecontroller"
	//"omega/database"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func createCourse(w http.ResponseWriter, r *http.Request){
	type Input struct{
		courseName string
		courseID string
		year string
		permission string
		userID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursecontroller.CreateCourse(input.courseName,input.courseID,input.year,input.permission,input.userID))
}

func getCourseList(w http.ResponseWriter, r *http.Request){
	type Input struct{
		userID string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetCourseList(input.userID))
}

func deleteCourse(w http.ResponseWriter, r *http.Request){
	type Input struct{
		courseCode string
		userID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursecontroller.DeleteCourse(input.courseCode,input.userID))
}

func test(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    //fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/",test)
	myRouter.HandleFunc("/createclass",createCourse).Methods("POST")
	myRouter.HandleFunc("/getallclass",getCourseList).Methods("POST")
	myRouter.HandleFunc("/deleteclass",deleteCourse).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}