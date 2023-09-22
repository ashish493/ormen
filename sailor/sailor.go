// Analogical to Worker
package sailor

import (
	"fmt"

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
func (s *Sailor) RunTask() {
	fmt.Println("Task will run ")
}

func (s *Sailor) StopTask() {
	fmt.Println("Task will be stopped started")
}

func (s *Sailor) StartTask() {
	fmt.Println("Task will get started")
}
