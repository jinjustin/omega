package main

import (
	"fmt"
	"net/http"
	"omega/classroomcreatorcontroller"
	//"omega/database"
	"omega/classroomdeletercontroller"
	"omega/classroomlistcontroller"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func createNewClass(w http.ResponseWriter, r *http.Request){
	type Input struct{
		ClassCode string
		ClassName string
		Year string
		Permission string
		UserID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(classroomcreatorcontroller.CreateNewClass(input.ClassName,input.ClassCode,input.Year,input.Permission,input.UserID))
}

func getClassList(w http.ResponseWriter, r *http.Request){
	type Input struct{
		UserID string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(classroomlistcontroller.GetClassroomList(input.UserID))
}

func deleteClass(w http.ResponseWriter, r *http.Request){
	type Input struct{
		ClassID string
		UserID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(classroomdeletercontroller.DeleteClassroom(input.ClassID,input.UserID))
}

func test(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    //fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/",test)
	myRouter.HandleFunc("/createclass",createNewClass).Methods("POST")
	myRouter.HandleFunc("/getallclass",getClassList).Methods("POST")
	myRouter.HandleFunc("/deleteclass",deleteClass).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}