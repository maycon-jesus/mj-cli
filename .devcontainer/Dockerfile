FROM mcr.microsoft.com/devcontainers/go:1.25

# Set Go cache directory
USER vscode
ENV GOCACHE=/home/vscode/.cache/go-build


RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/gopls@latest