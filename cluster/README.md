# Taints

kubectl taint nodes uoc-neo2core-01 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-01 memorytype=low:NoExecute
kubectl taint nodes uoc-neo2core-02 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-02 memorytype=low:NoExecute
kubectl taint nodes uoc-neo2core-03 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-03 memorytype=low:NoExecute

kubectl taint nodes uoc-neo2core-01 istio-ingress=compatible:NoSchedule
kubectl taint nodes uoc-neo2core-01 istio-ingress=compatible:NoExecute
kubectl taint nodes uoc-neo2core-02 istio-ingress=compatible:NoSchedule
kubectl taint nodes uoc-neo2core-02 istio-ingress=compatible:NoExecute
kubectl taint nodes uoc-neo2core-03 istio-ingress=compatible:NoSchedule
kubectl taint nodes uoc-neo2core-03 istio-ingress=compatible:NoExecute

# Labels

scylla.scylladb.com/node-type=scylla