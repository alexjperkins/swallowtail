#/bin/sh -x

# Publishes latest price information.
echo Calling s.market-data via gRPC. Publishes latest price information.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishLatestPriceInformation
