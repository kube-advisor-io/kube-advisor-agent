# syntax=docker/dockerfile:1

FROM golang:1.22 as base

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./
COPY dataproviders/ dataproviders/
COPY mqtt/ mqtt/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-advisor-agent

FROM scratch

COPY --from=base /kube-advisor-agent /kube-advisor-agent

# Run
CMD ["/kube-advisor-agent"]
