
Simple script to list and create Azure Resource Group.

```
az login

export AZURE_SUBSCRIPTION_ID=`az account show --query "id" -o tsv`

go build
```

#### \# list
```
./azure-rg-go list | jq
{
  "name": "demo",
  "details": {
    "id": "/subscriptions/.../resourceGroups/demo",
    "location": "westeurope",
    "state": "Succeeded",
    "tags": "no tags"
  }
}
```

#### \# create
```
./azure-rg-go create
2021/05/06 23:45:55 Created resource group: goRGtest

./azure-rg-go list | jq
{
  "name": "demo",
  "details": {
    "id": "/subscriptions/.../resourceGroups/demo",
    "location": "westeurope",
    "state": "Succeeded",
    "tags": "no tags"
  }
}
{
  "name": "goRGtest",
  "details": {
    "id": "/subscriptions/.../resourceGroups/goRGtest",
    "location": "westeurope",
    "state": "Succeeded",
    "tags": "no tags"
  }
}
```

#### \# delete
```
./azure-rg-go delete
2021/05/07 16:36:43 Deleted resource group: goRGtest
```