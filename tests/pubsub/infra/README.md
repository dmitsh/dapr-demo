# Testing environment

In these tests we are using RabbitMQ as the underlying pubsub service. However, you can use any other [Dapr pubsub component](https://docs.dapr.io/reference/components-reference/supported-pubsub/).

To keep track of sent and received messages, we have instrumented Prometheus client in our test application. Prometheus server is running alongside applications and scrapes the metrics.

We are running all components in a Kubernetes cluster, however, this test can be conducted in a self-hosted environment.

## Creating testing environment

### 1. Deploy Kubernetes cluster

Use an existent or deploy a new Kubernetes cluster.
Any Kubernetes implementation should work: [minikube](https://minikube.sigs.k8s.io/docs/), [kind](https://kind.sigs.k8s.io/), managed Kubernetes services (AKS, EKS, GCP, etc.), and many others.

### 2. Deploy RabbitMQ server

There are several ways to deploy RabbitMQ. For example, [RabbitMQ cluster operator](https://www.rabbitmq.com/kubernetes/operator/operator-overview.html) or [Bitnami helm chart](https://bitnami.com/stack/rabbitmq/helm), as shown below.

#### 2.1. RabbitMQ cluster operator

Deploy operator:
```bash
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
```
Deploy RabbitMQ server:
```bash
kubectl apply -f ../../../deployments/rabbitmq/rabbitmq-cluster.yml
```
For AKS, deploy RabbitMQ server with customized config:
```bash
kubectl apply -f ../../../deployments/rabbitmq/rabbitmq-cluster-aks.yml
```

#### 2.2. Bitnami helm chart

Add helm repo and deploy RabbitMQ server:
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami

helm install dapr-demo bitnami/rabbitmq \
  --set auth.username=guest \
  --set auth.password=guest \
  --set image.debug=true \
  --set persistence.size=1Gi
```
For AKS, use customized helm chart:
```bash
helm repo add azure-marketplace https://marketplace.azurecr.io/helm/v1/repo

helm install dapr-demo azure-marketplace/rabbitmq \
  --set auth.username=guest \
  --set auth.password=guest \
  --set image.debug=true \
  --set persistence.size=1Gi
```

### 3. Deploy Prometheus server

Use an existent or deploy a new Prometheus server. There are multiple ways to deploy Prometheus server. Feel free to use our setup:

```bash
kubectl create namespace monitoring
kubectl apply -f https://raw.githubusercontent.com/dmitsh/prometheus-deployment/main/deployments/rbac.yml
kubectl apply -f https://raw.githubusercontent.com/dmitsh/prometheus-deployment/main/deployments/configmap.yml
kubectl apply -f https://raw.githubusercontent.com/dmitsh/prometheus-deployment/main/deployments/deployment.yml
kubectl apply -f https://raw.githubusercontent.com/dmitsh/prometheus-deployment/main/deployments/service.yml
```
