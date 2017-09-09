package grtm

import (
	"fmt"
	"math/rand"
	"sync"
)

const (
	STOP = "__P:"
)

type GoroutineChannel struct {
	Gid  uint64
	Name string
	Msg  chan string
}

type GoroutineChannelMap struct {
	mutex      sync.Mutex
	Grchannels map[string]*GoroutineChannel
}

func (m *GoroutineChannelMap) unregister(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.Grchannels[name]; !ok {
		return fmt.Errorf("goroutine channel not find: %q", name)
	}
	delete(m.Grchannels, name)
	return nil
}

func (m *GoroutineChannelMap) register(name string) error {
	gchannel := &GoroutineChannel{
		Gid:  uint64(rand.Int63()),
		Name: name,
	}
	gchannel.Msg = make(chan string)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.Grchannels == nil {
		m.Grchannels = make(map[string]*GoroutineChannel)
	} else if _, ok := m.Grchannels[gchannel.Name]; ok {
		//fmt.Printf("goroutine channel already defined: %q\n", gchannel.Name)
		return fmt.Errorf("goroutine channel already defined: %q", gchannel.Name)
	}
	m.Grchannels[gchannel.Name] = gchannel
	return nil
}
