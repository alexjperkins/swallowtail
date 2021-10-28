package dto

// GetStatusRequest ...
type GetStatusRequest struct {
}

type GetStatusProxyResponse [1]int

func (r *GetStatusProxyResponse) Operative() int {
	return r[0]
}

// GetStatusResponse ...
type GetStatusResponse struct {
	// Defines if the Bitfinex platfrom is live.
	Operative int `json:"operative"`
	// Latency of the server.
	ServerLatency int `json:"-"`
}
