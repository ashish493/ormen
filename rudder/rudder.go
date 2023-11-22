//Analogical to scheduler

package rudder

import (
	"github.com/ashish493/ormen/mast"
	"github.com/ashish493/ormen/sail"
)

type Rudder interface {
	SelectCandidateNodes(t sail.Sail, nodes []*mast.Mast) []*mast.Mast
	Score(t sail.Sail, nodes []*mast.Mast) map[string]float64
	Pick(scores map[string]float64, candidates []*mast.Mast) *mast.Mast
}

type RoundRobin struct {
	Name       string
	LastWorker int
}

func (r *RoundRobin) SelectCandidateNodes(t sail.Sail, nodes []*mast.Mast) []*mast.Mast {
	return nodes
}

func (r *RoundRobin) Score(t sail.Sail, nodes []*mast.Mast) map[string]float64 {
	nodeScores := make(map[string]float64)

	var newWorker int
	if r.LastWorker+1 < len(nodes) {
		newWorker = r.LastWorker + 1
		r.LastWorker++
	} else {
		newWorker = 0
		r.LastWorker = 0
	}

	for idx, node := range nodes {
		if idx == newWorker {
			nodeScores[node.Name] = 0.1
		} else {
			nodeScores[node.Name] = 1.0
		}
	}

	return nodeScores
}

func (r *RoundRobin) Pick(scores map[string]float64, candidates []*mast.Mast) *mast.Mast {
	var bestNode *mast.Mast
	var lowestScore float64
	for idx, node := range candidates {
		if idx == 0 {
			bestNode = node
			lowestScore = scores[node.Name]
			continue
		}

		if scores[node.Name] < lowestScore {
			bestNode = node
			lowestScore = scores[node.Name]
		}
	}

	return bestNode
}
