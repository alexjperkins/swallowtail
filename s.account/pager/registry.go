package pager

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/monzo/terrors"
)

var (
	pagers  = map[int]Pager{}
	pagerMu sync.RWMutex
)

func register(id int, pager Pager) {
	pagerMu.Lock()
	defer pagerMu.Unlock()

	if _, ok := pagers[id]; ok {
		panic(fmt.Sprintf("Cannot register the same pager twice: %v", id))
	}

	pagers[id] = pager
}

// GetPagerByID is a concurrency safe way to retrieve the given pager by ID
func GetPagerByID(id int) (Pager, error) {
	pagerMu.RLock()
	defer pagerMu.RUnlock()
	pager, ok := pagers[id]
	if !ok {
		return nil, terrors.BadRequest("invalid-pager-id", "Failed to retreive pager; doesn't exist", map[string]string{
			"pager_id": strconv.Itoa(id),
		})
	}
	return pager, nil
}
