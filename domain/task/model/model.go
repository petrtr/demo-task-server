package model

import (
	"fmt"
	"strings"
)

type Task struct {
	Id    int  `json:"id"`
	Count uint `json:"count"`
}

func (t Task) String() string {
	return fmt.Sprintf("%d:%d", t.Id, t.Count)
}

type Tasks []Task

func (ts Tasks) String() string {

	lines := make([]string, 0, len(ts))
	for _, task := range ts {
		lines = append(lines, task.String())
	}

	return strings.Join(lines, "\n")
}
