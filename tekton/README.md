export NAMESPACE=knative-serving
kubectl create clusterrolebinding $NAMESPACE:knative-serving-namespaced-admin \
--clusterrole=knative-serving-namespaced-admin  --serviceaccount=$NAMESPACE:default