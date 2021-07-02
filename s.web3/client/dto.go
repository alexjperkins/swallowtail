package client

// SubscribeToInfuraRequest is the message to send to Infura to subscribe to
// their ethereum node.
type SubscribeToInfuraRequest struct {
	Id      int      `json:"id"`
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// SubscribeToInfuraResponse is response object when subscribing to the Infura client.
type SubscribeToInfuraResponse struct {
	Id      int               `json:"id"`
	JSONRPC string            `json:"jsonrpc"`
	Params  map[string]string `json:"params"`
}

type PendingTransactionEvent struct {
	// The block number; if pending this is null.
	Number int `json:number`

	// The hash of the block.
	Hash string `json:"hash"`

	// The hash of the parent block.
	ParentHash string `json:"parent_hash"`

	// Hash of the generated PoW; nil when pending.
	Nonce string `json:"nonce"`

	// SHA3 of the uncles data in the block.
	Sha3Uncles [32]byte `json:"sha3_uncles"`

	// The bloom filter for the logs of the block; nil when pending.
	LogsBloom string `json:"logs_bloom"`

	// The root of the transaction of the block.
	TransactionRoot string `json:"transaction_root"`

	// The root of the final state trie of the block.
	StateRoot string `json:"state_root"`

	// The root of the receipts.
	ReceiptsRoot string `json:"receipts_root"`

	// The address of the beneficiary to whom the mining rewards were given.
	Miner string `json:"miner"`

	// Extra data field.
	ExtraData string `json:"extra_data"`

	// The maximum gas allowed in this block.
	GasLimit int `json:"gas_limit"`

	// The total gas used by all transactions in this block.
	GasUsed int `json:"gas_used"`

	// The unix timestamp for when the the block was collated.
	Timestamp int `json:"timestamp"`
}
