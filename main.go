package main

import (
	"github.com/ashish493/ormen/sail"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"fmt"
	"time"

	"github.com/golang-collections/collections/queue"

	"github.com/ashish493/ormen/sailor"
)

func main() {
	// s := sail.Sail{
	// 	ID:     uuid.New(),
	// 	Name:   "Sail-1",
	// 	State:  sail.Pending,
	// 	Image:  "Image-1",
	// 	Memory: 1024,
	// 	Disk:   1,
	// }

	// se := sail.SailEvent{
	// 	ID:        uuid.New(),
	// 	State:     sail.Pending,
	// 	Timestamp: time.Now(),
	// 	Sail:      s,
	// }

	// fmt.Printf("sail: %v\n", s)
	// fmt.Printf("sail event: %v\n", se)

	// w := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    make(map[uuid.UUID]sail.Sail),
	// }

	// fmt.Printf("worker: %v\n", w)
	// // w.Collections()
	// // w.RunTask()
	// // w.StartTask()
	// // w.StopTask()

	// m := deck.Deck{
	// 	Pending: *queue.New(),
	// 	TaskDb:  make(map[string][]sail.Sail),
	// 	EventDb: make(map[string][]sail.SailEvent),
	// 	Workers: []string{w.Name},
	// }

	// fmt.Printf("manager: %v\n", m)
	// m.SelectSailor()
	// m.UpdateSails()
	// m.SendWork()

	// n := mast.Mast{
	// 	Name:   "Node-1",
	// 	Ip:     "192.168.1.1",
	// 	Cores:  4,
	// 	Memory: 1024,
	// 	Disk:   25,
	// 	Role:   "worker",
	// }
	// fmt.Printf("node: %v\n", n)

	// fmt.Printf("create a test container\n")

	// dockerTask, createResult := createContainer()

	// time.Sleep(time.Second * 5)
	// fmt.Printf("stopping container %s\n", createResult.ContainerId)
	// _ = stopContainer(dockerTask)

	db := make(map[uuid.UUID]sail.Sail)
	w := sailor.Sailor{
		Queue: *queue.New(),
		Db:    db,
	}

	t := sail.Sail{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: sail.Scheduled,
		Image: "strm/helloworld-http",
	}

	t2 := sail.Sail{
		ID:    uuid.New(),
		Name:  "test-container-2",
		State: sail.Scheduled,
		Image: "grafana/grafana",
	}

	// first time the worker will see the sail
	fmt.Println("starting sail")
	w.AddTask(t)
	w.AddTask(t2)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerId

	fmt.Printf("sail %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 50)

	fmt.Printf("stopping sail %s\n", t.ID)
	t.State = sail.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

}

func createContainer() (*sail.Docker, *sail.DockerResult) {
	c := sail.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := sail.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}
	fmt.Printf(
		"Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *sail.Docker) *sail.DockerResult {
	fmt.Println("config values", d.Config.Name)
	result := d.Stop(d.Config.Name)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf(
		"Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}
