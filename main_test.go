package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/gorilla/mux"
)

func Test_createClass(t *testing.T) {

	var jsonStr = []byte(`{"ClassCode": "55555555","ClassName": "Test Test","Year": "2563","Permission": "Public","UserID": "bNMv75"}`)

	req, err := http.NewRequest("POST", "/createClass", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/createclass", createNewClass).Methods("POST")
	ts := httptest.NewServer(r)
	
	defer ts.Close()
	res, err := http.Post(ts.URL + "/createclass")
	if err != nil {
	   t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
	   t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
 }