package syncBean

type StatusLevel int

const (
	FOLLOW StatusLevel = iota
	CANDIDATE
	LEADER
)
