#/bin/sh -x

# Runs get FTX venue status.
echo Calling s.ftx via gRPC. Gets venue status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-ftx:8000 \
	ftx.GetFTXStatus
