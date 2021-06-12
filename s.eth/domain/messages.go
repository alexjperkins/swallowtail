package domain

type SubscribeToInfuraRequest struct {
	Id      int      `json:"id"`
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type SubscribeToInfuraResponse struct {
	Id      int               `json:"id"`
	JSONRPC string            `json:"jsonrpc"`
	Params  map[string]string `json:"params"`
}
