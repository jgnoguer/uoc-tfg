# ScyllaDB config

Connect to scylladb pod (open shell)

cqlsh -u cassandra
(password cassandra)


# Local development

https://opensource.docs.scylladb.com/stable/operating-scylla/procedures/tips/best-practices-scylla-on-docker.html

docker run --name uoc-localscylla -d scylladb/scylla
docker run --name uoc-localscylla -p 9042:9042 --volume ~/uocWksp/repo/scylladb/local/master_scylla.yaml:/etc/scylla/scylla.yaml -d docker.io/scylladb/scylla

-- docker run --name uoc-localscylla-node2 -d scylladb/scylla --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' uoc-localscylla)"

docker inspect --format='{{ .NetworkSettings.IPAddress }}' uoc-localscylla
docker exec -it uoc-localscylla nodetool status

docker exec -it uoc-localscylla cqlsh

https://opensource.docs.scylladb.com/stable/getting-started/install-scylla/run-in-docker.html



