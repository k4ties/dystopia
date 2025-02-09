package rank

type Priority int

const (
	PriorityLoser Priority = -1
)

const (
	PriorityDefault Priority = iota * 10
	PriorityUnusual
	PrioritySpecial
	PriorityHelper
	PriorityModerator
	PriorityAdmin
	PriorityManager
	PriorityOwner Priority = iota * 1000
)
