# kappnav-operator-prototype
## Building the kappnav operator

This project was developed with Operator SDK v0.10.0.

Installation instructions for the Operator SDK CLI are here:
https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md

Remember to replace the RELEASE_VERSION in the instructions with v0.10.0. It may not compile with newer versions.
Prerequisites for the Operator SDK are listed here: https://github.com/operator-framework/operator-sdk

If you ever change the structs for the Kappnav CRD (located in kappnav_types.go) be sure to run:

`operator-sdk generate k8s
operator-sdk generate openapi`

This regenerates the CRD and the code that allows a Kappnav CR to be accessed programatically through the k8s APIs.

To build the project run:

`operator-sdk build kappnav.io/kappnav-operator:0.0.1`

## Installing the kappnav operator

To install the operator run:

`kubectl create -f deploy/crds/kappnav_v1_kappnav_crd.yaml
kubectl create -f deploy/crds/kappnav_v1_kappnav_cr.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml`

## Uninstalling the kappnav operator

To uninstall the operator run:

`kubectl delete -f deploy/crds/kappnav_v1_kappnav_crd.yaml
kubectl delete -f deploy/crds/kappnav_v1_kappnav_cr.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/operator.yaml`

## Adding additional CRDs to the operator

Additional CRDs can be added to the `deploy/crds/extensions` folder. These will be included in the Docker image. The Application CRD is always included in the image. When the operator is installed it will attempt to create each of the CRDs in k8s if they do not already exist.

## Adding additional status and action config maps to the operator

Additional action and status config maps should be added to the `deploy/maps/action` and `deploy/maps/status` folder. This supports the same templating language that is used in Helm charts. Variables are addressed by their field names in the Kappnav structs. For instance, the kubeEnv field from the CR would be addressed as `.Spec.Env.KubeEnv`
