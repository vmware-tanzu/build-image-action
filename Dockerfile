FROM --platform=${BUILDPLATFORM} golang:1.20 AS build-stage0

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

ARG VERSION=unknown

WORKDIR /root/
COPY . ./

RUN go build -trimpath -ldflags "-X github.com/vmware-tanzu/build-image-action/pkg/version.Version=$VERSION" -o builder main.go

FROM --platform=${BUILDPLATFORM} ubuntu:24.04

COPY --from=build-stage0 /root/builder /usr/bin/builder
COPY github-actions-entrypoint.sh /usr/bin/github-actions-entrypoint.sh

ENTRYPOINT [ "/usr/bin/github-actions-entrypoint.sh" ]
