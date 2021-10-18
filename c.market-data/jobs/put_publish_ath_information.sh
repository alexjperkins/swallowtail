#/bin/sh -x

# Publishes ATH information
echo Calling s.market-data via gRPC. Publishes ATH information.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishATHInformation
