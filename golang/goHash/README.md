## goHash

Golang Hash and Map

### Installation/update
```
go get -u github.com/ne7ermore/goHash
```

### Use

```
package main

import (
    "os"
    "path"
    . "github.com/ne7ermore/goHash"
)

func main() {
    fp, err := os.Getwd()
    if err != nil {
        t.Error(err)
    }
    testFile := path.Join(fp, "data", "test.txt")
    mapW := GetMapword()
    if err := mapW.LoadMapwords(testFile); err != nil {
        t.Error(err)
    }
    // map
    m := NewMap()
    vmap := map[string]string{
        "1": "one",
        "2": "two",
    }
    m.Add("test", vmap)
    // hash
    hash := New()
    hash.Add("test")
}
```

