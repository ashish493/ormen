//Analogical to scheduler

package rudder

type rudder interface {
	SelectCandidateNodes()
	Score()
	Pick()
}
