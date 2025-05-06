# Flux install

https://fluxcd.io/flux/installation/bootstrap/github/#github-organization

https://github.com/fluxcd/terraform-provider-flux/tree/main/examples/github-via-pat

Create a GitHub personal access token and export it as an env var

export GITHUB_TOKEN=<my-token>

The fine-grained PAT must be generated with the following permissions:

    Administration -> Access: Read-only
    Contents -> Access: Read and write
    Metadata -> Access: Read-only
    
# Affinity and tolerations
  
  https://fluxcd.io/flux/installation/configuration/vertical-scaling/#node-affinity-and-tolerations


-- kubectl label nodes uoc-cubie role=flux
kubectl label nodes uoc-neo2core-01 role=flux
kubectl label nodes uoc-neo2core-02 role=flux
kubectl label nodes uoc-neo2core-03 role=flux
-- kubectl label nodes uoc-rock3a-01 role=flux

## Install with cli

flux bootstrap github \
  --token-auth \
  --owner=jgnoguer \
  --repository=uoc-tfg \
  --branch=main \
  --path=flux/clusters/uoc \
  --personal



