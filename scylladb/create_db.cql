CREATE ROLE animalsapp WITH LOGIN = true and PASSWORD = '1f2d1e2e67df';



CREATE KEYSPACE IF NOT EXISTS uoc_animals WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};

-- CREATE KEYSPACE IF NOT EXISTS media_player WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': '3'}  AND durable_writes = true AND TABLETS = {'enabled': false};
GRANT MODIFY ON KEYSPACE uoc_animals TO animalsapp;

CREATE TABLE IF NOT EXISTS uoc_animals.media (id uuid,name text,contenttype text,location text,created_at timestamp,status int,size int, PRIMARY KEY (id));

CREATE TABLE IF NOT EXISTS uoc_animals.agent (id uuid,lastname text,firstname text,type tinyint,created_at timestamp,status tinyint, PRIMARY KEY (id));
