# Build stage
FROM brew.registry.redhat.io/rh-osbs/openshift-golang-builder:rhel_9_1.21 AS build-env
WORKDIR /gitsign
RUN git config --global --add safe.directory /gitsign
COPY . .
USER root
RUN git stash && \
    export GIT_VERSION=$(git describe --tags --always --dirty) && \
    git stash pop && \
    go mod download && \
    make -f Build.mak cross-platform && \
    gzip gitsign_cli_darwin_amd64 && \
    gzip gitsign_cli_linux_amd64 && \
    gzip gitsign_cli_windows_amd64.exe && \
    gzip gitsign_cli_darwin_arm64 && \
    gzip gitsign_cli_linux_arm64 && \
    gzip gitsign_cli_linux_ppc64le && \
    gzip gitsign_cli_linux_s390x && \
    ls -la

# Install Gitsign
FROM registry.access.redhat.com/ubi9/ubi-minimal@sha256:b7a3642d6245446da03d14482740be5f2fe58f30b9dfe001e89a39071a50edfc

LABEL description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.description="Gitsign is a source code signing tool that leverages simple, secure, and auditable signatures based on simple primitives and best practices."
LABEL io.k8s.display-name="Gitsign container image for Red Hat Trusted Artifact Signer"
LABEL io.openshift.tags="gitsign trusted-artifact-signer"
LABEL summary="Provides the gitsign CLI binary for signing and verifying container images."
LABEL com.redhat.component="gitsign"

COPY --from=build-env /gitsign/gitsign_cli_darwin_amd64.gz /usr/local/bin/gitsign_cli_darwin_amd64.gz
COPY --from=build-env /gitsign/gitsign_cli_linux_amd64.gz /usr/local/bin/gitsign_cli_linux_amd64.gz
COPY --from=build-env /gitsign/gitsign_cli_darwin_arm64.gz /usr/local/bin/gitsign_cli_darwin_arm64.gz
COPY --from=build-env /gitsign/gitsign_cli_linux_arm64.gz /usr/local/bin/gitsign_cli_linux_arm64.gz
COPY --from=build-env /gitsign/gitsign_cli_linux_ppc64le.gz /usr/local/bin/gitsign_cli_linux_ppc64le.gz
COPY --from=build-env /gitsign/gitsign_cli_linux_s390x.gz /usr/local/bin/gitsign_cli_linux_s390x.gz
COPY --from=build-env /gitsign/gitsign_cli_windows_amd64.exe.gz /usr/local/bin/gitsign_cli_windows_amd64.exe.gz

ENV HOME=/home
WORKDIR ${HOME}

RUN chown root:0 /usr/local/bin/gitsign_cli_darwin_amd64.gz && chmod g+wx /usr/local/bin/gitsign_cli_darwin_amd64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_linux_amd64.gz && chmod g+wx /usr/local/bin/gitsign_cli_linux_amd64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_windows_amd64.exe.gz && chmod g+wx /usr/local/bin/gitsign_cli_windows_amd64.exe.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_linux_arm64.gz && chmod g+wx /usr/local/bin/gitsign_cli_linux_arm64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_darwin_arm64.gz && chmod g+wx /usr/local/bin/gitsign_cli_darwin_arm64.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_linux_ppc64le.gz && chmod g+wx /usr/local/bin/gitsign_cli_linux_ppc64le.gz && \
    chown root:0 /usr/local/bin/gitsign_cli_linux_s390x.gz && chmod g+wx /usr/local/bin/gitsign_cli_linux_s390x.gz && \
    chgrp -R 0 /${HOME} && chmod -R g=u /${HOME}

LABEL com.redhat.component="gitsign"
# Makes sure the container stays running
CMD ["tail", "-f", "/dev/null"]