curl -v "http://broker-ingress.knative-eventing.svc.cluster.local/default/image-broker" \
-X POST \
-H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
-H "Ce-Specversion: 1.0" \
-H "Ce-Type: dev.jgnoguer.knative.uoc.imageadded" \
-H "Ce-Source: dev.jgnoguer.knative.uoc/mediastorage-service" \
-H "Content-Type: application/json" \
-d '{"msg":"Hello World from the curl pod."}'
exit


curl -v "http://broker-ingress.knative-eventing.svc.cluster.local/default/activities-broker" \
-X POST \
-H "Ce-Id: 12324324-88be-4077-9d7a-a3f162705f79" \
-H "Ce-Specversion: 1.0" \
-H "Ce-Type: dev.jgnoguer.knative.uoc.activitystarted" \
-H "Ce-Source: dev.jgnoguer.knative.uoc/activity-service" \
-H "Content-Type: application/json" \
-d '{"from":"jesus@jgnoguer.es"}'

curl -v "http://broker-ingress.knative-eventing.svc.cluster.local/default/demo-broker" \
-X POST \
-H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
-H "Ce-Specversion: 1.0" \
-H "Ce-Type: dev.jgnoguer.knative.uoc.activitystarted" \
-H "Ce-Source: dev.jgnoguer.knative.uoc/activity-service" \
-H "Content-Type: application/json" \
-d '{"msg":"Hello World from the curl pod."}'