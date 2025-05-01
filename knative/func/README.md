
# Notes about KNative

## Infrastructure

Getting all the services deployed

kubectl get ksvc -A

## Infra reconciliation

flux reconcile kustomization infra-knative

## Registry configuration

https://knative.dev/docs/serving/deploying-from-private-registry/

   80  doctl auth init --context uoc-do
   81  doctl account get
   82  doctl auth list
   83  doctl auth switch --context uoc-do
   84  doctl auth list
   85  doctl account get
   86  doctl registry kubernetes-manifest | kubectl apply -f -
   87  doctl registry kubernetes-manifest
   89  doctl registry kubernetes-manifest | kubectl apply -f -
  127  doctl registry kubernetes-manifest

kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "registry-dockerhub"}]}'


kubectl annotate kservice helloworld-go networking.knative.dev/ingress-class=


kubectl create secret docker-registry registry-dockerhub \
  --docker-server=docker.io/jgnoguer \
  --docker-email=jgarciano@uoc.edu \
  --docker-username=jgnoguer \
  --docker-password=<dockerpat>

  func build --registry docker.io/jgnoguer 

or

export FUNC_REGISTRY=docker.io/jgnoguer

## Enable persistent volumes

kubectl patch --namespace knative-serving configmap/config-features \
 --type merge \
 --patch '{"data":{"kubernetes.podspec-persistent-volume-claim": "enabled", "kubernetes.podspec-persistent-volume-write": "enabled"}}'

 ## Creating a REST service

 https://github.com/dewitt/knative-docs/blob/master/serving/samples/rest-api-go/README.md

 