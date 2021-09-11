# Cron: c.exchanges

This calls given endpoints via gRPC in order to publish heartbeats to the channel `#exchanges-pulse`.

The crontab runs every 5 minutes for:

- Binance
- FTX

It publishes certain statistics about the exchange such as:

- Latency
- Server time
- Assumed clock drift (relative to server time)
