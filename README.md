# Visitors Operator

This repo contains the source code to build a sample Kubernetes Operator named `Visitors Operator`. It's based on book [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/) and its code base [visitors-operator](https://github.com/jdob/visitors-operator), but refactored to work with the current version of [operator-sdk](https://sdk.operatorframework.io/docs/installation/).

Follow instructions below to deploy and run the Operator.

## Prerequisites

You must complete the prerequisites listed here before moving on to next step.

### Install Operator SDK

If you haven't installed Operator SDK, pls install it from GitHub release by following [this guide](https://sdk.operatorframework.io/docs/installation/#install-from-github-release).

After installation, run `operator-sdk version` in your CLI, you should see the similar output below:

```
operator-sdk version: "v1.22.2", commit: "da3346113a8a75e11225f586482934000504a60f", kubernetes version: "1.24.1", go version: "go1.18.4", GOOS: "linux", GOARCH: "amd64"
```

Otherwise, pls troubleshoot before returning here to continue.

### Prepare application images

The Visitors Operator will automatically deploy 3 kinds of servers consisting of [`visitors-webui`](./visitors-webui/), [`visitors-service`](./visitors-service/) and `mysql` underneath.

Except `mysql` image, you need to prepare images for `visitors-webui` and `visitors-service`.

1. Follow instructions in [build-visitors-webui-app-image](./visitors-webui/README.md#build-application-image) to generate and push `visitors-webui` image to DokcerHub.

1. Follow instructions in [build-visitors-service-app-image](./visitors-service/README.md#build-application-image) to generate and push `visitors-service` image to DokcerHub.

### Deploy an Azure Kubernetes Cluster

TODO.

## Running operator from local

TODO.

## Running operator from the Kubernetes cluster

TODO.

## References

* [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/)
* [github.com/jdob/visitors-operator](https://github.com/jdob/visitors-operator)
* [Installation Guide for the Operator SDK CLI](https://sdk.operatorframework.io/docs/building-operators/golang/installation/)
* [Go Operator Tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
