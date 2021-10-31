#/bin/sh -x

# Publishes solana nft price information for all assets.
echo Calling s.market-data via gRPC. Publishes solana nft price information.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-marketdata:8000 \
	marketdata.PublishSolanaNFTPriceInformation 
