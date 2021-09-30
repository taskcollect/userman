# https://github.com/microsoft/vscode-dev-containers/tree/v0.187.0/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.17"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
 RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
     && apt-get -y install --no-install-recommends nodejs iputils-ping 

# We need to be in the folder with the go package to run any go get's
WORKDIR /root/workspace/src
# PostgreSQL Driver
RUN go get github.com/lib/pq
# Better JSON Parser
RUN go get github.com/buger/jsonparser
# Sane argon2id library
RUN go get github.com/alexedwards/argon2id