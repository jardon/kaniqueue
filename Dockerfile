FROM golang:1.15

WORKDIR /kaniqueue
COPY . .

RUN go build -o kaniqueue

FROM gcr.io/kaniko-project/executor

FROM ubuntu:20.04

COPY --from=0 /kaniqueue /kaniqueue
COPY --from=1 /kaniko/executor /kaniko/executor
COPY --from=1 /kaniko/docker-credential-gcr /kaniko/docker-credential-gcr
COPY --from=1 /kaniko/docker-credential-ecr-login /kaniko/docker-credential-ecr-login
COPY --from=1 /kaniko/docker-credential-acr /kaniko/docker-credential-acr
COPY --from=1 /kaniko/ssl/certs/ /kaniko/ssl/certs/
COPY --from=1 /kaniko/.docker /kaniko/.docker
COPY --from=1 /etc/nsswitch.conf /etc/nsswitch.conf
ENV HOME /root
ENV USER root
ENV PATH /usr/local/bin:/kaniko
ENV SSL_CERT_DIR=/kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker/
ENV DOCKER_CREDENTIAL_GCR_CONFIG /kaniko/.config/gcloud/docker_credential_gcr_config.json
WORKDIR /workspace
RUN ["docker-credential-gcr", "config", "--token-source=env"]
ENTRYPOINT ["/kaniqueue/kaniqueue"]