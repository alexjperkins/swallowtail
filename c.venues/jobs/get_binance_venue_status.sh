#/bin/sh -x

# Runs get Binance venue status.
echo Calling s.binance via gRPC. Gets venue status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-binance:8000 \
	binance.GetStatus
