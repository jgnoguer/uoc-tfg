# Label nodes

https://operator.docs.scylladb.com/stable/installation/kubernetes/generic.html



kubectl label nodes uoc-rock3a-03 scylla.scylladb.com/node-type=scylla

kubectl label nodes uoc-rpicm4-02 scylla.scylladb.com/node-type=scylla

kubectl taint nodes uoc-rock3a-03 scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule
kubectl taint nodes uoc-rpicm4-02 scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule



Client discovery

https://operator.docs.scylladb.com/stable/resources/scyllaclusters/clients/discovery.html



# Local development

https://opensource.docs.scylladb.com/stable/operating-scylla/procedures/tips/best-practices-scylla-on-docker.html

docker run --name uoc-localscylla -d scylladb/scylla
docker run --name uoc-localscylla --volume ~/uocWksp/repo/scylladb/local/master_scylla.yaml:/etc/scylla/scylla.yaml -d docker.io/scylladb/scylla

-- docker run --name uoc-localscylla-node2 -d scylladb/scylla --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' uoc-localscylla)"

docker inspect --format='{{ .NetworkSettings.IPAddress }}' uoc-localscylla
docker exec -it uoc-localscylla nodetool status

docker exec -it uoc-localscylla cqlsh

https://opensource.docs.scylladb.com/stable/getting-started/install-scylla/run-in-docker.html



