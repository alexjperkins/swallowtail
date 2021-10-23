#/bin/sh -x

# Publishes funding rates information
echo Calling s.market-data via gRPC. Publishes funding rates information.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishFundingRatesInformation
