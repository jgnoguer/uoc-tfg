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



On nanopi core / nanopi zero2

curl -sfL https://get.k3s.io | K3S_URL=https://192.168.2.1:6443 K3S_TOKEN="THETOKEN" INSTALL_K3S_EXEC="agent --snapshotter=native" sh -


jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo cp /etc/rancher/k3s/k3s.yaml .
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo chown jgnoguer:jgnoguer k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ mv k3s.yaml ~
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ export KUBECONFIG=/home/jgnoguer/k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ kubectl get nodes



# Taints

--kubectl taint nodes uoc-neo2core-01 memorytype=low:NoSchedule
--kubectl taint nodes uoc-neo2core-02 memorytype=low:NoSchedule
--kubectl taint nodes uoc-neo2core-03 memorytype=low:NoSchedule

kubectl taint nodes uoc-zero2-01 uocnodetype=master:NoSchedule
kubectl taint nodes uoc-zero2-02 uocnodetype=master:NoSchedule
kubectl taint nodes uoc-zero2-03 uocnodetype=master:NoSchedule

# Labels

scylla.scylladb.com/node-type=scylla

kubectl label nodes uoc-neo2core-01 envoyLib=compatible
kubectl label nodes uoc-neo2core-02 envoyLib=compatible
kubectl label nodes uoc-neo2core-03 envoyLib=compatible


# Loadbalancer HA

https://www.google.com/search?client=firefox-b-lm&channel=entpr&q=k3s+load+balancer+external+ip