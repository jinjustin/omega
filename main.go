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
		UserID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	fmt.Println(input)
	w.Write(coursecontroller.CreateCourse(input.CourseName,input.CourseID,input.Year,input.Permission,input.UserID))
}

func getCourseList(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		UserID string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetCourseList(input.UserID))
}

func deleteCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		UserID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursecontroller.DeleteCourse(input.CourseCode,input.UserID))
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
    log.Fatal(http.ListenAndServe(":10000",c.Handler(myRouter)))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	handleRequests()
}