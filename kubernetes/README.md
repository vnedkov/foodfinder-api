# Deploying the service
## Prerequisites
* A Running Kubernetes cluster and a client configured (kubectl) with the right permissions
* An Elasticsearch instance, accessible via HTTP(S). 

## Create a namespace
```sh
kubectl apply -f namespace.yaml
```

## Create a secret to store Elasticsearch connection parameters
Here is an example. Replace the values according your setup:
```shell
kubectl create secret generic elasticsearch-secret \
-n foodfinder # in the same namespace where the service will be deployed\
--from-literal=ELASTICSEARCH_URL='http://localhost:9200' # Must use single quotes to escape special characters \
--from-literal=ELASTICSEARCH_USER=elastic \
--from-literal=ELASTICSEARCH_PASSWORD='elastic-user-password'
--from-literal=ELASTICSEARCH_INDEX=foodfinder # The name of the index to search
```
** Please note, anyone else having access to the pod can see the secret material. This approach is acceptable for a homelab, but not in a cloud environment with multiple people having access.

## Deploy the service along with Cilium policies allowing access from outside the cluster
```sh
kubectl apply -f deploy.yaml
```