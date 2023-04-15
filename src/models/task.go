package models

import (
	"github.com/scarpart/distributed-task-scheduler/src/types"
)

type Task struct {
	Id              uint64            `json:"id" gorm="primary_key"`
	TaskName        string            `json:"taskName"`
	TaskDescription int               `json:"taskDescription"`
	Status          int               `json:"status"`
	Priority        int               `json:"priority"`
	Dependencies    types.Uint64Array `json:"dependencies"`
	NodeID          int               `json:"nodeId"`
	Command         string            `json:"command"`
}
