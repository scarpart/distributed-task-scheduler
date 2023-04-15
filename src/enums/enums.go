package enums

type Status int

const (
	Running = iota
	Ready
	Blocked
	Done
)

type Priority int

const (
	Urgent = iota
	Important
	Normal
	NotImportant
)
