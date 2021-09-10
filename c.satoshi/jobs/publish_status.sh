#/bin/sh -x

# Runs get Satoshi status.
echo Calling s.satoshi via gRPC. Gets Satoshi Status.

exec grpcurl -plaintext -d \
	'{}' \
	swallowtail-s-satoshi:8000 \
	satoshi.PublishStatus
