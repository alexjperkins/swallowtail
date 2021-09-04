#/bin/sh -x

# Runs enforce subscriptions job
echo Calling s.payments via gRPC. Enforces subscriptions.

exec grpcurl -plaintext -d \
	'{ "actor_id": "cron:enforce-subscriptions" }' \
	swallowtail-s-payments:8000 \
	payments.EnforceSubscriptions
