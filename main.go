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
	fmt.Println("Endpoint Hit: returnAllJobs")
    json.NewEncoder(w).Encode(Jobs)
}

func createJob (w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job Job
	json.Unmarshal(reqBody, &job)
	Jobs = append(Jobs, job)
	fmt.Printf("Job recieved: %s\n", job.Destination)
}

func runJob(job Job) {
	cmd := exec.Command("/kaniko/executor", "--context", "git://" + job.Context, "--dockerfile", job.Dockerfile, "--destination", job.Destination)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	cmd.Run()

	err := os.RemoveAll("/kaniko/buildcontext")
    if err != nil {
        log.Fatal(err)
    }

	log.Println(stdBuffer.String())
	_, Jobs = Jobs[0], Jobs[1:]
}

func processRequests() {
	for {
		if (len(Jobs) > 0) {
			runJob(Jobs[0])
		}
	}
}

func handleRequests() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/jobs", returnAllJobs).Methods("GET")
	router.HandleFunc("/jobs", createJob).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Starting Kaniqueue server....")
	go processRequests()
    handleRequests()
}
