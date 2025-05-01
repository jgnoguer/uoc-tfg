-- Active: 1746104101425@@192.168.2.126@9042
CREATE KEYSPACE IF NOT EXISTS uoc_animals WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};
CREATE KEYSPACE IF NOT EXISTS media_player WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': '3'}  AND durable_writes = true;