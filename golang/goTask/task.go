// Copyright 2017 All rights reserved.
// Author: ne7ermore.
package goTask

import (
	"log"
	"reflect"
	"time"
)

type Task struct {

	// the time task begin running
	StartTime time.Time

	// duration between last run and next run
	Duration time.Duration

	// function task store
	Func interface{}

	// params of task function
	Params []reflect.Value
}

// create a new task
func NewTask() *Task {
	return &Task{}
}

// set the time wich task run with
func (t *Task) SetTaskTime(taskTime time.Time) *Task {
	t.StartTime = taskTime
	return t
}

// set the duration between last and next task run
func (t *Task) SetDuration(d time.Duration) *Task {
	t.Duration = d
	return t
}

// set the function store in task
func (t *Task) SetFunc(f interface{}) *Task {
	t.Func = f
	return t
}

// set the params wich task function with
func (t *Task) SetParams(params ...interface{}) *Task {
	v := make([]reflect.Value, len(params))
	for index, p := range params {
		v[index] = reflect.ValueOf(p)
	}
	t.Params = v
	return t
}

// task run
func (t *Task) Run() {
	f := reflect.ValueOf(t.Func)
	if f.Kind() != reflect.Func {
		log.Fatal("Task.Func type is invalid")
	}
	go func() {
		nextTime, tic := start(t.StartTime, t.Duration)
		for {
			<-tic.C
			f.Call(t.Params)
			nextTime, tic = start(nextTime, t.Duration)
		}
	}()
}

func start(t time.Time, d time.Duration) (time.Time, *time.Ticker) {
	if !t.After(time.Now()) {
		if !t.Add(d).After(time.Now()) {
			t = time.Now().Add(d)
		} else {
			t = t.Add(d)
		}
	}
	diff := t.Sub(time.Now())
	return t, time.NewTicker(diff)
}
