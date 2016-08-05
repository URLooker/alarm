package cache

import (
	"sync"

	"github.com/urlooker/web/model"
)

type SafeStrategyMap struct {
	sync.RWMutex
	M map[int64]model.Strategy
}

var StrategyMap = &SafeStrategyMap{M: make(map[int64]model.Strategy)}

func (this *SafeStrategyMap) ReInit(m map[int64]model.Strategy) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyMap) Get(key int64) (model.Strategy, bool) {
	this.RLock()
	defer this.RUnlock()
	s, exists := this.M[key]
	return s, exists
}
