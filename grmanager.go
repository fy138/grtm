package grtm

import (
	"fmt"
	"strconv"
	"strings"
)

type GrManager struct {
	grchannelMap *GoroutineChannelMap
	ErrChan      chan error
	NotiChan     chan string
}

func NewGrManager() *GrManager {
	gm := &GoroutineChannelMap{}
	errchan := make(chan error)
	notichan := make(chan string)
	return &GrManager{grchannelMap: gm, ErrChan: errchan, NotiChan: notichan}
}

func (gm *GrManager) StopLoopGoroutine(name string) {
	stopChannel, ok := gm.grchannelMap.Grchannels[name]
	if !ok {
		gm.ErrChan <- fmt.Errorf("not found goroutine name :" + name)
		return
	}
	gm.grchannelMap.Grchannels[name].Msg <- STOP + strconv.Itoa(int(stopChannel.Gid))
	//return nil
}

func (gm *GrManager) NewLoopGoroutine(name string, fc interface{}, args ...interface{}) {
	go func(this *GrManager, n string, fc interface{}, args ...interface{}) {
		//register channel
		err := this.grchannelMap.register(n)
		if err != nil {
			gm.ErrChan <- err
			return
		}
		for {
			select {
			case info := <-this.grchannelMap.Grchannels[name].Msg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(this.grchannelMap.Grchannels[name].Gid)) {
					if signal == "__P" {
						gm.NotiChan <- fmt.Sprintf("gid[%s]quit", gid)
						err := this.grchannelMap.unregister(name)
						if err != nil {
							gm.ErrChan <- err
						}
						return
					} else {
						gm.ErrChan <- fmt.Errorf("unknown signal")
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
			gm.ErrChan <- err
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		err = gm.grchannelMap.unregister(name)
		if err != nil {
			gm.ErrChan <- err
			return
		}
	}(name, fc, args...)

}
func (gm *GrManager) GetAllTask() map[string]*GoroutineChannel {
	return gm.grchannelMap.Grchannels
}
func (gm *GrManager) GetTaskTotal() int {
	return len(gm.grchannelMap.Grchannels)
}
