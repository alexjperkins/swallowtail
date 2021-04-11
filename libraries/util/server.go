package util

import "fmt"

func ServiceRPCAddress(hostname, port string) (string, error) {
	return fmt.Sprintf("%s:%s", hostname, port), nil
}
