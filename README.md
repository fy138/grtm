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
	pool := grtm.NewPool(10)

	for i := 100; i >= 1; i-- {
		fmt.Println("I=", i)
		//通过通道来限制goroutine 数量，下面这一行不要忘记了
		pool.LimitChan <- true //importan
		pool.AddTask(Download, i, "test", "name")
		/*如果你觉得上面传参数比较麻烦，那么可以把
		pool.AddTask(Download, i, "test")
		替换为
		go func(i int, str string) {
			Download2(i, str)
			defer func() {
				<-pool.LimitChan
			}()
		}(i, "test")
		*/

	}
	time.Sleep(time.Second * 20) //防止主线程提前退出
}

func Download(args ...interface{}) {
	time.Sleep(2 * time.Second)
	fmt.Printf("Download:%d =>%s =>%s \n", args[0].([]interface{})[0].(int), args[0].([]interface{})[1].(string), args[0].([]interface{})[2].(string))
}
func Download2(i int, str string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("Download:%d =>%s \n", i, str)
}


```
```bash
>grtm_test3.exe
I= 100
I= 99
I= 98
I= 97
I= 96
I= 95
I= 94
I= 93
I= 92
I= 91
I= 90
go goroutines: 3
go goroutines: 13
Download:100 =>test =>name
I= 89
Download:98 =>test =>name
I= 88
Download:99 =>test =>name
I= 87
Download:97 =>test =>name
I= 86
Download:93 =>test =>name
I= 85
Download:96 =>test =>name
I= 84
Download:91 =>test =>name
I= 83
Download:95 =>test =>name
I= 82
Download:92 =>test =>name
I= 81
Download:94 =>test =>name
I= 80
go goroutines: 13
go goroutines: 13
Download:90 =>test =>name
I= 79
Download:89 =>test =>name
I= 78
Download:88 =>test =>name
I= 77
Download:87 =>test =>name
I= 76
Download:86 =>test =>name
I= 75
Download:85 =>test =>name
I= 74
Download:84 =>test =>name
Download:82 =>test =>name
Download:83 =>test =>name
Download:81 =>test =>name
go goroutines: 13
I= 73
I= 72
I= 71
I= 70
go goroutines: 13
Download:80 =>test =>name
I= 69
Download:79 =>test =>name
I= 68
Download:78 =>test =>name
I= 67
Download:77 =>test =>name
I= 66
Download:76 =>test =>name
I= 65
Download:75 =>test =>name
I= 64
Download:74 =>test =>name
I= 63
go goroutines: 13
Download:73 =>test =>name
I= 62
Download:72 =>test =>name
I= 61
Download:71 =>test =>name
I= 60
go goroutines: 13
Download:70 =>test =>name
I= 59
Download:69 =>test =>name
I= 58
Download:68 =>test =>name
I= 57
Download:67 =>test =>name
I= 56
Download:66 =>test =>name
I= 55
Download:65 =>test =>name
I= 54
```