# Cron: c.satoshi

This calls given endpoints via gRPC in order to publish heartbeats to the channel `#satoshi-pulse`.

The crontab runs every 5 minutes.

It publishes certain statistics about the exchange such as:

- Aliveness
- Consumer health
