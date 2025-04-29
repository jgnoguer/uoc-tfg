export KUBECONFIG=/home/jgnoguer/uocWksp/kubectl/uoc-cubie.yaml

kubectl -n scylla patch service/scylladb-uoc-client -p '{"metadata": {"annotations": {"networking.gke.io/load-balancer-type": "Internal"}}, "spec": {"type": "LoadBalancer"}}'
kubectl -n scylla wait --for=jsonpath='{.status.loadBalancer.ingress}' service/scylladb-uoc-client
kubectl -n scylla get service/scylladb-uoc-client -o='jsonpath={.status.loadBalancer.ingress[0].ip}'