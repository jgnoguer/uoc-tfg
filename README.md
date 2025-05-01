# uoc-tfg
TFG - Aplicacions i sistemes distribuïts


## Longhorn

https://longhorn.io/kb/tip-only-use-storage-on-a-set-of-nodes/

## k3s install

### Server

curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --disable traefik" sh -

--disable traefik

   34  curl -sfL https://get.k3s.io | K3S_URL=https://uoc-cubie:6443 K3S_TOKEN=token sh -
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

K3S_TOKEN=
check /var/lib/rancher/k3s/server/node-token

Kubeconfig

 /etc/rancher/k3s/k3s.yaml


On nanopi core
curl -sfL https://get.k3s.io | K3S_URL=https://uoc-cubie:6443 K3S_TOKEN=theks3token INSTALL_K3S_EXEC="agent --snapshotter=native" sh -


jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo cp /etc/rancher/k3s/k3s.yaml .
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ sudo chown jgnoguer:jgnoguer k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ mv k3s.yaml ~
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ export KUBECONFIG=/home/jgnoguer/k3s.yaml 
jgnoguer@kiwi:~/uocWksp/repo/knative/func/uoc-test$ kubectl get nodes


# Cubie images

https://github.com/Misaka-Nnnnq/Radxa_A5E_Firmware

# KNative func

https://github.com/knative/func/releases

## Registry

## ScyllaDB

https://operator.docs.scylladb.com/stable/index.html
