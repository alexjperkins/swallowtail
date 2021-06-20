package pager

import (
	"fmt"
	"sync"

	"github.com/monzo/terrors"

	accountproto "swallowtail/s.account/proto"
)

var (
	pagers  = map[accountproto.PagerType]Pager{}
	pagerMu sync.RWMutex
)

func register(id accountproto.PagerType, pager Pager) {
	pagerMu.Lock()
	defer pagerMu.Unlock()

	if _, ok := pagers[id]; ok {
		panic(fmt.Sprintf("Cannot register the same pager twice: %v", id))
	}

	pagers[id] = pager
}

// GetPagerByID is a concurrency safe way to retrieve the given pager by ID
func GetPagerByID(id accountproto.PagerType) (Pager, error) {
	pagerMu.RLock()
	defer pagerMu.RUnlock()
	pager, ok := pagers[id]
	if !ok {
		return nil, terrors.BadRequest("invalid-pager-id", "Failed to retreive pager; doesn't exist", map[string]string{
			"pager_id": string(id),
		})
	}
	return pager, nil
}
