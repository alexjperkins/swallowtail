#/bin/sh -x

# Runs get Binance exchange status.
echo Calling s.binance via gRPC. Gets Exchange Status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-binance:8000 \
	binance.GetStatus
