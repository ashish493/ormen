// Analogical to Worker
package sailor

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ashish493/ormen/sail"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Sailor struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]sail.Sail
	Stats     *Stats
	SailCount int
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
		s.Stats = GetStats()
		s.SailCount = s.Stats.TaskCount
		time.Sleep(15 * time.Second)
	}
}

func (s *Sailor) GetTasks() []*sail.Sail {
	tasks := []*sail.Sail{}
	for _, t := range s.Db {
		tasks = append(tasks, &t)
	}
	return tasks
}

func (s *Sailor) Collections() {
	fmt.Println("Metrics Collection Func")
}
func (s *Sailor) AddTask(t sail.Sail) {
	s.Queue.Enqueue(t)
}

func (s *Sailor) runTask() sail.DockerResult {
	fmt.Println("Task will run ")
	t := s.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return sail.DockerResult{Error: nil}
	}

	taskQueued := t.(sail.Sail)

	taskPersisted := s.Db[taskQueued.ID]
	fmt.Println(taskQueued, "taskqueued")
	// if taskPersisted == nil {
	// 	taskPersisted = taskQueued
	// 	s.Db[taskQueued.ID] = taskQueued
	// }

	var result sail.DockerResult
	if sail.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case sail.Scheduled:
			result = s.StartTask(taskQueued)
		case sail.Completed:
			result = s.StopTask(taskQueued)
		default:
			result.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
		return result
	}
	return result

}

func (s *Sailor) StopTask(t sail.Sail) sail.DockerResult {
	config := sail.NewConfig(&t)
	d := sail.NewDocker(config)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v", t.ContainerID, result.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = sail.Completed
	s.Db[t.ID] = t
	log.Printf("Stopped and removed container %v for sail %v", t.ContainerID, t.ID)

	return result
}

func (s *Sailor) StartTask(t sail.Sail) sail.DockerResult {
	fmt.Println("Task will get started")
	t.StartTime = time.Now().UTC()
	config := sail.NewConfig(&t)
	d := sail.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Err running sail %v: %v\n", t.ID, result.Error)
		t.State = sail.Failed
		s.Db[t.ID] = t
		return result
	}

	t.ContainerID = result.ContainerId
	t.State = sail.Running
	s.Db[t.ID] = t

	return result
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
	for id, t := range w.Db {
		if t.State == sail.Running {
			resp := w.InspectTask(t)
			if resp.Error != nil {
				fmt.Printf("ERROR: %v", resp.Error)
			}

			if resp.Container == nil {
				log.Printf("No container for running task %s", id)
				if state, ok := w.Db[id]; ok {
					state.State = sail.Failed
				}

			}

			if resp.Container.State.Status == "exited" {
				log.Printf("Container for task %s in non-running state %s", id, resp.Container.State.Status)
				if state, ok := w.Db[id]; ok {
					state.State = sail.Failed
				}
			}

			// task is running, update exposed ports
			if state, ok := w.Db[id]; ok {
				state.HostPorts = resp.Container.NetworkSettings.NetworkSettingsBase.Ports
			}
			// w.Db[id].HostPorts = resp.Container.NetworkSettings.NetworkSettingsBase.Ports
		}
	}
}
