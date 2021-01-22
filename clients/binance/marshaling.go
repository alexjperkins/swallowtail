package binance

import (
	"encoding/json"
	"swallowtail/clients/binance/domain"
	"swallowtail/transport"
)

func wsMsgToBinanceMsg(msg *transport.WsMessage) (*domain.BinanceMsg, error) {
	bMsg := &domain.BinanceMsg{}
	if err := json.Unmarshal(msg.Raw, bMsg); err != nil {
		return nil, err
	}
	return &domain.BinanceMsg{}, nil
}
