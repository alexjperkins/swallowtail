package arbitrage

import "sync"

var (
	arbitrageClients = map[string]ArbitrageClient{}
	mtx              sync.Mutex
)

func register(clientID string, client ArbitrageClient) {
	mtx.Lock()
	defer mtx.Unlock()
}

func getAllArbitrageClients() []ArbitrageClient {
	mtx.Lock()
	defer mtx.Unlock()
	cs := []ArbitrageClient{}
	for _, c := range arbitrageClients {
		cs = append(cs, c)
	}
	return cs
}
