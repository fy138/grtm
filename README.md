# about grtm 
[![Build Status](https://travis-ci.org/fy138/grtm.svg?branch=master)](https://travis-ci.org/fy138/grtm)

grtm is a tool to manage golang goroutines.use this can start or stop a long loop goroutine.
* by fy138
* 增加了查询任务的数量
* 增加了当前任务列表
* 去掉了 nosignal 提示，方便使用在产品中
* 增加线程池功能
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
output

```bash
hello+1
hello+2
hello+1
hello+2
hello+1
hello+2
task name:myfunc,task id:5577006791947779410,task name2:myfunc
task name:myfunc2,task id:8674665223082153551,task name2:myfunc2
NumTask:2
hello+1
hello+2
hello+1
hello+2
```
增加线程池功能，限制任务的线程数量
```golang
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fy138/grtm"
)

func main() {
	go func() {
		for {
			//get  goroutine total
			fmt.Println("go goroutines:", runtime.NumGoroutine())
			time.Sleep(time.Second * 1)
		}

	}()
	//建立线程池
	pool_1 := grtm.NewPool(3)
	pool_2 := grtm.NewPool(2)
	for i := 100; i >= 1; i-- {
		fmt.Println("I=", i)
		//通过通道来限制goroutine 数量
		/* 下面是第一种调用方法 */
		pool_1.LimitChan <- true //importan
		pool_1.AddTask(Download, i, "Download_1")

		/* 下面是第二种调用方法 */
		pool_2.LimitChan <- true //importan
		go func(i int, str string) {
			Download2(i, str)
			//函数执行完释放通道
			defer func() {
				<-pool_2.LimitChan
			}()
		}(i, "Download_2")

	}
	time.Sleep(time.Second * 20) //防止主线程提前退出
}

func Download(args ...interface{}) {
	time.Sleep(2 * time.Second)
	fmt.Printf("%s => %d \n", args[0].([]interface{})[1].(string), args[0].([]interface{})[0].(int))
}
func Download2(i int, str string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("%s => %d \n", str, i)
}


```
```bash
I= 100
go goroutines: 4
I= 99
I= 98
go goroutines: 9
Download_2 => 100
I= 97
Download_1 => 100
go goroutines: 9
Download_2 => 99
I= 96
Download_1 => 99
Download_1 => 98
go goroutines: 8
Download_2 => 98
I= 95
Download_1 => 97
go goroutines: 8
Download_2 => 97
I= 94
Download_1 => 96
go goroutines: 8
Download_2 => 96
I= 93
Download_1 => 95
go goroutines: 8
Download_2 => 95
I= 92
Download_1 => 94
go goroutines: 8
Download_2 => 94
I= 91

```