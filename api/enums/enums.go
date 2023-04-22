package enums

type TaskStatus int

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

type NodeStatus int

const (
	OnFree = iota
	OnBusy
	Off
)
