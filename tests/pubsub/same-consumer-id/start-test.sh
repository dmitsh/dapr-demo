#!/bin/bash

set -e

echo "Deploying Dapr component configs"
kubectl apply -f dapr-pubsub-rabbitmq.yml
echo "Deploying subscribers"
kubectl apply -f sub.yml
sleep 5
echo "Deploying publisher"
kubectl apply -f pub.yml
