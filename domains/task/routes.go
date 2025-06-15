package task

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var tasks = []Task{
	{
		ID:          "1",
		Title:       "Task 1",
		Description: "Description 1",
	},
	{
		ID:          "2",
		Title:       "Task 2",
		Description: "Description 2",
	},
	{
		ID:          "3",
		Title:       "Task 3",
		Description: "Description 3",
	},
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("GET /tasks")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(tasks)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling tasks"))
		return
	}

	w.Write(response)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("GET /tasks/:id")

	// extract the id from the path
	id := r.URL.Path[len("/tasks/"):]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid task ID"))
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			response, err := json.Marshal(task)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error marshalling task"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found"))
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading body"))
		return
	}

	// validate the fields
	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error unmarshalling body"))
		return
	}

	// create the task
	tasks = append(tasks, task)

	// prepare response
	response, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling task"))
		return
	}

	// return the task
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
