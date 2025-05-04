export KUBECONFIG=/home/jgnoguer/uocWksp/kubectl/uoc-cubie.yaml

SCYLLADB_CONFIG="$( mktemp -d )" 

cat <<EOF > "${SCYLLADB_CONFIG}/credentials"
[PlainTextAuthProvider]
username = cassandra
password = cassandra
EOF
chmod 600 "${SCYLLADB_CONFIG}/credentials"

#SCYLLADB_DISCOVERY_EP="$( kubectl -n scylla get service/scylladb-uoc-client -o='jsonpath={.spec.clusterIP}' )"
SCYLLADB_DISCOVERY_EP="$( kubectl -n scylla get service/scylladb-uoc-client -o='jsonpath={.status.loadBalancer.ingress[0].ip}' )"
kubectl -n scylla get configmap/scylladb-uoc-local-serving-ca -o='jsonpath={.data.ca-bundle\.crt}' > "${SCYLLADB_CONFIG}/serving-ca-bundle.crt"
kubectl -n scylla get secret/scylladb-uoc-local-user-admin -o='jsonpath={.data.tls\.crt}' | base64 -d > "${SCYLLADB_CONFIG}/admin.crt"
kubectl -n scylla get secret/scylladb-uoc-local-user-admin -o='jsonpath={.data.tls\.key}' | base64 -d > "${SCYLLADB_CONFIG}/admin.key"

cat <<EOF > "${SCYLLADB_CONFIG}/cqlshrc"
[authentication]
credentials = ${SCYLLADB_CONFIG}/credentials
[connection]
hostname = ${SCYLLADB_DISCOVERY_EP}
port = 9142
ssl=true
factory = cqlshlib.ssl.ssl_transport_factory
[ssl]
validate=true
certfile=${SCYLLADB_CONFIG}/serving-ca-bundle.crt
usercert=${SCYLLADB_CONFIG}/admin.crt
userkey=${SCYLLADB_CONFIG}/admin.key
EOF

docker run -it --rm --entrypoint=cqlsh \
-v="${SCYLLADB_CONFIG}:${SCYLLADB_CONFIG}:ro,Z" \
-v="${SCYLLADB_CONFIG}/cqlshrc:/root/.cassandra/cqlshrc:ro,Z" \
docker.io/scylladb/scylla:5.4.3

echo Fin