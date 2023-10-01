// Analogical to Task
package sail

import (
	"time"

	"github.com/docker/go-connections/nat" //Utility package to interact with Docker network connections
	"github.com/google/uuid"               // package to generate UUIDs
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Sail struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type SailEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Sail      Sail
}
