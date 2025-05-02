# Notes

kubectl label nodes uoc-neo2core-01 longhorn=compatible
kubectl label nodes uoc-neo2core-02 longhorn=compatible
kubectl label nodes uoc-neo2core-03 longhorn=compatible



flux get helmrelease longhorn-release -n longhorn-system

kubectl -n longhorn-system get pod

kubectl -n longhorn-system get svc

https://longhorn.io/docs/1.8.1/deploy/accessing-the-ui/longhorn-ingress/
