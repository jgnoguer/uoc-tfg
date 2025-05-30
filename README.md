# uoc-tfg
TFG - Aplicacions i sistemes distribuïts



# Cubie images

https://github.com/Misaka-Nnnnq/Radxa_A5E_Firmware

# KNative func

https://github.com/knative/func/releases

## Registry

## ScyllaDB

https://operator.docs.scylladb.com/stable/index.html

NS=`kubectl get ns |grep Terminating | awk 'NR==1 {print $1}'` && kubectl get namespace "$NS" -o json   | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/"   | kubectl replace --raw /api/v1/namespaces/$NS/finalize -f -


## Longhorn

https://longhorn.io/kb/tip-only-use-storage-on-a-set-of-nodes/


192.168.2.123
192.168.2.126
192.168.2.137
192.168.2.222 

## curl

kubectl run mycurlpod --image=curlimages/curl -i --tty -- sh

