# about grtm 
[![Build Status](https://travis-ci.org/fy138/grtm.svg?branch=master)](https://travis-ci.org/fy138/grtm)

grtm is a tool to manage golang goroutines.use this can start or stop a long loop goroutine.
*by fy138
*增加了查询任务的数量
*增加了当前任务列表
*去掉了 nosignal 提示，方便使用在产品中
## Getting started
```bash
go get github.com/fy138/grtm
```

## Create normal goroutine

```golang
package main

import (
        "fmt"
        "github.com/fy138/grtm"
        "time"
       )

func normal() {
    fmt.Println("i am normal goroutine")
}

func main() {
        gm := grtm.NewGrManager()
        gm.NewGoroutine("normal", normal)
        fmt.Println("main function")
        time.Sleep(time.Second * time.Duration(5))
}
~
```

## Create normal goroutine function with params

```golang
package main

import (
        "fmt"
        "github.com/fy138/grtm"
        "time"
       )

func normal() {
    fmt.Println("i am normal goroutine")
}

func funcWithParams(args ...interface{}) {
    fmt.Println(args[0].([]interface{})[0].(string))
    fmt.Println(args[0].([]interface{})[1].(string))
}

func main() {
        gm := grtm.NewGrManager()
        gm.NewGoroutine("normal", normal)
        fmt.Println("main function")
        gm.NewGoroutine("funcWithParams", funcWithParams, "hello", "world")
        time.Sleep(time.Second * time.Duration(5))
}
```

## Create long loop goroutine then stop it

```golang
package main

import (
        "fmt"
        "github.com/fy138/grtm"
        "time"
       )

func myfunc() {
    fmt.Println("do something repeat by interval 4 seconds")
        time.Sleep(time.Second * time.Duration(4))
}

func main() {
gm := grtm.NewGrManager()
        gm.NewLoopGoroutine("myfunc", myfunc)
        fmt.Println("main function")
        time.Sleep(time.Second * time.Duration(40))
        fmt.Println("stop myfunc goroutine")
        gm.StopLoopGoroutine("myfunc")
        time.Sleep(time.Second * time.Duration(80))
}
```

output

```bash
main function
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
stop myfunc goroutine
gid[5577006791947779410] quit

```
*fy138
```golang
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fy138/grtm"
)

func myfunc(me interface{}) {
	fmt.Println("hello+" + me.(string))
	time.Sleep(time.Second * 2)
}
func main() {
	gm := grtm.NewGrManager()

	gm.NewLoopGoroutine("myfunc", myfunc, "1")
	gm.NewLoopGoroutine("myfunc2", myfunc, "2")
	fmt.Println("main function")
	fmt.Printf("NumGoroutine:%d\n", runtime.NumGoroutine())

	for {
		for k, v := range gm.GetAllTask() {
			fmt.Printf("task name:%s,task id:%d,task name2:%s\n", k, v.Gid, v.Name)
		}
		fmt.Printf("NumTask:%d\n", gm.GetTaskTotal())
		time.Sleep(time.Second * 5)
	}
}


```
