package accountproto

const (
	// Valid actor IDs.
	ActorSystemPayments    = "actor-system-payments"
	ActorSystemTradeEngine = "actor-system-tradeengine"
	ActorManual            = "manual"
)

const (
	// Request context.
	RequestContextOrderRequest = "order-request" // request from a user indirectly via an order.
	RequestContextUserRequest  = "user-request"  // direct request from a user.
)

const (
	// SubAccountUnknown defines the constant used for exchanges that don't support subaccounts.
	SubAccountUnknown = "UNKNOWN"
)
