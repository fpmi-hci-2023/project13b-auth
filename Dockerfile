ARG GO_VERSION=1.21.3

################
# BUILD target #
################
FROM golang:${GO_VERSION}-alpine AS build
COPY .. /build
WORKDIR /build
ARG GITHUB_TOKEN
RUN if [ ! -d "vendor" ]; then apk add git openssh; git config --global url."https://git:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; GOPRIVATE=github.com/fpmi-hci-2023/* go mod vendor; fi

ARG TARGETOS TARGETARCH commit version
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /bin/hci-auth -mod=vendor -ldflags "-X main.Commit=$commit -X main.Version=$version"

##############
# DEV target #
##############
FROM alpine:latest AS dev

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 50320
ENTRYPOINT ["/bin/hci-auth"]

COPY charts /helm-charts
COPY --from=build /bin/hci-auth /bin/hci-auth
