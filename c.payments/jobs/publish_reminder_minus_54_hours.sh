#/bin/sh -x

# Runs enforce subscriptions job
echo Calling s.payments via gRPC. Publishes Reminder minus 40H.

exec grpcurl -plaintext -d \
	'{ "actor_id": "cron:publish-reminder", "reminder_type": "MINUS_54_HOURS" }' \
	swallowtail-s-payments:8000 \
	payments.PublishSubscriptionReminder
