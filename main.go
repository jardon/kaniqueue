package main

import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"os/exec"
	"os"
	"io"
	"bytes"
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
	fmt.Printf("Job recieved: %s", job.Destination)
	cmd := exec.Command("/kaniko/executor", "--context", "git://" + job.Context, "--dockerfile", job.Dockerfile, "--destination", job.Destination)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	log.Println(stdBuffer.String())
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
