
# Notes about KNative

## Infrastructure

Getting all the services deployed

kubectl get ksvc -A

kubectl --namespace kourier-system get service kourier

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


kubectl annotate kservice helloworld-go networking.knative.dev/ingress-class=


kubectl create secret docker-registry registry-dockerhub \
  --docker-server=docker.io/jgnoguer \
  --docker-email=jgarciano@uoc.edu \
  --docker-username=jgnoguer \
  --docker-password=<dockerpat>

kubectl create secret docker-registry registry-ghrc \
--docker-server=ghcr.io/jgnoguer \
--docker-email=jgarciano@uoc.edu \
--docker-username=jgnoguer \
--docker-password=<dockerpat>

--namespace=knative-samples \

kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "registry-ghrc"}]}'
kubectl patch serviceaccount default -n knative-samples -p '{"imagePullSecrets": [{"name": "registry-ghrc"}]}'

  func build --registry ghcr.io/jgnoguer 

or

export FUNC_REGISTRY=ghcr.io/jgnoguer

## Enable persistent volumes

See knative kustomize

 ## Creating a REST service

 https://github.com/dewitt/knative-docs/blob/master/serving/samples/rest-api-go/README.md

 ## Check

kubectl --namespace istio-system get service istio-ingressgateway

 curl -H "Host: mediastorage.default.192.168.2.1.sslip.io" http://192.168.2.1:31287 -v
 curl -H "Host: agents.default.192.168.2.1.sslip.io" http://192.168.2.1:31287 -v

 kubectl get ksvc -A

 ## Create a new function

 func create -l go hello

 ## Create a new service 

 kn service create gitopstest --image knativesamples/helloworld --target=/user/knfiles/test.yaml


kn service apply -f mediastorage-service.yaml
kn service apply -f agents-service.yaml