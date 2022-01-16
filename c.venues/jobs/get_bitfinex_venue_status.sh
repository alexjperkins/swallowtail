#/bin/sh -x

# Runs get Bitfinex venue status.
echo Calling s.bitfinex via gRPC. Gets venue status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-bitfinex:8000 \
	bitfinex.GetBitfinexStatus
