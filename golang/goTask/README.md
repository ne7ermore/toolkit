## goTask: A Golang Task Scheduling Package.

Golang Task scheduling package, simplely and human-friendly

### Installation/update
```
go get -u github.com/ne7ermore/goTask
```

### Use
```
package main

import (
    "fmt"
    "github.com/ne7ermore/goTask"
    "time"
)

const (
    DURATION        time.Duration = 60 * time.Minute
    HOUR_TO_TASK    int           = 1
    MINUTE_TO_TASK  int           = 0
    SECOND_TO_TASK  int           = 0
)

func main() {
    tk := goTask.NewTask()
    tk.SetTaskTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TASK, MINUTE_TO_TASK, SECOND_TO_TASK, 0, time.Local)).
        SetDuration(DURATION).
        SetParams("tj", 1987, time.Now()).
        SetFunc(show).
        Run()
}

func show(s string, i int, t time.Time) {
    fmt.Println(s, i, t)
}
```
