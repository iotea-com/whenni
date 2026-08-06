# start from the latest Debian image
FROM ubuntu:22.04

ARG TARGETPLATFORM
ARG GO_VERSION=1.24.2

# Set environment variables necessary for CGO
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Prevent apt from prompting for input
ENV DEBIAN_FRONTEND=noninteractive

# Install packages and clean up in one layer to keep the image size small
RUN apt-get update && apt-get install -y --no-install-recommends \
    supervisor \
    libc6-dev \
    gcc \
    ca-certificates \
    curl  && \
    # Clean up
    apt-get clean && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

# Switch to bash for all subsequent RUN commands
SHELL ["/usr/bin/bash", "-c"]

# Install Go with the specified version based on the target platform
RUN case "${TARGETPLATFORM}" in \
    "linux/amd64") \
    curl -LO "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz \
    ;; \
    "linux/arm64") \
    curl -LO "https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz" \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz \
    && rm go${GO_VERSION}.linux-arm64.tar.gz \
    ;; \
    esac

# Set up Go environment
ENV GOPATH=/go \
    PATH=/usr/local/go/bin:/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install air
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && ./install.sh -d

WORKDIR /tmp
# Install nvm and Node.js
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash && \
    export NVM_DIR="$HOME/.nvm" && \
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && \
    nvm install v20.10.0 && \
    npm install -g npm@10.2.4 && \
    npm install -g nx pnpm && \
    corepack enable

# Repo is bind-mounted at /gruent (see deploy/engine/local/docker-compose.yaml)
WORKDIR /gruent

# Command to run the executable
CMD chmod +x ./deploy/engine/local/config/start.sh && ./deploy/engine/local/config/start.sh