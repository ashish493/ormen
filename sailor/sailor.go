// Analogical to Worker
package sailor

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ashish493/ormen/sail"
	"github.com/ashish493/ormen/stats"
	"github.com/ashish493/ormen/store"
	"github.com/golang-collections/collections/queue"
)

type Sailor struct {
	Name  string
	Queue queue.Queue
	Db    store.Store
	//Db        map[uuid.UUID]*task.Task
	Stats     *stats.Stats
	SailCount int
}

func New(name string, taskDbType string) *Sailor {
	w := Sailor{
		Name:  name,
		Queue: *queue.New(),
	}

	var s store.Store
	var err error
	switch taskDbType {
	case "memory":
		s = store.NewInMemoryTaskStore()
	case "persistent":
		filename := fmt.Sprintf("%s_tasks.db", name)
		s, err = store.NewTaskStore(filename, 0600, "tasks")
	}
	if err != nil {
		log.Printf("unable to create new task store: %v", err)
	}
	w.Db = s
	return &w
}

func (w *Sailor) RunTasks() {
	for {
		if w.Queue.Len() != 0 {
			result := w.runTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}

}

func (s *Sailor) CollectStats() {
	for {
		log.Println("Collecting stats")
		s.Stats = stats.GetStats()
		s.SailCount = s.Stats.TaskCount
		time.Sleep(15 * time.Second)
	}
}

func (w *Sailor) GetTasks() []*sail.Sail {
	taskList, err := w.Db.List()
	if err != nil {
		log.Printf("error getting list of tasks: %v", err)
		return nil
	}

	return taskList.([]*sail.Sail)
}

func (s *Sailor) Collections() {
	fmt.Println("Metrics Collection Func")
}
func (s *Sailor) AddTask(t sail.Sail) {
	s.Queue.Enqueue(t)
}

func (w *Sailor) runTask() sail.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("[worker] No sails in the queue")
		return sail.DockerResult{Error: nil}
	}

	taskQueued := t.(sail.Sail)
	fmt.Printf("[worker] Found task in queue: %v:\n", taskQueued)

	err := w.Db.Put(taskQueued.ID.String(), &taskQueued)
	if err != nil {
		msg := fmt.Errorf("error storing task %s: %v", taskQueued.ID.String(), err)
		log.Println(msg)
		return sail.DockerResult{Error: msg}
	}

	result, err := w.Db.Get(taskQueued.ID.String())
	if err != nil {
		msg := fmt.Errorf("error getting task %s from database: %v", taskQueued.ID.String(), err)
		log.Println(msg)
		return sail.DockerResult{Error: msg}
	}

	taskPersisted := *result.(*sail.Sail)

	if taskPersisted.State == sail.Completed {
		return w.StopTask(taskPersisted)
	}

	var dockerResult sail.DockerResult
	if sail.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case sail.Scheduled:
			if taskQueued.ContainerID != "" {
				dockerResult = w.StopTask(taskQueued)
				if dockerResult.Error != nil {
					log.Printf("%v\n", dockerResult.Error)
				}
			}
			dockerResult = w.StartTask(taskQueued)
		default:
			fmt.Printf("This is a mistake. taskPersisted: %v, taskQueued: %v\n", taskPersisted, taskQueued)
			dockerResult.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		dockerResult.Error = err
		return dockerResult
	}
	return dockerResult
}

func (w *Sailor) StartTask(t sail.Sail) sail.DockerResult {
	config := sail.NewConfig(&t)
	d := sail.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Err running sail %v: %v\n", t.ID, result.Error)
		t.State = sail.Failed
		w.Db.Put(t.ID.String(), &t)
		return result
	}

	t.ContainerID = result.ContainerId
	t.State = sail.Running
	w.Db.Put(t.ID.String(), &t)

	return result
}

func (w *Sailor) StopTask(t sail.Sail) sail.DockerResult {
	config := sail.NewConfig(&t)
	d := sail.NewDocker(config)

	stopResult := d.Stop(t.ContainerID)
	if stopResult.Error != nil {
		log.Printf("%v\n", stopResult.Error)
	}
	removeResult := d.Remove(t.ContainerID)
	if removeResult.Error != nil {
		log.Printf("%v\n", removeResult.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = sail.Completed
	w.Db.Put(t.ID.String(), &t)
	log.Printf("Stopped and removed container %v for sail %v\n", t.ContainerID, t.ID)

	return removeResult
}

func (w *Sailor) InspectTask(t sail.Sail) sail.DockerInspectResponse {
	config := sail.NewConfig(&t)
	d := sail.NewDocker(config)
	return d.Inspect(t.ContainerID)
}

func (w *Sailor) UpdateTasks() {
	for {
		log.Println("Checking status of tasks")
		w.updateTasks()
		log.Println("Task updates completed")
		log.Println("Sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}
}

func (w *Sailor) updateTasks() {
	// for each task in the worker's datastore:
	// 1. call InspectTask method
	// 2. verify task is in running state
	// 3. if task is not in running state, or not running at all, mark task as `failed`
	tasks, err := w.Db.List()
	if err != nil {
		log.Printf("error getting list of tasks: %v", err)
		return
	}
	for _, t := range tasks.([]*sail.Sail) {
		if t.State == sail.Running {
			resp := w.InspectTask(*t)
			if resp.Error != nil {
				fmt.Printf("ERROR: %v", resp.Error)
			}

			if resp.Container == nil {
				log.Printf("No container for running task %s", t.ID)
				t.State = sail.Failed
				w.Db.Put(t.ID.String(), t)
			}

			if resp.Container.State.Status == "exited" {
				log.Printf("Container for task %s in non-running state %s", t.ID, resp.Container.State.Status)
				t.State = sail.Failed
				w.Db.Put(t.ID.String(), t)
			}

			// task is running, update exposed ports
			t.HostPorts = resp.Container.NetworkSettings.NetworkSettingsBase.Ports
			w.Db.Put(t.ID.String(), t)
		}
	}
}
