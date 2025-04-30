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

kubectl label nodes uoc-neo2core-01 envoyLib=compatible
kubectl label nodes uoc-neo2core-02 envoyLib=compatible
kubectl label nodes uoc-neo2core-03 envoyLib=compatible

kubectl label nodes uoc-zero2-01 istiogatewaylb=compatible
kubectl label nodes uoc-zero2-02 istiogatewaylb=compatible
kubectl label nodes uoc-zero2-03 istiogatewaylb=compatible

