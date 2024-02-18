package model

import (
	"fmt"
	"sort"
	"strings"
)

const TaskStringFormat = "%d:%d"

type Task struct {
	Id    int  `json:"id"`
	Count uint `json:"count"`
}

func (t Task) String() string {
	return fmt.Sprintf(TaskStringFormat, t.Id, t.Count)
}

type Tasks []Task

func (ts Tasks) String() string {

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Id < ts[j].Id
	})

	lines := make([]string, 0, len(ts))
	for _, task := range ts {
		lines = append(lines, task.String())
	}

	return strings.Join(lines, "\n")
}
