# Label nodes

https://operator.docs.scylladb.com/stable/installation/kubernetes/generic.html


kubectl label nodes uoc-rock3a-01 scylla.scylladb.com/node-type=scylla

Client discovery

https://operator.docs.scylladb.com/stable/resources/scyllaclusters/clients/discovery.html


CREATE KEYSPACE IF NOT EXISTS media_player WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': '3'}  AND durable_writes = true AND TABLETS = {'enabled': false};
CREATE TABLE IF NOT EXISTS media_player.playlist (id uuid,title text,album text,artist text,created_at timestamp,PRIMARY KEY (id, created_at)) WITH CLUSTERING ORDER BY (created_at DESC);