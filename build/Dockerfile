FROM registry.access.redhat.com/ubi8-dev-preview/ubi-minimal:8.0

ENV OPERATOR=/usr/local/bin/gitorchestrator \
    USER_UID=1001 \
    USER_NAME=gitorchestrator

# install operator binary
COPY build/_output/bin/gitorchestrator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
