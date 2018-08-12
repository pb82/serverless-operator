# Serverless Operator

A Kubernetes / Openshift Operator, based on the [Operator SDK](https://github.com/operator-framework/operator-sdk)
that manages (creates and deletes) [Openwhisk](https://openwhisk.apache.org) actions.

## Installation on Openshift

*Prerequisite*: Openwhisk deployed to a namespace. You can follow the [instructions](https://github.com/projectodd/openwhisk-openshift) to set up Openwhisk.

Download the resources:

```sh
$ curl https://raw.githubusercontent.com/pb82/serverless-operator/master/deploy/operator.yaml > operator.yaml
$ curl https://raw.githubusercontent.com/pb82/serverless-operator/master/deploy/rbac.yaml > rbay.yaml
$ https://raw.githubusercontent.com/pb82/serverless-operator/master/deploy/crd.yaml > crd.yaml
```

Create the resources:

```sh
$ oc create -f rbac.yaml
$ oc create -f crd.yaml
$ oc create -f operator.yaml --namespace=<Namespace containing Openwhisk>
```

*NOTE*: You need to have admin permissions to create the role bindings and custom resource.

## Usage

You interact with the Operator by creating instances of the following custom resource:

```yaml
apiVersion: "serverless.pb82.com/v1alpha1"
kind: "ServerlessAction"
metadata:
  name: "<Resource name>"
spec:
  name: "<Action name>"
  kind: "<Action type, e.g. nodejs:6>"
  code: "<Code for the Action>"
  username: "<Openwhisk Username>"
  password: "<Openwhisk Password>"
  namespace: "<Openwhisk Namespace>"
```

Create a copy of the above and save it as `cr.yaml`. Make sure to fill out the fields, then run:

```sh
$ oc create -f cr.yaml
```

This will trigger the operator to create the action. In order to delete the action from Openwhisk you delete
the custom resource:

```sh
$ oc delete -f cr.yaml
```

## Development

*Prerequisite*: The Operator SDK needs to be installed on your machine. Follow [the docs](https://github.com/operator-framework/operator-sdk#quick-start) to set it up.

You can run the Operator without deploying it to your namespace by cloning this repository and running:

```sh
$ dep ensure
$ operator-sdk build docker.io/pb82/serverless-operator
$ operator-sdk up local --namespace=<Namespace containing Openwhisk>
```

