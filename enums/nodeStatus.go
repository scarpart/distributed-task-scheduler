package enums

type NodeStatus int

const (
	OnFree NodeStatus = iota
	OnBusy
	Off
)
