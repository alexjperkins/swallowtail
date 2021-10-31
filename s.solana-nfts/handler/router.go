package handler

import (
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

// SolanaNFTsService ...
type SolanaNFTsService struct {
	*solananftsproto.UnimplementedSolananftsServer
}
