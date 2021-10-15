#/bin/sh -x

# Runs get FTX exchange status.
echo Calling s.ftx via gRPC. Gets Exchange Status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-ftx:8000 \
	ftx.GetFTXStatus
