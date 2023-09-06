# Build stage
FROM registry.access.redhat.com/ubi9/go-toolset@sha256:52ab391730a63945f61d93e8c913db4cc7a96f200de909cd525e2632055d9fa6 AS build-env
WORKDIR /gitsign
RUN git config --global --add safe.directory /gitsign
COPY . .
USER root
RUN make build-gitsign

# Install Gitsign
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build-env /gitsign/gitsign /usr/local/bin/gitsign
RUN chown root:0 /usr/local/bin/gitsign && chmod g+wx /usr/local/bin/gitsign

# Configure home directory
ENV HOME=/home
RUN chgrp -R 0 /${HOME} && chmod -R g=u /${HOME}

WORKDIR ${HOME}

# Makes sure the container stays running
CMD ["tail", "-f", "/dev/null"]
