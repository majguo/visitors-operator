# Visitors Operator

This repo contains the source code to build a sample Kubernetes Operator named `Visitors Operator`. It's based on book [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/) and its code base [visitors-operator](https://github.com/jdob/visitors-operator), but refactored to work with the current version of [operator-sdk](https://sdk.operatorframework.io/docs/installation/).

Follow instructions below to deploy and run the Operator.

## Prerequisites

TODO.

1. Deploy an AKS cluster.
1. Prepare app images for [`visitors-webui`](./visitors-webui/) and [`visitors-service`](./visitors-service/).

## Running operator from local

TODO.

## Running operator from the Kubernetes cluster

TODO.

## References

* [Kubernetes Operators](https://www.oreilly.com/library/view/kubernetes-operators/9781492048039/)
* [github.com/jdob/visitors-operator](https://github.com/jdob/visitors-operator)
* [Installation Guide for the Operator SDK CLI](https://sdk.operatorframework.io/docs/building-operators/golang/installation/)
* [Go Operator Tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
