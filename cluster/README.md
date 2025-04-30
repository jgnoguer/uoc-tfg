# Taints

kubectl taint nodes uoc-neo2core-01 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-01 memorytype=low:NoExecute
kubectl taint nodes uoc-neo2core-02 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-02 memorytype=low:NoExecute
kubectl taint nodes uoc-neo2core-03 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-03 memorytype=low:NoExecute

kubectl taint nodes uoc-neo2core-01 istiogateway=compatible:NoSchedule
kubectl taint nodes uoc-neo2core-02 istiogateway=compatible:NoSchedule
kubectl taint nodes uoc-neo2core-03 istiogateway=compatible:NoSchedule



# Labels

scylla.scylladb.com/node-type=scylla

kubectl label nodes uoc-neo2core-01 istiogateway=compatible
kubectl label nodes uoc-neo2core-02 istiogateway=compatible
kubectl label nodes uoc-neo2core-03 istiogateway=compatible

kubectl label nodes uoc-cubie traefikports=compatible
kubectl label nodes uoc-r2splus-01 traefikports=compatible
kubectl label nodes uoc-r2splus-02 traefikports=compatible
