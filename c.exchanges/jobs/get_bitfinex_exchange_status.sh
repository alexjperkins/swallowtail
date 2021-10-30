#/bin/sh -x

# Runs get Bitfinex exchange status.
echo Calling s.bitfinex via gRPC. Gets Exchange Status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-bitfinex:8000 \
	bitfinex.GetBitfinexStatus
