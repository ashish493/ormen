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
	SailCount int
}

func (s *Sailor) Collections() {
	fmt.Println("Metrics Collection Func")
}
func (s *Sailor) AddTask(t sail.Sail) {
	s.Queue.Enqueue(t)
}

func (s *Sailor) RunTask() sail.DockerResult {
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
