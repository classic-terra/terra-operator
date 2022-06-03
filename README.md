Terra-Operator (v1)
======================================

The Terra-Operator is a community driven project focused on creating a Kubernetes native operator that will simplify the process of deploying TerradNodes and Validators via Kubectl to any k8s cluster. The idea is to make it alot simpler for users to easily spin up a cluster with a few commands on any k8s resources available to them thus allowing the hashing power of our networks to grow (Note: It supports both Classic and V2).

## Getting started
These instructions will help you setup the Terra-Operator on your k8s cluster. If you find yourself in a situation where one of more tools might not be working please reach out to us for assistance on how to proceed, post an [issue in our repository](https://github.com/terra-rebels/terra-operator/issues), fix it yourself & update the code via a [pull request](https://github.com/terra-rebels/terra-operator/pulls) or reach out to us on [Discord](https://discord.gg/zW43ghuMpa).

### Prerequisites
* [Go v. 1.18+](https://go.dev/dl/)
* [Operator-sdk](https://sdk.operatorframework.io/docs/installation/)
* [MiniKube](https://minikube.sigs.k8s.io/docs/start/)

### Installing Terra-Operator
In order to install the Terra-Operator the above prerequisites must be meet by the host machine and if you wish to run a full node (e.i a Validator) the machine must meet the following requirements: https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html. Once you have verified your system meets the minimum requirements the process of getting the Terra-Operator installed is fairly simple. :)

#### Cloning Terra-Operator repo
Clone the Terra-Operator from GitHub using the following command:

```
git clone https://github.com/terra-rebels/terra-operator.git
```

#### Apply Terra-Operator yaml
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

This should yield something like this: `terra-operator   1/1     1            1           16m`

Congratulations you have now installed the Terra-Operator on your k8s cluster.

### TerradNode CRD (v1alpha)
The TerradNode CRD is a custom resource definition managed by the Terra-Operator which provides the base layer for any terra node running on Kubernetes. Its job is simply to spin up a `terrad` daemon running in a initialized state with Tendermint consensus (BPOS) and networking components targeting the `ChainId` using the desired `NodeImage` client.

The TerradNode consists of a `pod` running the `NodeImage` client vs a desired version of the Terra blockchain identified by the `ChainId` with the following containerPorts exposed: `1317` (LCD), `26656` (P2P), `26657` (RPC) & `26660` (Prometheus). Furthermore it kick-starts the `terrad start` command to ensure the node is initialized and running either as a `light-node` or a `full-node` depending on the `IsFullNode` value of the TerradNodeSpec. 

#### How to install TerradNode CRD
From the root of the Terra-Operator repo run the following command:

```
minikube kubectl apply -f ./deploy/crds/terra.rebels.info_terradnodes_crd.yaml
```

Verify that kubectl prints the following message: `customresourcedefinition.apiextensions.k8s.io/terradnodes.terra.rebels.info created`

#### How to create an TerradNode instance
From the root of the Terra-Operator repo run the following command:

```
minikube kubectl apply -f ./deploy/crds/terra.rebels.info_v1alpha1_terradnode_cr.yaml
```

Verify that kubectl prints the following message: `terradnode.terra.rebels.info/example-terradnode created`

##### TerradNodeSpec configuration options
The `terra.rebels.info_v1alpha1_terradnode_cr.yaml` example supports the following configuration options:

```
spec:
  nodeImage: terradnode-container-image (required - string)
  isFullNode: terradnode-light-or-full (optional - bool)
  postStartCommand: terradnode-post-start-command (optional - string)
  dataVolume: (optional)
    name: my-nfs-share (required - string)
    nfs:
      server: my.nfs.share (required - string)
      path: /my-nfs-share/ (required - string)
```

### Validator CRD (v1alpha)
The Validator CRD is a custom resource definition managed by the Terra-Operator that mounts a Validator on top of a TerradNode resource and runs it in a bonded mode using the configured Application Oracle Key (`terrad tx create-validator --from` arg). A Validators responsibility is to spin up a `terrad` daemon running as a `full-node`, mount it on a volume containing the desired blockchain snapshot (can be found at https://quicksync.io/networks/terra.html) and bootstraps a `PostStartupScript` command on the TerradNode ContainerSpec that executes the required commands to succesful launch the Terra Validator with the desired ValidatorSpec.

Lastly when running in "public-mode" (`IsPublic: true`) it also creates a `service` which can route ingress traffic to the containerPorts from clients outside your Kubernetes cluster.

#### How to install Validator CRD
From the root of the Terra-Operator repo run the following command:

```
minikube kubectl apply -f ./deploy/crds/terra.rebels.info_validators_crd.yaml
```

Verify that kubectl prints the following message: `customresourcedefinition.apiextensions.k8s.io/validators.terra.rebels.info created`

#### How to create an Validator instance
From the root of the Terra-Operator repo run the following command:

```
minikube kubectl apply -f ./deploy/crds/terra.rebels.info_v1alpha1_validator_cr.yaml
```

Verify that kubectl prints the following message: `validator.terra.rebels.info/example-validator created`

##### ValidatorSpec configuration options
The `terra.rebels.info_v1alpha1_validator_cr.yaml` example supports the following configuration options:

```
  chainId: validator-chain-id (required - string)
  nodeImage: validator-node-image (required - string)
  name: validator-moniker (required - string)
  fromKeyName: validator-application-oracle-key-name (required - string)
  initialCommissionRate: validator-commission-rate (required - string)
  maximumCommission: validator-commission-max-rate (required - string)
  commissionChangeRate: validator-commission-max-change-rate (required - string)
  minimumSelfBondAmount: validator-min-self-delegation (required - string)
  initialSelfBondAmount: validator-amount (required - string)
  isPublic: validator-is-accessible-from-outside-cluster (optional - bool)
  description: validator-description (optional - string)
  website: validator-website (optional - string)
  dataVolume: (optional)
    name: my-nfs-share (required - string)
    nfs:
      server: my.nfs.share (required - string)
      path: /my-nfs-share/ (required - string)
```

## Want to help make our documentation better?
 * Want to **log an issue**? Feel free to visit our [GitHub site](https://github.com/terra-rebels/terra-operator/issues).
 
