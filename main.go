package main

import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"os/exec"
)

type Job struct {
	Context string `json:"context"`
	Dockerfile string `json:"dockerfile"`
	Destination string `json:"destination"`
}

var Jobs = []Job{}

func returnAllJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Jobs)
}

func createJob (w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job Job
	json.Unmarshal(reqBody, &job)
	Jobs = append(Jobs, job)
	cmd := exec.Command("/kaniko/executor", "--context", job.Context, "--dockerfile", job.Dockerfile, "--destination", job.Destination)
	go cmd.Run()
}

func handleRequests() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/jobs", returnAllJobs).Methods("GET")
	router.HandleFunc("/jobs", createJob).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Starting Kaniqueue server....")
    handleRequests()
}
