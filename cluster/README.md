# Setting up a K3s HA Cluster

## k3s install

See https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/kubernetes-cluster-setup/k3s-for-rancher
See https://docs.k3s.io/datastore/ha-embedded


### Server

#### First master server

curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --disable traefik --snapshotter=native --cluster-init --tls-san 192.168.2.1" K3S_TOKEN="THETOKEN" sh -

Get main server node token:

cat /var/lib/rancher/k3s/server/token

#### Second and third server

Install the second server node:

curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --disable traefik --snapshotter=native --tls-san 192.168.2.1 --server https://uoc-zero2-01:6443" K3S_TOKEN="THETOKEN" sh -

curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --disable traefik --snapshotter=native --server https://uoc-zero2-01:6443" K3S_TOKEN="THETOKEN" sh -

K3S_TOKEN=
check /var/lib/rancher/k3s/server/node-token

Kubeconfig

 cat /etc/rancher/k3s/k3s.yaml


### Agents

curl -sfL https://get.k3s.io | K3S_URL=https://192.168.2.1:6443 K3S_TOKEN="THETOKEN" sh -

   35  ps -aux
   36  reboot
   37  k3s kubectl get nodes
   38  k3s agent
   39  k3s agent --server rock-3a
   40  k3s agent --server rock-3a --token theks3token
   41  cat /var/log/k3s.log

   43  service k3s-agent sta
   45  cat /var/lib/rancher/k3s/agent/containerd/containerd.log
   46  service k3s-agent logs
   47  journalctl -u k3s-agent
   48  journalctl -u k3s-agent > logs.txt
   49  more logs.txt 

   52  df -h

https://docs.k3s.io/quick-start

#### RaspberryPi cm4

https://learn.umh.app/course/resolving-cgroup-v2-memory-issues-when-running-umh-lite-in-docker-on-raspberry-pi/

#### Nanopi core / nanopi zero2

curl -sfL https://get.k3s.io | K3S_URL=https://192.168.2.1:6443 K3S_TOKEN="thetoken" INSTALL_K3S_EXEC="agent --snapshotter=native" sh -


jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo cp /etc/rancher/k3s/k3s.yaml .
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo chown jgnoguer:jgnoguer k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ mv k3s.yaml ~
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ export KUBECONFIG=/home/jgnoguer/k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ kubectl get nodes



# Taints

kubectl taint nodes uoc-neo2core-01 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-02 memorytype=low:NoSchedule
kubectl taint nodes uoc-neo2core-03 memorytype=low:NoSchedule


kubectl taint nodes uoc-orangezeroplus2-01 memorytype=low:NoSchedule

kubectl taint nodes uoc-zero2-01 uocnodetype=master:NoSchedule
kubectl taint nodes uoc-zero2-02 uocnodetype=master:NoSchedule
kubectl taint nodes uoc-zero2-03 uocnodetype=master:NoSchedule
kubectl taint nodes uoc-zero2-01 uocnodetype=master:NoExecute
kubectl taint nodes uoc-zero2-02 uocnodetype=master:NoExecute
kubectl taint nodes uoc-zero2-03 uocnodetype=master:NoExecute

kubectl taint nodes uoc-rock3a-01 scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule
kubectl taint nodes uoc-rock3a-02 scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule
kubectl taint nodes uoc-rock3a-03 scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule


# Labels

kubectl label nodes uoc-rock3a-01 scylla.scylladb.com/node-type=scylla
kubectl label nodes uoc-rock3a-02 scylla.scylladb.com/node-type=scylla
kubectl label nodes uoc-rock3a-03 scylla.scylladb.com/node-type=scylla


kubectl label nodes uoc-neo2core-01 envoyLib=compatible
kubectl label nodes uoc-neo2core-02 envoyLib=compatible
kubectl label nodes uoc-neo2core-03 envoyLib=compatible

kubectl label nodes uoc-cubie-02 mosquitto=compatible


# Loadbalancer HA

https://www.google.com/search?client=firefox-b-lm&channel=entpr&q=k3s+load+balancer+external+ip


2025-05-18T10:15:28.886205Z	info	Envoy command: [-c etc/istio/proxy/envoy-rev.json --drain-time-s 45 --drain-strategy immediate --local-address-ip-version v4 --file-flush-interval-msec 1000 --disable-hot-restart --allow-unknown-static-fields -l warning --component-log-level misc:error --skip-deprecated-logs --concurrency 2]
2025-05-18T10:15:28.892844Z	info	sds	Starting SDS grpc server
2025-05-18T10:15:28.893430Z	info	sds	Starting SDS server for workload certificates, will listen on "var/run/secrets/workload-spiffe-uds/socket"
14 external/com_github_google_tcmalloc/tcmalloc/system-alloc.cc:769] MmapAligned() failed - unable to allocate with tag (hint=0xd2b40000000, size=1073741824, alignment=1073741824) - is something limiting address placement?
14 external/com_github_google_tcmalloc/tcmalloc/system-alloc.cc:776] Note: the allocation may have failed because TCMalloc assumes a 48-bit virtual address space size; you may need to rebuild TCMalloc with TCMALLOC_ADDRESS_BITS defined to your system's virtual address space size
14 external/com_github_google_tcmalloc/tcmalloc/arena.cc:56] CHECK in Alloc: FATAL ERROR: Out of memory trying to allocate internal tcmalloc data (bytes=131072, object-size=16384); is something preventing mmap from succeeding (sandbox, VSS limitations)?
2025-05-18T10:15:28.903926Z	error	Envoy exited with error: signal: aborted