package main

import (
	"github.com/ashish493/ormen/deck"
	"github.com/ashish493/ormen/mast"
	"github.com/ashish493/ormen/sail"
	"github.com/google/uuid"

	"fmt"
	"time"

	"github.com/golang-collections/collections/queue"

	"github.com/ashish493/ormen/sailor"
)

func main() {
	s := sail.Sail{
		ID:     uuid.New(),
		Name:   "Sail-1",
		State:  sail.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	se := sail.SailEvent{
		ID:        uuid.New(),
		State:     sail.Pending,
		Timestamp: time.Now(),
		Sail:      s,
	}

	fmt.Printf("task: %v\n", s)
	fmt.Printf("task event: %v\n", se)

	w := sailor.Sailor{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]sail.Sail),
	}

	fmt.Printf("worker: %v\n", w)
	w.Collections()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := deck.Deck{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]sail.Sail),
		EventDb: make(map[string][]sail.SailEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("manager: %v\n", m)
	m.SelectSailor()
	m.UpdateSails()
	m.SendWork()

	n := mast.Mast{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", n)
}
