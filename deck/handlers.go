package deck

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashish493/ormen/sail"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)

	te := sail.SailEvent{}
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Printf(msg)
		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}

	a.Manager.AddTask(te)
	log.Printf("Added sail %v\n", te.Sail.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Sail)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Manager.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Printf("No taskID passed in request.\n")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	_, ok := a.Manager.TaskDb[tID]
	if !ok {
		log.Printf("No sail with ID %v found", tID)
		w.WriteHeader(404)
	}

	te := sail.SailEvent{
		ID:        uuid.New(),
		State:     sail.Completed,
		Timestamp: time.Now(),
	}
	taskToStop := a.Manager.TaskDb[tID]
	// we need to make a copy so we are not modifying the sail in the datastore
	taskCopy := *taskToStop
	taskCopy.State = sail.Completed
	te.Sail = taskCopy
	a.Manager.AddTask(te)

	log.Printf("Added sail event %v to stop task %v\n", te.ID, taskToStop.ID)
	w.WriteHeader(204)
}
