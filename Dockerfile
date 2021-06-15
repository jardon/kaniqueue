FROM golang:1.15

WORKDIR /kaniqueue
COPY . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags '-extldflags "-static"' -o kaniqueue

FROM gcr.io/kaniko-project/executor

COPY --from=0 /kaniqueue/kaniqueue /kaniko/kaniqueue
# COPY ./config.json /kaniko/.docker/config.json
ENV HOME /root
ENV USER root
ENV PATH $PATH:/usr/local/bin:/kaniko
ENV SSL_CERT_DIR=/kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker/
ENV DOCKER_CREDENTIAL_GCR_CONFIG /kaniko/.config/gcloud/docker_credential_gcr_config.json
WORKDIR /workspace
EXPOSE 10000
RUN ["docker-credential-gcr", "config", "--token-source=env"]
ENTRYPOINT ["/kaniko/kaniqueue"]