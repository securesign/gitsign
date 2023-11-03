# Build stage
FROM registry.access.redhat.com/ubi9/go-toolset@sha256:82d9bc5d3ceb43635288880f26207201e55d1c688a60ebbfff4f54d4963a62a1 AS build-env
WORKDIR /gitsign
RUN git config --global --add safe.directory /gitsign
COPY . .
USER root
RUN git stash && \
    export GIT_VERSION=$(git describe --tags --always --dirty) && \
    git stash pop && \
    make -f Build.mak gitsign-cli-darwin-amd64 && \
    make -f Build.mak gitsign-cli-linux-amd64 && \
    make -f Build.mak gitsign-cli-windows && \
    gzip gitsign_cli_darwin_amd64 && \
    gzip gitsign_cli_linux_amd64 && \
    gzip gitsign_cli_windows_amd64.exe && \
    ls -la

# Install Gitsign
FROM registry.access.redhat.com/ubi9/ubi-minimal@sha256:b40f52aa68b29634ff45429ee804afbaa61b33de29ae775568933c71610f07a4

LABEL description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.display-name="Gitsign container image for Red Hat Trusted Artifact Signer"
LABEL io.openshift.tags="gitsign trusted-artifact-signer"
LABEL summary="Provides the gitsign CLI binary for signing and verifying container images."
LABEL com.redhat.component="gitsign"

COPY --from=build-env /gitsign/gitsign_cli_darwin_amd64.gz /usr/local/bin/gitsign_cli_darwin_amd64.gz
COPY --from=build-env /gitsign/gitsign_cli_linux_amd64.gz /usr/local/bin/gitsign_cli_linux_amd64.gz
COPY --from=build-env /gitsign/gitsign_cli_windows_amd64.exe.gz /usr/local/bin/gitsign_cli_windows_amd64.exe.gz

ENV HOME=/home
WORKDIR ${HOME}

RUN chown root:0 /usr/local/bin/gitsign_cli_darwin_amd64.gz && chmod g+wx /usr/local/bin/gitsign_cli_darwin_amd64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_linux_amd64.gz && chmod g+wx /usr/local/bin/gitsign_cli_linux_amd64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_windows_amd64.exe.gz && chmod g+wx /usr/local/bin/gitsign_cli_windows_amd64.exe.gz && \
    chgrp -R 0 /${HOME} && chmod -R g=u /${HOME}

LABEL com.redhat.component="gitsign"
# Makes sure the container stays running
CMD ["tail", "-f", "/dev/null"]
