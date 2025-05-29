# Knative configuration

kubectl get ksvc -A


flux reconcile kustomization infra-knative


## Build

export CR_PAT=ghp_*****
export FUNC_REGISTRY=ghcr.io/jgnoguer
echo $CR_PAT | docker login ghcr.io -u ****** --password-stdin

func build -v
func deploy -v
