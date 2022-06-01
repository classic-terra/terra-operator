Terra-Operator (v1)
======================================

The Terra-Operator is a community driven project focused on creating a Kubernetes native operator that will simplify the process of deploying TerradNodes and Validators to via Kubectl to any k8s cluster. The idea is to make it alot simpler for users to easily spin up a cluster with a few commands on any k8s resources available to them thus allowing the hashing power of our networks to grow (Note: It supports both Classic and V2).

## Getting started
These instructions will help you setup the Terra-Operator on your k8s cluster. If you find yourself in a situation where one of more tools might not be working please reach out to us for assistance on how to proceed, post an [issue in our repository](https://github.com/terra-rebels/terra-operator/issues), fix it yourself & update the code via a [pull request](https://github.com/terra-rebels/terra-operator/pulls) or reach out to us on [Discord](https://discord.gg/zW43ghuMpa).

### Prerequisites
* [Go v. 1.18+](https://go.dev/dl/)
* [Operator-sdk](https://sdk.operatorframework.io/docs/installation/)
* [MiniKube](https://minikube.sigs.k8s.io/docs/start/)

### Installing Terra-Operator
In order to get install the Terra-Operator the above prerequisites must be meet by the host machine and if you wish the run a full node (e.i a Validator) the machine must meet the following requirements: https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html. Once you have verified your system meets the required guidelines the process of getting the Terra-Operator installed is described below.

#### Cloning Terra-Operator repo
Clone the Terra-Operator from GitHub using the following command:

```
git clone https://github.com/terra-rebels/terra-operator.git
```

#### Apply Terra-operator yaml
Navigate to the deploy directory and apply the yaml files using the following commands:

```
cd deploy
minikube kubectl apply -f ./
```

#### Verify that validator is installed succesfully
Verify that Terra-Operator is running using the following command:

```
minikube kubectl get Deployment terra-operator -n terra
```

Which should yield something like this: `terra-operator   1/1     1            1           16m`

Congratulations you have now installed the Terra-Operator on your k8s cluster.

### TODO: Section on creating a TerradNode custom resource

- How to install CRD
- How to install CR (incl. configuration options)
- How to add a shared volume

### TODO: Section on creating a Validator custom resource

- How to install CRD
- How to install CR (incl. configuration options)
- How to add a shared volume

## Want to help make our dcoumentation better?
 * Want to **log an issue**? Feel free to visit our [GitHub site](https://github.com/terra-rebels/terra-operator/issues).
 
