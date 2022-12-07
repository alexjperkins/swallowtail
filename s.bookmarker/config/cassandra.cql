CREATE KEYSPACE IF NOT EXISTS bookmarker WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'birdperch-1': '3'
} AND durable_writes = true;
