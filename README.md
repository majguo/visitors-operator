# Visitors Operator

This repo contains the source code to build a sample Kubernetes Operator named `Visitors Operator`. It's based on book [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/) and its code base [visitors-operator](https://github.com/jdob/visitors-operator), but refactored to work with the current version of [operator-sdk](https://sdk.operatorframework.io/docs/installation/).

Follow instructions below to deploy and run the Operator.

## Prerequisites

You must complete the prerequisites listed in this section before moving on to next step.

1. You will need an Azure subscription. If you don't have one, you can get one for free for one year [here](https://azure.microsoft.com/free).
1. Install [Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli?view=azure-cli-latest&preserve-view=true).
1. Install [Docker](https://docs.docker.com/get-docker/) for your OS.

### Install Operator SDK

If you haven't installed Operator SDK, pls install it from GitHub release by following [this guide](https://sdk.operatorframework.io/docs/installation/#install-from-github-release).

> **Note**
> The Operator was originally developed with operator-sdk v1.18.1, and upgraded to v1.22.2. So please make sure download `v1.22.2` for installation:
>
> ```
> export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.22.2
> ```

After installation, run `operator-sdk version` in your CLI, you should see the similar output below:

```
operator-sdk version: "v1.22.2", commit: "da3346113a8a75e11225f586482934000504a60f", kubernetes version: "1.24.1", go version: "go1.18.4", GOOS: "linux", GOARCH: "amd64"
```

Otherwise, pls troubleshoot before returning here to continue.

### Deploy an Azure Kubernetes Service cluster

Execute the following commands to deploy an Azure Kubernetes Service (AKS) cluster.

```
RESOURCE_GROUP_NAME=jiangma-aks-cluster
az group create --name $RESOURCE_GROUP_NAME --location eastus

CLUSTER_NAME=myAKSCluster
az aks create --resource-group $RESOURCE_GROUP_NAME --name $CLUSTER_NAME --node-count 1 --generate-ssh-keys --enable-managed-identity
```

If you encoutner the issue `Please run 'az login' to setup account.`, login to Azure first.


Then connect to the created AKS cluster.

```
az aks get-credentials --resource-group $RESOURCE_GROUP_NAME --name $CLUSTER_NAME --overwrite-existing
```

To verify the connection to your cluster, use the `kubectl get` command to return a list of the cluster nodes.

```
kubectl get nodes
```

You should see one cluster node returned if connection successfuly set up. Otherwise troubleshoot before continuing. You can find more details about deploying an AKS cluster using Azure CLI with [Quickstart: Deploy an Azure Kubernetes Service cluster using the Azure CLI](https://docs.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli).

## Prepare application images

The Visitors Operator will automatically deploy 3 kinds of servers consisting of [`visitors-webui`](./visitors-webui/), [`visitors-service`](./visitors-service/) and `mysql` underneath.

Except `mysql` image, `visitors-webui` and `visitors-service` depend on images generated from [`visitors-webui`](./visitors-webui/) and [`visitors-service`](./visitors-service/).

The Operator will use default images `majguo/visitors-webui` for `visitors-webui`, and `majguo/visitors-service` for `visitors-service`, if you don't specify them in the CRD, e.g., [apps_v1alpha1_visitorsapp.yaml](./config/samples/apps_v1alpha1_visitorsapp.yaml):

```
apiVersion: apps.example.com/v1alpha1
kind: VisitorsApp
metadata:
  name: visitorsapp-sample
spec:
  size: 3
```
 
Alternatively, you can also specify images by yourself, e.g., [apps_v1alpha1_visitorsapp_all.yaml](./config/samples/apps_v1alpha1_visitorsapp_all.yaml):

```
apiVersion: apps.example.com/v1alpha1
kind: VisitorsApp
metadata:
  name: visitorsapp-sample-all
spec:
  size: 3
  title: "Custom Dashboard Title"
  frontendImage: "majguo/visitors-webui:1.0.0"
  backendImage: "majguo/visitors-service:1.0.0"
```

Here're steps on how you can build, tag, push and specify images.

1. Follow instructions in [build-visitors-webui-app-image](./visitors-webui/README.md#build-application-image).

1. Follow instructions in [build-visitors-service-app-image](./visitors-service/README.md#build-application-image).

1. Replace values for `frontendImage` and `backendImage` in [apps_v1alpha1_visitorsapp_all.yaml](./config/samples/apps_v1alpha1_visitorsapp_all.yaml) with your own images.

## Run operator

You have several approaches to run the operator, depending on your purpose.

First, clone this repo and change to root directory of the local clone.

### Run operator locally outside the cluster

To run locally for development purposes and outside of a cluster, run the following target using `make`.

```
make install run
```

### Run operator as a Deployment inside the cluster

By default, a new namespace is created with name <project-name>-system, e.g., `visitors-operator-system`, and will be used for the deployment.

Run the following to deploy the operator. This will also install the RBAC manifests from `config/rbac`.

```
make docker-build docker-push IMG="<DockerHub-account>/visitors-operator"
make deploy IMG="<DockerHub-account>/visitors-operator"
```

Once you excuted the above commands, you can use `fast-deploy` target to only deploy operator next time.

```
make fast-deploy IMG="<DockerHub-account>/visitors-operator"
```

## Deploy sample CR

To verify if the operator works as expected, run the commands below to deploy the sample CR(s).

```
# Deploy a sample CR with required fields
kubectl apply -f config/samples/apps_v1alpha1_visitorsapp.yaml

# Deploy a sample CR with all fields
kubectl apply -f config/samples/apps_v1alpha1_visitorsapp_all.yaml
```

Execute the following commands to check if deployment succeeded.

```
kubectl get VisitorsApp
kubectl get deployment
```

You should see different resources output in the CLI with **READY** state. Then run the following command to monitor status of pods and wait until they become **Running**.

```
kubectl get pod -o wide -w
```

Press `CONTROL-C` to stop monitoring.

Finally run the following command to get public IP address and port for front-end service.

```
kubectl get servie
```

Copy value of `EXTERNAL-IP` and open it with port `3000` in the browser. You should see one entry with **Service IP** and **Client IP** listed in the table of the page. Refreshing the page will adding more similar entries. 

## Cleanup

Once you completed the walkthrough of this operator learning module, clean up the resources as below.

First, delete CR(s) deployed before.

```
# Delete the sample CR with required fields if you deployed before
kubectl delete -f config/samples/apps_v1alpha1_visitorsapp.yaml

# Delete the sample CR with all fields if you deployed before
kubectl delete -f config/samples/apps_v1alpha1_visitorsapp_all.yaml
```

Then stop running of the operator. 

* Press `CONTROL-C` if you ran operator locally outside the cluster. 

* Execute the following command if you ran operator as a Deployment inside the cluster.

  ```
  make undeploy
  ```
   
Finally delete the AKS cluster.

```
az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes
```

## References

* [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/)
* [github.com/jdob/visitors-operator](https://github.com/jdob/visitors-operator)
* [Installation Guide for the Operator SDK CLI](https://sdk.operatorframework.io/docs/building-operators/golang/installation/)
* [Go Operator Tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
* [Quickstart: Deploy an Azure Kubernetes Service cluster using the Azure CLI](https://docs.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli)
