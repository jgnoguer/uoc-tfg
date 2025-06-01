# Config

kubectl label nodes uoc-rock3a-01 camelk=compatible
kubectl label nodes uoc-rock3a-02 camelk=compatible
kubectl label nodes uoc-rock3a-03 camelk=compatible

kubectl create secret generic mail-credentials --from-file=mail-credentials.properties