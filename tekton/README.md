export NAMESPACE=knative-serving
kubectl create clusterrolebinding $NAMESPACE:knative-serving-namespaced-admin \
--clusterrole=knative-serving-namespaced-admin  --serviceaccount=$NAMESPACE:default


apiVersion: v1
kind: Secret
metadata:
  name: basic-user-pass
  annotations:
    tekton.dev/git-0: https://github.com # Described below
type: kubernetes.io/basic-auth
stringData:
  username: <cleartext username>
  password: <cleartext password>


apiVersion: v1
kind: ServiceAccount
metadata:
  name: build-bot
secrets:
  - name: basic-user-pass


https://tekton.dev/docs/pipelines/additional-configs/

Edit feature-flags config map: coschedule: pipelineruns

  https://tekton.dev/docs/pipelines/workspaces/#specifying-workspace-order-in-a-pipeline-and-affinity-assistants