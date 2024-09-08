# syntax=docker/dockerfile:1

FROM golang:1.22 AS base

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY dataproviders/ dataproviders/
COPY mqtt/ mqtt/
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-advisor-agent


FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch

COPY --from=base /kube-advisor-agent /kube-advisor-agent
COPY default_config.yaml /

CMD ["/kube-advisor-agent"]
