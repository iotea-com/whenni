FROM debian:stable-slim AS builder

# -----------------------------------------------------------------------------------------------------
#                                                                                      Build Arguements
# ---------------------------------------------------------------------------------------------------*/
ARG GO_VERSION=1.24.2
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# -----------------------------------------------------------------------------------------------------
#                                                                                         Prerequisites
# ---------------------------------------------------------------------------------------------------*/
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    libc6-dev \
    curl \
    git \
    gcc

# -----------------------------------------------------------------------------------------------------
#                                                                                                    Go
# ---------------------------------------------------------------------------------------------------*/
ENV CGO_ENABLED=1
ENV GO111MODULE=on

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

# -----------------------------------------------------------------------------------------------------
#                                                                                            Node, Yarn
# ---------------------------------------------------------------------------------------------------*/

# Install dependencies required for nvm and Node.js
RUN apt-get install -y zsh && chsh -s $(which zsh)

# Set the SHELL to use zsh for RUN commands
SHELL ["/usr/bin/zsh", "-c"]

# Install nvm, Node.js, and yarn globally
ENV NVM_DIR=/root/.nvm
RUN touch ~/.zshrc && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash && \
    echo 'source $NVM_DIR/nvm.sh' >> ~/.zshrc && \
    source ~/.zshrc && \
    nvm install v20.10.0 && \
    npm install -g npm@10.2.4 yarn corepack

# -----------------------------------------------------------------------------------------------------
#                                                                                                 Build
# ---------------------------------------------------------------------------------------------------*/
WORKDIR /iotea

RUN source ~/.zshrc && \
    yarn add prisma

# First we only copy and generate the Go prisma client once, 
# so we can use docker caching and speed up development builds
COPY ./prisma/schema.prisma ./prisma/schema.prisma
COPY ./prisma/gen-client.sh ./prisma/gen-client.sh
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# Create prisma client
WORKDIR /iotea/prisma
RUN source ~/.zshrc && \
    ./gen-client.sh

# Download dependencies first (this layer will be cached)
WORKDIR /iotea
RUN go mod download

# Copy the rest of the project
COPY ./libs/engine/dependencies ./libs/engine/dependencies
COPY ./libs/engine/channels ./libs/engine/channels
COPY ./libs/protocols ./libs/protocols
COPY ./libs/http ./libs/http
COPY ./libs/val ./libs/val
COPY ./libs/id ./libs/id
COPY ./libs/secrets ./libs/secrets
COPY ./services/http-api ./services/http-api

# Go to the http-api service directory
WORKDIR /iotea/services/http-api

# Build-time variables
ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GIT_DIRTY
ARG BUILD_CREATOR

# Create a checksum of the source code to determine if we need to update build info
RUN find . -type f -name "*.go" -exec sha256sum {} \; | sort | sha256sum > /tmp/source_checksum

# Get current build time in Mountain Time and update build info only if source changed
RUN BUILD_TIME=$(TZ='America/Denver' date '+%Y-%m-%d_%I:%M:%S%p_MST') && \
    GIT_BRANCH_ESCAPED=$(echo "${GIT_BRANCH}" | sed 's/[&/]/\\&/g') && \
    GIT_COMMIT_ESCAPED=$(echo "${GIT_COMMIT}" | sed 's/[&/]/\\&/g') && \
    GIT_DIRTY_ESCAPED=$(echo "${GIT_DIRTY}" | sed 's/[&/]/\\&/g') && \
    BUILD_CREATOR_ESCAPED=$(echo "${BUILD_CREATOR}" | sed 's/[&/]/\\&/g') && \
    sed -i "s/{{BUILD_VERSION}}/${GIT_BRANCH_ESCAPED}/g" config/config.go && \
    sed -i "s/{{BUILD_COMMIT}}/${GIT_COMMIT_ESCAPED}/g" config/config.go && \
    sed -i "s/{{BUILD_DIRTY}}/${GIT_DIRTY_ESCAPED}/g" config/config.go && \
    sed -i "s/{{BUILD_TIME}}/${BUILD_TIME}/g" config/config.go && \
    sed -i "s/{{BUILD_CREATOR}}/${BUILD_CREATOR_ESCAPED}/g" config/config.go

# Build the application
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o http-api .

# -----------------------------------------------------------------------------------------------------
#                                                                                              Run-time
# ---------------------------------------------------------------------------------------------------*/
FROM debian:stable-slim

# Copy build arguments to runtime stage
ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GIT_DIRTY
ARG BUILD_CREATOR
ARG GO_VERSION

# Add labels to the final image
LABEL branch=${GIT_BRANCH} \
    commit=${GIT_COMMIT} \
    dirty=${GIT_DIRTY} \
    build-creator=${BUILD_CREATOR} \
    go-version=${GO_VERSION}

# Create a non-root user with a specific UID/GID
RUN groupadd -r iotea --gid=1000 && \
    useradd -r -g iotea --uid=1000 --create-home --shell=/sbin/nologin iotea && \
    # Add docker group with GID 999 (common docker group GID)
    groupadd -r docker --gid=999 && \
    usermod -aG docker iotea

# Setup system and clean up in the same layer to keep image size small
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    dumb-init \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create startup script
# We use dumb-init here because:
# 1. Our API uses Fiber's prefork mode which doesn't work when the app has PID 1 in containers
# 2. dumb-init takes PID 1 and properly handles signal forwarding to child processes
# 3. This ensures both proper prefork operation and graceful shutdown when k8s sends termination signals

RUN echo '#!/bin/sh\n\
    # Add any startup tasks here\n\
    exec /usr/bin/dumb-init -- /http-api "$@"' > /start.sh && \
    chmod +x /start.sh

# Copy binary
COPY --from=builder /iotea/services/http-api/http-api /http-api

# Set proper permissions
RUN chown iotea:iotea /http-api /start.sh && \
    chmod 755 /http-api

# Switch to non-root user
USER iotea

# Start
WORKDIR /home/iotea
ENTRYPOINT ["/start.sh"]
