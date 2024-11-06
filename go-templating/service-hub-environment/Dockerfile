# Build stage I: 
# FROM golang:1.19 AS build
FROM mcr.microsoft.com/oss/go/microsoft/golang:1.23.1-1-bullseye AS build

WORKDIR /go-templating

COPY . .

# The output will be in the current directory.
RUN go build -o maintain_service go-templating/maintain-service/cmd/maintain-service

# Build stage II: 
FROM mcr.microsoft.com/openjdk/jdk:17-ubuntu
COPY --from=build /go-templating/maintain_service /app/maintain_service
COPY --from=build /go-templating/maintain-service/internal/resourcelinks/templates /app/templates


# Install required packages
RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    wget \
    default-jre \
    python3-pip \
    npm \
    curl

# Download and install kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

# Download and install Helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 \
    && chmod +x get_helm.sh \
    && ./get_helm.sh

#Download both platform's versions of golang
RUN wget https://golang.org/dl/go1.23.1.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.23.1.linux-arm64.tar.gz && \
    mv /usr/local/go /usr/local/go-arm64

RUN wget https://golang.org/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz && \
    mv /usr/local/go /usr/local/go-amd64

WORKDIR /app

ARG TARGETPLATFORM

#Copy over version of go depending on platform
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
    cp -r /usr/local/go-amd64 /usr/local/go; \
    else \
    cp -r /usr/local/go-arm64 /usr/local/go; \
    fi && \
    rm -rf /usr/local/go-amd64 && \
    rm -rf /usr/local/go-arm64

#Set required environment variables
ENV GO111MODULE=on
ENV GOBIN=/usr/local/go/bin
ENV GOPATH=/usr/local/go
ENV PATH=$PATH:/usr/local/go/bin

#Install necessary go packages.
RUN go install github.com/bufbuild/buf/cmd/buf@v1.29.0
RUN go install go.uber.org/mock/mockgen@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install npm packages
RUN npm install -g swagger-nodegen-cli

# Start a shell session (optional)
CMD ["/bin/bash"]