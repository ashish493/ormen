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

	"github.com/ashish493/ormen/mast"
	"github.com/ashish493/ormen/rudder"
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
	WorkerNodes   []*mast.Mast
	Rudder        rudder.Rudder
}

func New(workers []string, schedulerType string) *Deck {
	taskDb := make(map[uuid.UUID]*sail.Sail)
	eventDb := make(map[uuid.UUID]*sail.SailEvent)
	workerTaskMap := make(map[string][]uuid.UUID)
	taskWorkerMap := make(map[uuid.UUID]string)

	var nodes []*mast.Mast
	for worker := range workers {
		workerTaskMap[workers[worker]] = []uuid.UUID{}

		nAPI := fmt.Sprintf("http://%v", workers[worker])
		n := mast.NewNode(workers[worker], nAPI, "worker")
		nodes = append(nodes, n)
	}

	var s rudder.Rudder
	switch schedulerType {
	// case "greedy":
	// 	s = &scheduler.Greedy{Name: "greedy"}
	case "roundrobin":
		s = &rudder.RoundRobin{Name: "roundrobin"}
		// default:
		// 	s = &scheduler.Epvm{Name: "epvm"}
	}

	return &Deck{
		Pending:       *queue.New(),
		Workers:       workers,
		TaskDb:        taskDb,
		EventDb:       eventDb,
		WorkerTaskMap: workerTaskMap,
		TaskWorkerMap: taskWorkerMap,
		WorkerNodes:   nodes,
		Rudder:        s,
	}
}

func (m *Deck) SelectWorker(t sail.Sail) (*mast.Mast, error) {
	candidates := m.Rudder.SelectCandidateNodes(t, m.WorkerNodes)
	if candidates == nil {
		msg := fmt.Sprintf("No available candidates match resource request for task %v", t.ID)
		err := errors.New(msg)
		return nil, err
	}
	scores := m.Rudder.Score(t, candidates)
	if scores == nil {
		return nil, fmt.Errorf("no scores returned to task %v", t)
	}
	selectedNode := m.Rudder.Pick(scores, candidates)

	return selectedNode, nil
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

func (m *Deck) stopTask(worker string, taskID string) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/tasks/%s", worker, taskID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("error creating request to delete task %s: %v", taskID, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error connecting to worker at %s: %v", url, err)
		return
	}

	if resp.StatusCode != 204 {
		log.Printf("Error sending request: %v", err)
		return
	}

	log.Printf("task %s has been scheduled to be stopped", taskID)
}

func (d *Deck) SendWork() {
	if d.Pending.Len() > 0 {
		e := d.Pending.Dequeue()
		te := e.(sail.SailEvent)
		d.EventDb[te.ID] = &te
		log.Printf("Pulled %v off pending queue", te)

		taskWorker, ok := d.TaskWorkerMap[te.Sail.ID]
		if ok {
			persistedTask := d.TaskDb[te.Sail.ID]
			if te.State == sail.Completed && sail.ValidStateTransition(persistedTask.State, te.State) {
				d.stopTask(taskWorker, te.Sail.ID.String())
				return
			}

			log.Printf("invalid request: existing task %s is in state %v and cannot transition to the completed state", persistedTask.ID.String(), persistedTask.State)
			return
		}

		t := te.Sail
		w, err := d.SelectWorker(t)
		if err != nil {
			log.Printf("error selecting worker for task %s: %v", t.ID, err)
			return
		}

		// e := d.Pending.Dequeue()
		// te := e.(sail.SailEvent)
		// t := te.Sail
		d.WorkerTaskMap[w.Name] = append(d.WorkerTaskMap[w.Name], te.Sail.ID)
		d.TaskWorkerMap[t.ID] = w.Name

		url := fmt.Sprintf("http://%s/tasks", w.Name)

		log.Printf("Pulled %v off pending queue", t)
		d.EventDb[te.ID] = &te
		// d.WorkerTaskMap[w] = append(d.WorkerTaskMap[w], te.Sail.ID)
		// d.TaskWorkerMap[t.ID] = w

		t.State = sail.Scheduled
		d.TaskDb[t.ID] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Unable to marshal task object: %v.", t)
		}

		url = fmt.Sprintf("http://%s/tasks", w)
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
