# Notes

kubectl label nodes uoc-rock3a-01 longhorn=compatible
kubectl label nodes uoc-rock3a-01 longhorn-ui=compatible
kubectl label nodes uoc-rock3a-02 longhorn=compatible
kubectl label nodes uoc-rock3a-02 longhorn-ui=compatible
kubectl label nodes uoc-rock3a-03 longhorn=compatible
kubectl label nodes uoc-rock3a-03 longhorn-ui=compatible
kubectl label nodes uoc-rpicm4-01 longhorn=compatible
kubectl label nodes uoc-rpicm4-01 longhorn-ui-
kubectl label nodes uoc-rpicm4-02 longhorn=compatible
kubectl label nodes uoc-rpicm4-02 longhorn-ui=-


flux get helmrelease longhorn-release -n longhorn-system

kubectl -n longhorn-system get pod

kubectl -n longhorn-system get svc

https://longhorn.io/docs/1.8.1/deploy/accessing-the-ui/longhorn-ingress/

USER=<USERNAME_HERE>; PASSWORD=<PASSWORD_HERE>; echo "${USER}:$(openssl passwd -stdin -apr1 <<< ${PASSWORD})" >> auth

kubectl -n longhorn-system create secret generic basic-auth --from-file=auth

helm uninstall longhorn-release -n longhorn-system
helm install -f flux/infrastructure/longhorn/values.yaml longhorn-release longhorn-repo/longhorn -n longhorn-system

flux reconcile kustomization infra-longhorn
flux get helmrelease longhorn-release -n longhorn-system

kubectl -n longhorn-system edit settings.longhorn.io deleting-confirmation-flag

kubectl -n longhorn-system port-forward service/longhorn-frontend 8080:80



## stuck release

flux suspend hr longhorn-release -n longhorn-system
flux resume hr longhorn-release -n longhorn-system