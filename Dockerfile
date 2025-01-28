# syntax=docker/dockerfile:1

FROM golang:1.22 AS base

ARG TARGETPLATFORM

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY providers/ providers/
COPY resources/ resources/
COPY mqtt/ mqtt/
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-advisor-agent

# this image will only contain our kube-advisor-agent binary and its default config.
FROM scratch

COPY --from=base /kube-advisor-agent /kube-advisor-agent
COPY default_config.yaml /

CMD ["/kube-advisor-agent"]
