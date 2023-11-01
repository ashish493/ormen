package main

import (
	"github.com/ashish493/ormen/cmd"
)

func main() {
	cmd.Execute()

	//chapter -3
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

	//chapter -4
	// db := make(map[uuid.UUID]sail.Sail)
	// w := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    db,
	// }

	// t := sail.Sail{
	// 	ID:    uuid.New(),
	// 	Name:  "test-container-1",
	// 	State: sail.Scheduled,
	// 	Image: "strm/helloworld-http",
	// }

	// t2 := sail.Sail{
	// 	ID:    uuid.New(),
	// 	Name:  "test-container-2",
	// 	State: sail.Scheduled,
	// 	Image: "grafana/grafana",
	// }

	// // first time the worker will see the sail
	// fmt.Println("starting sail")
	// w.AddTask(t)
	// w.AddTask(t2)
	// result := w.RunTask()
	// if result.Error != nil {
	// 	panic(result.Error)
	// }

	// t.ContainerID = result.ContainerId

	// fmt.Printf("sail %s is running in container %s\n", t.ID, t.ContainerID)
	// fmt.Println("Sleepy time")
	// time.Sleep(time.Second * 50)

	// fmt.Printf("stopping sail %s\n", t.ID)
	// t.State = sail.Completed
	// w.AddTask(t)
	// result = w.RunTask()
	// if result.Error != nil {
	// 	panic(result.Error)
	// }

	//chapter 7
	// host := os.Getenv("CUBE_HOST")
	// port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))
	// host := "localhost"
	// port := 8080
	// fmt.Println("Starting Cube worker")

	// w := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    make(map[uuid.UUID]sail.Sail),
	// }

	// api := sailor.Api{Address: host, Port: port, Worker: &w}
	// // fmt.Println(host, port, "host:port")
	// // fmt.Print(api)
	// go w.RunTasks()
	// go w.CollectStats()
	// go api.Start()

	// workers := []string{fmt.Sprintf("%s:%d", host, port)}
	// m := deck.New(workers)

	// for i := 0; i < 3; i++ {
	// 	t := sail.Sail{
	// 		ID:    uuid.New(),
	// 		Name:  fmt.Sprintf("test-container-%d", i),
	// 		State: sail.Scheduled,
	// 		Image: "strm/helloworld-http",
	// 	}
	// 	te := sail.SailEvent{
	// 		ID:    uuid.New(),
	// 		State: sail.Running,
	// 		Sail:  t,
	// 	}
	// 	m.AddTask(te)
	// 	m.SendWork()
	// }

	// go func() {
	// 	for {
	// 		fmt.Printf("[Manager] Updating tasks from %d workers\n", len(m.Workers))
	// 		m.UpdateTasks()
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }()

	// for {
	// 	for _, t := range m.TaskDb {
	// 		fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }

	// whost := os.Getenv("CUBE_WORKER_HOST")
	// wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))
	// whost := "localhost"
	// wport := 8080

	// mhost := "localhost"
	// mport := 8081

	// mhost := os.Getenv("CUBE_MANAGER_HOST")
	// mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	// w1 := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    make(map[uuid.UUID]sail.Sail),
	// }
	// wapi1 := sailor.Api{Address: whost, Port: wport, Worker: &w1}

	// w2 := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    make(map[uuid.UUID]sail.Sail),
	// }
	// wapi2 := sailor.Api{Address: whost, Port: wport + 1, Worker: &w2}

	// w3 := sailor.Sailor{
	// 	Queue: *queue.New(),
	// 	Db:    make(map[uuid.UUID]sail.Sail),
	// }
	// wapi3 := sailor.Api{Address: whost, Port: wport + 2, Worker: &w3}

	// go w1.RunTasks()
	// go w1.UpdateTasks()
	// go wapi1.Start()

	// go w2.RunTasks()
	// go w2.UpdateTasks()
	// go wapi2.Start()

	// go w3.RunTasks()
	// go w3.UpdateTasks()
	// go wapi3.Start()

	// fmt.Println("Starting Cube manager")

	// workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	// m := deck.New(workers, "roundrobin")
	// mapi := deck.Api{Address: mhost, Port: mport, Manager: m}

	// go m.ProcessTasks()
	// go m.UpdateTasks()
	// go m.DoHealthChecks()

	// mapi.Start()
}

// func createContainer() (*sail.Docker, *sail.DockerResult) {
// 	c := sail.Config{
// 		Name:  "test-container-1",
// 		Image: "postgres:13",
// 		Env: []string{
// 			"POSTGRES_USER=cube",
// 			"POSTGRES_PASSWORD=secret",
// 		},
// 	}
// 	dc, _ := client.NewClientWithOpts(client.FromEnv)
// 	d := sail.Docker{
// 		Client: dc,
// 		Config: c,
// 	}
// 	result := d.Run()
// 	if result.Error != nil {
// 		fmt.Printf("%v\n", result.Error)
// 		return nil, nil
// 	}
// 	fmt.Printf(
// 		"Container %s is running with config %v\n", result.ContainerId, c)
// 	return &d, &result
// }

// func stopContainer(d *sail.Docker) *sail.DockerResult {
// 	fmt.Println("config values", d.Config.Name)
// 	result := d.Stop(d.Config.Name)
// 	if result.Error != nil {
// 		fmt.Printf("%v\n", result.Error)
// 		return nil
// 	}
// 	fmt.Printf(
// 		"Container %s has been stopped and removed\n", result.ContainerId)
// 	return &result
// }
