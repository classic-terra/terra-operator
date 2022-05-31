Terra-operator (v1)
======================================

The Terra-Operator is a community driven project focused on creating a Kubernetes native operator that will simplify the process of deploying TerradNodes and Validators to via Kubectl to any k8s cluster. The idea is to make it alot simpler for users to easily spin up a cluster with a few commands on any k8s resources available to them thus allowing the hashing power of our networks to grow (yes it supports both Classic and V2).

## Getting started
These instructions will help you setup the Terra-operator on your k8s cluster. If you find yourself in a situation where one of more tools might not be working please reach out to us for assistance on how to proceed, post an [issue in our repository](https://github.com/terra-rebels/terra-operator/issues), fix it yourself and update the kata via a [pull request](https://github.com/terra-rebels/terra-operator/pulls) or reach out to us on [Discord]().

### Prerequisites
* [Go v. 1.18+](https://go.dev/dl/)
* [Operator-sdk](https://sdk.operatorframework.io/docs/installation/)
* [MiniKube](https://minikube.sigs.k8s.io/docs/start/) (or any other k8s cluster)



### 1. Create a kata directory
First we setup a directory for our exercise files. It's pretty straight forward:

```
mkdir kata1
cd kata1
```

### 2. Create a new Web API .NET Core project
Then we create some boilerplate code to test ApplicationInsights:

```
dotnet new webapi
```

Just to explain: <br/>
`dotnet` - is the dotnet CLI <br/>
`new` - instructs the dotnet CLI to create a new project in the current directory.<br/>
`webapi` - tells the dotnet CLI which project template to use


## Want to help make our dcoumentation better?
 * Want to **log an issue**? Feel free to visit our [GitHub site](https://github.com/terra-rebels/terra-operator/issues).
 