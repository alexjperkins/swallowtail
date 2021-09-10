package handler

import satoshiproto "swallowtail/s.satoshi/proto"

type SatoshiService struct {
	*satoshiproto.UnimplementedSatoshiServer
}
