package domain

import (
	"encoding/json"
	"swallowtail/libraries/transport"
)

// WsMsgToBinanceEvent converts a websocket message to a BinanceEvent based on the constructor
// provided
func WsMsgToBinanceEvent(msg *transport.WsMessage, constructor func() interface{}) (interface{}, error) {
	bEvent := constructor()
	if err := json.Unmarshal(msg.Raw, bEvent); err != nil {
		return nil, err
	}
	return &bEvent, nil
}
