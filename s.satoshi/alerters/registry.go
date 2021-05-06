package alerters

import "sync"

var (
	SatohsiAlerterRegistry []Alerter
	registryMtx            sync.RWMutex
)

func register(alerter Alerter) {
	registryMtx.Lock()
	defer registryMtx.Unlock()
	SatohsiAlerterRegistry = append(SatohsiAlerterRegistry, alerter)
}
