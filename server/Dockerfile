# Build stage I: 
# FROM golang:1.19 AS build
FROM mcr.microsoft.com/oss/go/microsoft/golang:1.21.8-1-bullseye AS build

WORKDIR /server
COPY . .


ARG PAT
RUN echo "export GOPRIVATE='dev.azure.com'" >> ~/.bashrc
RUN git config --global url."https://$PAT@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf "https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"


RUN go mod tidy

# The output will be in the current directory.
# The path is the complete package path. No main name in the path.
RUN go build dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/cmd/client
RUN go build dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/cmd/demoserver
RUN go build dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/cmd/server


# Build stage II: 
# FROM scratch
# FROM busybox
FROM mcr.microsoft.com/openjdk/jdk:17-ubuntu
COPY --from=build /server/client /
COPY --from=build /server/demoserver /
COPY --from=build /server/server /

CMD ["/server"]