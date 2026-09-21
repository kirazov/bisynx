FROM mcr.microsoft.com/devcontainers/go:1-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    make \
    build-essential \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/

USER vscode

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest