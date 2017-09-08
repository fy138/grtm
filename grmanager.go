package grtm

import (
	"fmt"
	"strconv"
	"strings"
)

type GrManager struct {
	grchannelMap *GoroutineChannelMap
}

func NewGrManager() *GrManager {
	gm := &GoroutineChannelMap{}
	return &GrManager{grchannelMap: gm}
}

func (gm *GrManager) StopLoopGoroutine(name string) error {
	stopChannel, ok := gm.grchannelMap.Grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	gm.grchannelMap.Grchannels[name].Msg <- STOP + strconv.Itoa(int(stopChannel.Gid))
	return nil
}

func (gm *GrManager) NewLoopGoroutine(name string, fc interface{}, args ...interface{}) {
	go func(this *GrManager, n string, fc interface{}, args ...interface{}) {
		//register channel
		err := this.grchannelMap.register(n)
		if err != nil {
			return
		}
		for {
			select {
			case info := <-this.grchannelMap.Grchannels[name].Msg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(this.grchannelMap.Grchannels[name].Gid)) {
					if signal == "__P" {
						fmt.Println("gid[" + gid + "] quit")
						this.grchannelMap.unregister(name)
						return
					} else {
						fmt.Println("unknown signal")
					}
				}
			default:
				//fmt.Println("no signal")
			}

			if len(args) > 1 {
				fc.(func(...interface{}))(args)
			} else if len(args) == 1 {
				fc.(func(interface{}))(args[0])
			} else {
				fc.(func())()
			}
		}
	}(gm, name, fc, args...)
}

func (gm *GrManager) NewGoroutine(name string, fc interface{}, args ...interface{}) {
	go func(n string, fc interface{}, args ...interface{}) {
		//register channel
		err := gm.grchannelMap.register(n)
		if err != nil {
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		gm.grchannelMap.unregister(name)
	}(name, fc, args...)

}
func (gm *GrManager) GetAllTask() map[string]*GoroutineChannel {
	return gm.grchannelMap.Grchannels
}
func (gm *GrManager) GetTaskTotal() int {
	return len(gm.grchannelMap.Grchannels)
}
