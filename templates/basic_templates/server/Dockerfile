# Build stage I: 
# FROM golang:1.19 AS build
FROM mcr.microsoft.com/oss/go/microsoft/golang:1.23.2-1-bullseye AS build

WORKDIR /server
COPY . .

<<if contains .envInformation.goModuleNamePrefix "dev.azure.com">>
ARG PAT
RUN echo "export GOPRIVATE='dev.azure.com'" >> ~/.bashrc
RUN git config --global url."https://$PAT@<<.envInformation.goModuleNamePrefix|trimGitSuffix>>".insteadOf "https://<<.envInformation.goModuleNamePrefix|trimGitSuffix>>"
<<else>>
ARG AKS_GOPROXY_TOKEN
ARG GOPROXY
ARG GOPRIVATE
ARG GONOPROXY
<<end>>

RUN go mod tidy

# The output will be in the current directory.
# The path is the complete package path. No main name in the path.
RUN go build <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/client
RUN go build <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/demoserver
RUN go build <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/server


# Build stage II: 
# FROM scratch
# FROM busybox
FROM mcr.microsoft.com/openjdk/jdk:17-ubuntu
COPY --from=build /server/client /
COPY --from=build /server/demoserver /
COPY --from=build /server/server /

CMD ["/server"]
