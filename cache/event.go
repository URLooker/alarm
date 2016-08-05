package cache

import (
	"sync"

	"github.com/urlooker/web/model"
)

type SafeEventMap struct {
	sync.RWMutex
	M map[string]*model.Event
}

var LastEvents = &SafeEventMap{M: make(map[string]*model.Event)}

func (this *SafeEventMap) Get(key string) (*model.Event, bool) {
	this.RLock()
	defer this.RUnlock()
	event, exists := this.M[key]
	return event, exists
}

func (this *SafeEventMap) Set(key string, event *model.Event) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = event
}

func (this *SafeEventMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}
