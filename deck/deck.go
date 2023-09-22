//Analogicall to Manager

package deck

import (
	"fmt"

	"github.com/ashish493/ormen/sail"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Deck struct {
	Pending       queue.Queue
	TaskDb        map[string][]sail.Sail
	EventDb       map[string][]sail.SailEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Deck) SelectSailor() {
	fmt.Println("I will select an appropriate worker")
}
func (m *Deck) UpdateSails() {
	fmt.Println("I will update tasks")
}
func (m *Deck) SendWork() {
	fmt.Println("I will send work to workers")
}
