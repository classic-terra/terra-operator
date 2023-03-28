# terra-operator

## Description
The terra-operator is a project focused on creating a Kubernetes native operator that will simplify the process of deploying Validator nodes to any k8s cluster.

## Getting started
These instructions will help you setup the terra-operator on your k8s cluster. If you find yourself in a situation where one of more tools might not be working please reach out to us for assistance on how to proceed, post an [issue in our repository](https://github.com/terra-rebels/terra-operator/issues), fix it yourself & update the code via a [pull request](https://github.com/terra-rebels/terra-operator/pulls) or reach out to us on [Discord](https://discord.gg/zW43ghuMpa).

### Prerequisites
* [Go v. 1.18+](https://go.dev/dl/)
* [Operator-sdk](https://sdk.operatorframework.io/docs/installation/)
* [KIND](https://sigs.k8s.io/kind)

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=public.ecr.aws/classic-terra/terraclassic.operator
```
	
2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=public.ecr.aws/classic-terra/terraclassic.operator
```

3. Install Instances of Custom Resources:

```sh
kubectl apply -f config/test-samples/
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
 * Want to **log an issue**? Feel free to visit our [GitHub site](https://github.com/terra-rebels/terra-operator/issues).

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
