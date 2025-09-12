# Deploying the service
## Prerequisites
* A Running Kubernetes cluster and a client configured (kubectl) with the right permissions
* An Elasticsearch instance, accessible via HTTP(S). 

## Create a namespace
Namespace should be created outside of the deployment to avoid destruction of other resources (secrets i.e.) when deleting the application deployment.
```sh
kubectl apply -f namespace.yaml
```

## Create a secret to store Elasticsearch connection parameters
Here is an example. Replace the values according your setup. You must use single quotes to escape special characters:
```shell
kubectl create secret generic elasticsearch-secret \
-n foodfinder \
--from-literal=ELASTICSEARCH_URL='http://localhost:9200' \
--from-literal=ELASTICSEARCH_USER=elastic \
--from-literal=ELASTICSEARCH_PASSWORD='elastic-user-password' \
--from-literal=ELASTICSEARCH_INDEX=foodfinder # The name of the index to search
```
** Please note, anyone else having access to the pod can see the secret material. This approach is acceptable for a homelab, but not in a cloud environment with multiple people having access.

## Deploy the service along with Cilium policies allowing access from outside the cluster
```sh
kubectl apply -f deploy.yaml
```