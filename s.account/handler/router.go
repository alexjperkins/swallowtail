package handler

import (
	accountproto "swallowtail/s.account/proto"
)

// AccountService ...
type AccountService struct {
	accountproto.UnimplementedAccountServer
}
