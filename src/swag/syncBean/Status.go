package syncBean

// StatusLevel is the level of status
type StatusLevel int

const (
	FOLLOW StatusLevel = iota
	CANDIDATE
	LEADER
)
