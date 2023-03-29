#!/bin/bash

kubectl delete statefulset --all
kubectl delete pvc --selector=node=atreides
kubectl delete pvc --selector=node=harkonnen