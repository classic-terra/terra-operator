#!/bin/bash

kubectl apply -f config/net-samples/atreides/volume.yaml
kubectl apply -f config/net-samples/atreides

kubectl apply -f config/net-samples/harkonnen/volume.yaml
kubectl apply -f config/net-samples/harkonnen

go run ./... network --replica 20