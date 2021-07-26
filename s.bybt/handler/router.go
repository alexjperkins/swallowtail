package handler

import (
	bybtproto "swallowtail/s.bybt/proto"
)

// ByBtService ...
type ByBtService struct {
	*bybtproto.UnimplementedBybtServer
}
