// Analogicall to Node

package mast

import (
	"github.com/ashish493/ormen/stats"
)

type Mast struct {
	Name            string
	Ip              string
	Api             string
	Memory          int64
	MemoryAllocated int64
	Disk            int64
	DiskAllocated   int64
	Stats           stats.Stats
	Role            string
	TaskCount       int
}

func NewNode(name string, api string, role string) *Mast {
	return &Mast{
		Name: name,
		Api:  api,
		Role: role,
	}
}
