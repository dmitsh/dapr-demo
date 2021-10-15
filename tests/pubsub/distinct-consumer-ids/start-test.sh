#!/bin/bash

set -e

echo "Deploying Dapr component configs"
kubectl apply -f dapr-pubsub-rabbitmq1.yml
kubectl apply -f dapr-pubsub-rabbitmq2.yml
echo "Deploying first application"
kubectl apply -f pubsub1.yml
echo "Waiting 30 seconds"
sleep 30
echo "Deploying second application"
kubectl apply -f pubsub2.yml
