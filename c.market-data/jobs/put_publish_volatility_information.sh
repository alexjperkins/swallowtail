#/bin/sh -x

# Publishes volatility information for all assets.
echo Calling s.market-data via gRPC. Publishes volatility information.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishVolatilityInformation
