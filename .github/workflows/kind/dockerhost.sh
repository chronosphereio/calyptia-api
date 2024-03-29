#!/bin/bash

set -ex

cat <<EOF | kubectl apply -f -
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dockerhost
  labels:
    k8s-app: dockerhost
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: dockerhost
  template:
    metadata:
      labels:
        k8s-app: dockerhost
    spec:
      containers:
        - name: dockerhost
          image: qoomon/docker-host
          securityContext:
            capabilities:
              add: ["NET_ADMIN", "NET_RAW"]
          env:
            - name: DOCKER_HOST
              value: 172.17.0.1
---
apiVersion: v1
kind: Service
metadata:
  name: dockerhost
spec:
  clusterIP: None
  selector:
    k8s-app: dockerhost
EOF

cat ~/.kube/config
