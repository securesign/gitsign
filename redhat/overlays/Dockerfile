# Build stage
FROM registry.access.redhat.com/ubi9/go-toolset@sha256:52ab391730a63945f61d93e8c913db4cc7a96f200de909cd525e2632055d9fa6 AS build-env
WORKDIR /gitsign
RUN git config --global --add safe.directory /gitsign
COPY . .
USER root
RUN make build-gitsign

# Install Gitsign
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.display-name="Gitsign container image for Red Hat Trusted Artifact Signer"
LABEL io.openshift.tags="gitsign trusted-artifact-signer"
LABEL summary="Provides the gitsign CLI binary for signing and verifying container images."
LABEL com.redhat.component="gitsign"

COPY --from=build-env /gitsign/gitsign /usr/local/bin/gitsign
RUN chown root:0 /usr/local/bin/gitsign && chmod g+wx /usr/local/bin/gitsign

# Configure home directory
ENV HOME=/home
RUN chgrp -R 0 /${HOME} && chmod -R g=u /${HOME}

WORKDIR ${HOME}

LABEL com.redhat.component="gitsign"
# Makes sure the container stays running
CMD ["tail", "-f", "/dev/null"]
