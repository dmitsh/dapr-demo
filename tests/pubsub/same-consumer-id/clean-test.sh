#!/bin/bash

set -x

echo "Deleting test applications"
kubectl delete -f sub.yml
kubectl delete -f pub.yml
kubectl delete -f dapr-pubsub-rabbitmq.yml
