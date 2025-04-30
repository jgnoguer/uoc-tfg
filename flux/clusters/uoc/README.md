# For flux

flux bootstrap github \
  --token-auth \
  --owner=jgnoguer \
  --repository=uoc-tfg \
  --branch=main \
  --path=flux/clusters/uoc \
  --personal

  # Affinity and tolerations
  
  https://fluxcd.io/flux/installation/configuration/vertical-scaling/#node-affinity-and-tolerations



kubectl label nodes uoc-cubie role=flux
kubectl label nodes uoc-zero2-01 role=flux
kubectl label nodes uoc-zero2-02 role=flux
kubectl label nodes uoc-zero2-03 role=flux
kubectl label nodes uoc-rock3a-01 role=flux

