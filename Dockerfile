FROM golang

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set working directory
WORKDIR /form3_rest_api_client

# Run tests
CMD CGO_ENABLED=0 go test ./...