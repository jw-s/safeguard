#!/bin/bash

CN=safeguard.kube-system.svc
OUTPUT_DIR=./certs

set -e

rm -rf ./certs/*

openssl genrsa -out ${OUTPUT_DIR}/${CN}.key 2048
openssl req -new -key ${OUTPUT_DIR}/${CN}.key -out ${OUTPUT_DIR}/${CN}.csr -subj "/CN=${CN}"
openssl x509 -req -in ${OUTPUT_DIR}/${CN}.csr -CA ~/.minikube/certs/ca.pem -CAkey ~/.minikube/certs/ca-key.pem -CAcreateserial -out ${OUTPUT_DIR}/${CN}.crt -days 500 -sha256