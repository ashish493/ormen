//Analogicall to Manager

package deck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ashish493/ormen/sail"
	"github.com/ashish493/ormen/sailor"
	"github.com/docker/go-connections/nat"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Deck struct {
	Pending       queue.Queue
	TaskDb        map[uuid.UUID]*sail.Sail
	EventDb       map[uuid.UUID]*sail.SailEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
	LastWorker    int
}

func New(workers []string) *Deck {
	taskDb := make(map[uuid.UUID]*sail.Sail)
	eventDb := make(map[uuid.UUID]*sail.SailEvent)
	workerTaskMap := make(map[string][]uuid.UUID)
	taskWorkerMap := make(map[uuid.UUID]string)
	for worker := range workers {
		workerTaskMap[workers[worker]] = []uuid.UUID{}
	}

	return &Deck{
		Pending:       *queue.New(),
		Workers:       workers,
		TaskDb:        taskDb,
		EventDb:       eventDb,
		WorkerTaskMap: workerTaskMap,
		TaskWorkerMap: taskWorkerMap,
	}
}

func (d *Deck) SelectWorker() string {
	fmt.Println("I will select an appropriate sailor.ErrResponse")
	var newWorker int
	if d.LastWorker+1 < len(d.Workers) {
		newWorker = d.LastWorker + 1
		d.LastWorker++
	} else {
		newWorker = 0
		d.LastWorker = 0
	}

	return d.Workers[newWorker]
}

func (m *Deck) UpdateTasks() {
	for {
		log.Println("Checking for task updates from workers")
		m.updateTasks()
		log.Println("Task updates completed")
		log.Println("Sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}
}

func (m *Deck) updateTasks() {
	for _, worker := range m.Workers {
		log.Printf("Checking worker %v for task updates", worker)
		url := fmt.Sprintf("http://%s/tasks", worker)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error connecting to %v: %v", worker, err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error sending request: %v", err)
		}

		d := json.NewDecoder(resp.Body)
		var tasks []*sail.Sail
		err = d.Decode(&tasks)
		if err != nil {
			log.Printf("Error unmarshalling tasks: %s", err.Error())
		}

		for _, t := range tasks {
			log.Printf("Attempting to update task %v", t.ID)

			_, ok := m.TaskDb[t.ID]
			if !ok {
				log.Printf("Task with ID %s not found\n", t.ID)
				return
			}
			if m.TaskDb[t.ID].State != t.State {
				m.TaskDb[t.ID].State = t.State
			}

			m.TaskDb[t.ID].StartTime = t.StartTime
			m.TaskDb[t.ID].FinishTime = t.FinishTime
			m.TaskDb[t.ID].ContainerID = t.ContainerID
		}

	}
}

func (m *Deck) ProcessTasks() {
	for {
		log.Println("Processing any tasks in the queue")
		m.SendWork()
		log.Println("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

func (d *Deck) SendWork() {
	if d.Pending.Len() > 0 {
		w := d.SelectWorker()

		e := d.Pending.Dequeue()
		te := e.(sail.SailEvent)
		t := te.Sail
		log.Printf("Pulled %v off pending queue", t)
		d.EventDb[te.ID] = &te
		d.WorkerTaskMap[w] = append(d.WorkerTaskMap[w], te.Sail.ID)
		d.TaskWorkerMap[t.ID] = w

		t.State = sail.Scheduled
		d.TaskDb[t.ID] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Unable to marshal task object: %v.", t)
		}

		url := fmt.Sprintf("http://%s/tasks", w)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error connecting to %v: %v", w, err)
			d.Pending.Enqueue(t)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := sailor.ErrResponse{}
			err := d.Decode(&e)
			if err != nil {
				fmt.Printf("Error decoding response: %s\n", err.Error())
				return
			}
			log.Printf("Response error (%d): %s", e.HTTPStatusCode, e.Message)
			return
		}

		t = sail.Sail{}
		err = d.Decode(&t)
		if err != nil {
			fmt.Printf("Error decoding response: %s\n", err.Error())
			return
		}
		log.Printf("%#v\n", t)
	} else {
		log.Println("No work in the queue")
	}
}

func (m *Deck) GetTasks() []*sail.Sail {
	tasks := []*sail.Sail{}
	for _, t := range m.TaskDb {
		tasks = append(tasks, t)
	}
	return tasks
}

func (m *Deck) AddTask(te sail.SailEvent) {
	m.Pending.Enqueue(te)
}

func (m *Deck) checkTaskHealth(t sail.Sail) error {
	log.Printf("Calling health check for task %s: %s\n", t.ID, t.HealthCheck)

	w := m.TaskWorkerMap[t.ID]
	hostPort := getHostPort(t.HostPorts)
	worker := strings.Split(w, ":")
	if hostPort == nil {
		log.Printf("Have not collected task %s host port yet. Skipping.\n", t.ID)
		return nil
	}
	url := fmt.Sprintf("http://%s:%s%s", worker[0], *hostPort, t.HealthCheck)
	log.Printf("Calling health check for task %s: %s\n", t.ID, url)
	resp, err := http.Get(url)
	if err != nil {
		msg := fmt.Sprintf("[manager] Error connecting to health check %s", url)
		log.Println(msg)
		return errors.New(msg)
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Error health check for task %s did not return 200\n", t.ID)
		log.Println(msg)
		return errors.New(msg)
	}

	log.Printf("Task %s health check response: %v\n", t.ID, resp.StatusCode)

	return nil
}
func getHostPort(ports nat.PortMap) *string {
	for k, _ := range ports {
		return &ports[k][0].HostPort
	}
	return nil
}

func (m *Deck) DoHealthChecks() {
	for {
		log.Println("Performing task health check")
		m.doHealthChecks()
		log.Println("Task health checks completed")
		log.Println("Sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}
}

func (m *Deck) doHealthChecks() {
	for _, t := range m.TaskDb {
		if t.State == sail.Running && t.RestartCount < 3 {
			err := m.checkTaskHealth(*t)
			if err != nil {
				if t.RestartCount < 3 {
					m.restartTask(t)
				}
			}
		} else if t.State == sail.Failed && t.RestartCount < 3 {
			m.restartTask(t)
		}
	}
}

func (m *Deck) restartTask(t *sail.Sail) {
	// Get the worker where the task was running
	w := m.TaskWorkerMap[t.ID]
	t.State = sail.Scheduled
	t.RestartCount++
	// We need to overwrite the existing task to ensure it has
	// the current state
	m.TaskDb[t.ID] = t

	te := sail.SailEvent{
		ID:        uuid.New(),
		State:     sail.Running,
		Timestamp: time.Now(),
		Sail:      *t,
	}
	data, err := json.Marshal(te)
	if err != nil {
		log.Printf("Unable to marshal task object: %v.", t)
	}

	url := fmt.Sprintf("http://%s/tasks", w)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error connecting to %v: %v", w, err)
		m.Pending.Enqueue(t)
		return
	}

	d := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		e := sailor.ErrResponse{}
		err := d.Decode(&e)
		if err != nil {
			fmt.Printf("Error decoding response: %s\n", err.Error())
			return
		}
		log.Printf("Response error (%d): %s", e.HTTPStatusCode, e.Message)
		return
	}

	newTask := sail.Sail{}
	err = d.Decode(&newTask)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
		return
	}
	log.Printf("%#v\n", t)
}
