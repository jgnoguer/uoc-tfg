# Secrets

kubectl create secret generic mqtt-credentials --from-file=mqtt-credentials.properties



mosquitto_pub -m "{ "msg": "message from mosquitto_pub client"} " -t "sensor-topic" -u jgnoguer -h 10.43.26.129 -P uocAn1m4ls