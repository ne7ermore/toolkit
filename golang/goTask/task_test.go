package goTask

import (
	"fmt"
	"testing"
	"time"
)

const (
	INTERVAL_PERIOD time.Duration = 24 * time.Hour
	HOUR_TO_TASK    int           = 1
	MINUTE_TO_TASK  int           = 0
	SECOND_TO_TASK  int           = 0
)

func TestTask(t *testing.T) {
	tk := NewTask()
	tk.SetTaskTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TASK, MINUTE_TO_TASK, SECOND_TO_TASK, 0, time.Local)).
		SetDuration(INTERVAL_PERIOD).
		SetParams("aaa", 1, time.Now()).
		SetFunc(show).
		Run()
}

func show(s string, i int, t time.Time) {
	fmt.Println(s)
	fmt.Println(i)
	fmt.Println(t)
}
