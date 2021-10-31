package client

import "context"

// vendorClient is an abstraction to group all different vendor clients
// the sum of behaviour must equal the interface defined.
type vendorClient struct {
	*solanartClient
	*magicEdenClient
}

// Ping ...
func (v *vendorClient) Ping(ctx context.Context) error {
	return nil
}
