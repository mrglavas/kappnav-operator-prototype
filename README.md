# kappnav-operator-prototype

## Pulling from GitHub

After extracting the project, rename the base dir from `kappnav-operator-prototype` to `kappnav-operator`.

Also run the following commands:

1. chmod +x ./build/bin/entrypoint
2. chmod +x ./build/bin/user_setup

These are temporary workarounds.

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

To build the project run:

operator-sdk build kappnav.io/kappnav-operator:0.0.1

## Installing the kappnav operator

To install the operator run:

1. kubectl create -f deploy/crds/kappnav_v1_kappnav_crd.yaml
2. kubectl create -f deploy/crds/kappnav_v1_kappnav_cr.yaml
3. kubectl create -f deploy/service_account.yaml
4. kubectl create -f deploy/role.yaml
5. kubectl create -f deploy/role_binding.yaml
6. kubectl create -f deploy/operator.yaml

Note: You should modify `kappnav_v1_kappnav_cr.yaml` for your Kube environment.

## Uninstalling the kappnav operator

To uninstall the operator run:

1. kubectl delete -f deploy/crds/kappnav_v1_kappnav_crd.yaml
2. kubectl delete -f deploy/crds/kappnav_v1_kappnav_cr.yaml
3. kubectl delete -f deploy/service_account.yaml
4. kubectl delete -f deploy/role.yaml
5. kubectl delete -f deploy/role_binding.yaml
6. kubectl delete -f deploy/operator.yaml

## Adding additional CRDs to the operator

Additional CRDs can be added to the `deploy/crds/extensions` folder. These will be included in the Docker image. The Application CRD is always included in the image. When the operator is installed it will attempt to create each of the CRDs in k8s if they do not already exist.

## Adding additional status and action config maps to the operator

Additional action and status config maps should be added to the `deploy/maps/action` and `deploy/maps/status` folder. This supports the same templating language that is used in Helm charts. Variables are addressed by their field names in the Kappnav structs. For instance, the kubeEnv field from the CR would be addressed as `.Spec.Env.KubeEnv`. Action and status config maps will be initially created when a CR is installed.

## Adding additional logic to the controller

If you are adding additional logic for managing resources, provide an implemenation of the `NewKappnavExtension` function in the `utils/extensions.go` file that returns an instance of `KappnavExtension`.
