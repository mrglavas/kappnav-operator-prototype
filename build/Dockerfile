FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/kappnav-operator \
    USER_UID=1001 \
    USER_NAME=kappnav-operator

# install operator binary
COPY build/_output/bin/kappnav-operator ${OPERATOR}
COPY deploy/maps/ maps/
COPY deploy/crds/extensions crds/

# get application CRD from Kubernetes Application SIG
RUN curl -fsSLO --compressed https://raw.githubusercontent.com/kubernetes-sigs/application/master/config/crds/app_v1beta1_application.yaml \
    && mv app_v1beta1_application.yaml crds/

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
