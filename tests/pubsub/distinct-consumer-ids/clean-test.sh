#!/bin/bash

set -x

echo "Deleting test applications"
kubectl delete -f pubsub1.yml
kubectl delete -f pubsub2.yml
kubectl delete -f dapr-pubsub-rabbitmq1.yml
kubectl delete -f dapr-pubsub-rabbitmq2.yml