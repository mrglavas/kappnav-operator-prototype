# kappnav-operator-prototype 

## Building the kappnav operator

This project was developed with Operator SDK v0.10.0.

Installation instructions for the Operator SDK CLI are here:
https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md

Remember to replace the RELEASE_VERSION in the instructions with v0.10.0. It may not compile with newer versions.
Prerequisites for the Operator SDK are listed here: https://github.com/operator-framework/operator-sdk

If you ever change the structs for the Kappnav CRD (located in kappnav_types.go) be sure to run:

1. operator-sdk generate k8s
2. operator-sdk generate openapi

This regenerates the CRD and the code that allows a Kappnav CR to be accessed programatically through the k8s APIs.

To build the project run: `./build.sh`

## Installing the kappnav operator

To install the operator on OKD run:

1. kubectl create namespace kappnav
2. kubectl create -f kappnav.yaml -n kappnav

or for finer grained control on resource creation run:

1. kubectl create -f deploy/crds/kappnav_v1_kappnav_crd.yaml
2. kubectl create -f deploy/crds/kappnav_v1_kappnav_cr.yaml
3. kubectl create -f deploy/service_account.yaml
4. kubectl create -f deploy/role.yaml
5. kubectl create -f deploy/role_binding.yaml
6. kubectl create -f deploy/operator.yaml

Note: You should modify `kappnav_v1_kappnav_cr.yaml` for your Kube environment.

## Uninstalling the kappnav operator

To uninstall the operator on OKD run:

1. kubectl delete -f kappnav-delete-CR.yaml -n kappnav --now
2. kubectl delete -f kappnav-delete.yaml -n kappnav
3. kubectl delete namespace kappnav

or for finer grained control on resource deletion run:

1. kubectl delete -f deploy/crds/kappnav_v1_kappnav_crd.yaml
2. kubectl delete -f deploy/crds/kappnav_v1_kappnav_cr.yaml
3. kubectl delete -f deploy/service_account.yaml
4. kubectl delete -f deploy/role.yaml
5. kubectl delete -f deploy/role_binding.yaml
6. kubectl delete -f deploy/operator.yaml

## Default values

Default values for the operator's configuration are stored in `deploy/default_values.yaml`. This CR file is included in the Docker image and is read each time a Kappnav CR is reconciled by the operator to fill in defaults for values that were not specified in the CR.

## Adding additional CRDs to the operator

Additional CRDs can be added to the `deploy/crds/extensions` folder. These will be included in the Docker image. The Application CRD is always included in the image. When the operator is installed it will attempt to create each of the CRDs in k8s if they do not already exist.

## Adding additional action, sections and status config maps to the operator

Additional action, sections and status config maps should be added to the `deploy/maps/action`, `deploy/maps/sections` and `deploy/maps/status` folders respectively. This supports the same templating language that is used in Helm charts. Variables are addressed by their field names in the Kappnav structs. For instance, the kubeEnv field from the CR would be addressed as `.Spec.Env.KubeEnv`. Action, sections and status config maps will be initially created when a CR is installed.

## Adding additional logic to the controller

If you are adding additional logic for managing resources, provide an implemenation of the `NewKappnavExtension` function in the `utils/extensions.go` file that returns an instance of `KappnavExtension`.
