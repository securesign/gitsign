# Build stage
FROM registry.access.redhat.com/ubi9/go-toolset@sha256:330c52d81d5bde432fb59c4943fcb5143940ceb460f99c1ac8e0a9ea1f8f77e8 AS build-env
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
FROM registry.access.redhat.com/ubi9/ubi-minimal@sha256:06d06f15f7b641a78f2512c8817cbecaa1bf549488e273f5ac27ff1654ed33f0

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
