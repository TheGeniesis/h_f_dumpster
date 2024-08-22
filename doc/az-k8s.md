# az-k8s

## k8s

docker build --no-cache --progress=plain -t my-image .

### events

```shell
export NAMESPACE=<nas> && kubectl -n ${NAMESPACE} get events --sort-by=.metadata.creationTimestamp
```

### logs

```shell
NAMESPACE=<ns> kubectl -n ${NAMESPACE} logs $(kubectl -n ${NAMESPACE} get pods | grep "<pattern_in_name>" | awk '{print $1}')
```

### Delete pods

```shell
NAMESPACE=<ns>; kubectl -n ${NAMESPACE} delete pod $(echo $(kubectl -n ${NAMESPACE} get pods | grep "<pattern_in_name>" | awk '{print $1}'))
```

#### Force delete

```shell
export NAMESPACE=<ns> && kubectl -n ${NAMESPACE} delete --grace-period=0 --force pod $(echo $(kubectl -n ${NAMESPACE} get pods | grep "<pattern_in_name>" | awk '{print $1}'))
```

### Delete all releases in NS

```shell
NAMESPACE=<ns> a=$(export NUM_DAYS=3 && ./helm ls -o json -f '^.*' --namespace ${NAMESPACE} | jq -r --argjson num_days "$NUM_DAYS" '.[] | select (.updated | sub("\\..*";"Z") | sub("\\s";"T") | fromdate < now - $num_days*86400)' | jq '.name' | xargs -I{} | grep -v -e "<excluded_release_pattern>") && ./helm -n ${NAMESPACE} delete $(echo $a)
```

### Deploy

```shell
RELEASE_NAMESPACE=<ns> ./helmfile --environment=<app_env> apply
```

### Login azure

```shell
az account set --subscription "<subscription>"
az aks get-credentials --resource-group <rg> --name <k8s-name>

```

### Generate databricks access token

```shell

az account set --subscription <subscription>
az login --service-principal -u <application_id> -p <password> --tenant <tenantid>
export DATABRICKS_TOKEN=$(az account get-access-token --resource 2ff814a6-3304-4ab8-85cb-cd0e6f879c1d --query "accessToken" --output tsv)

curl -X POST '<databricks_host>/api/2.0/token/create' --http1.1 --header "Authorization: Bearer ${DATABRICKS_TOKEN}"  --data-raw '{"lifetime_seconds":7776000,"comment":"token"}' -vvv

az login --service-principal -u <sp_app_id> -p <sp_password> --tenant <tenant_id>
```

### Login Rancher

```shell
cp ~/.kube/config ~/.kube/config.bak && KUBECONFIG=~/.kube/config:<path_to_download>/Downloads/devl-cloudops-k8s.yaml kubectl config view --flatten > /tmp/config && mv /tmp/config ~/.kube/config && kubectl config use-context devl-cloudops-k8s
```

### Docker login and promote

```shell
    AZURE_REGISTRY=<registry>
    IMAGE=busybox
    VERSION=latest
    full_image_name=$(echo "${IMAGE}" | sed 's/^\///')
    docker pull ${full_image_name}:${VERSION}
    docker tag ${full_image_name}:${VERSION} ${AZURE_REGISTRY}/test/${IMAGE}:${VERSION}
    az acr login -n <username>
    docker push ${AZURE_REGISTRY}/test/${IMAGE}:${VERSION}
```

### Generate secrets

```shell
kubectl -n <ns> create secret generic <secret_name> --from-literal=<field_name>=<field_value> --from-literal=<password_field_name>='<password>'
```

## Debug pod

```shell
kubectl -n <ns> run -it busybox --rm --overrides='
{
  "apiVersion": "v1",
  "spec": {
    "securityContext": {
      "runAsNonRoot": true,
      "runAsUser": 1001,
      "runAsGroup": 1001,
      "fsGroup": 1001
    },
    "containers": [
      {
        "name": "busybox",
        "image": "busybox",
        "stdin": true,
        "stdinOnce": true,
        "tty": true,
        "securityContext": {
          "allowPrivilegeEscalation": false
        }
      }
    ]
  }
}
' --image=busybox --restart=Never -- bash
```

## Check openssl cert

```sh
# create a temp pod (will fail command, but still running)
kubectl run --rm utils -it --restart=Never --image alpine/openssl -- s_client "-connect=<smtp_host>:25" "-starttls=smtp"
kubectl exec -it utils -- sh
openssl s_client -connect <smtp_host>:25 -starttls smtp
curl -vLk <smtp_host>:25
kubectl delete pod utils
```
