#/bin/sh -x

# Runs get Binance exchange status.
echo Calling s.binance via gRPC. Gets Exchange Status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishLatestPriceInformation
