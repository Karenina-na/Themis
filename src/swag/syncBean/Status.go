package syncBean

// StatusLevel is the level of status
type StatusLevel int

// StatusLevel
const (
	FOLLOW StatusLevel = iota
	CANDIDATE
	LEADER
)
