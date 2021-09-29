#!/bin/bash

set -e
set -o pipefail

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"

echo $DIR

echo "Checking kubectl"
which kubectl

echo "Deploying rabbitmq cluster operator"
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
POD=$(kubectl -n rabbitmq-system get pod -l app.kubernetes.io/component=rabbitmq-operator -o jsonpath={.items..metadata.name})
kubectl wait -n rabbitmq-system --for=condition=Ready pod/$POD

echo "Deploying rabbitmq cluster"
kubectl apply -f ../deployments/pubsub/rabbitmq-cluster.yml
kubectl wait --for=condition=Ready pod/dapr-demo-server-0

echo "Deploying Dapr rabbitmq pubsub component"

USER=$(kubectl get secret dapr-demo-default-user -o jsonpath="{.data.username}" | base64 --decode)
PASSWORD=$(kubectl get secret dapr-demo-default-user -o jsonpath="{.data.password}" | base64 --decode)

cat ../deployments/pubsub/dapr-pubsub-rabbitmq.yml.templ | sed "s/<USERNAME>/${USER}/g" | sed "s/<PASSWORD>/${PASSWORD}/g" | cat > ../deployments/pubsub/dapr-pubsub-rabbitmq.yml

kubectl apply -f ../deployments/pubsub/dapr-pubsub-rabbitmq.yml

echo "Deploying pubsub app"
kubectl apply -f ../deployments/pubsub/pub.yml
kubectl apply -f ../deployments/pubsub/sub.yml

echo "Scaling up subscribers"
POD=$(kubectl get pod -l app=pubsub-sub -o jsonpath={.items..metadata.name})
kubectl wait --for=condition=Ready pod/$POD
kubectl scale deployment.apps/sub --replicas=9
