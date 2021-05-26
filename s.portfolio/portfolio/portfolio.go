package portfolio

import (
	"context"
	"swallowtail/s.portfolio/owner"
	"swallowtail/s.portfolio/sync"
)

type Portfolio interface {
	Sync(context.Context)
	Done()
	OwnerID() string
}

func New(owner *owner.Owner, syncer sync.Syncer) Portfolio {
	return &portfolio{
		Owner:  owner,
		Syncer: syncer,
	}
}

type portfolio struct {
	Owner  *owner.Owner
	Syncer sync.Syncer
}

func (p *portfolio) Sync(context.Context) {
}

func (p *portfolio) Done() {
}

func (p *portfolio) OwnerID() string {
	return p.Owner.Name
}
