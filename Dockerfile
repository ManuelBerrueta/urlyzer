# syntax=docker/dockerfile:1

# Multi Stage for a lean image

# Build App
FROM golang:1.20-alpine3.19 AS build-task

WORKDIR /app

COPY *.go ./

# Build
RUN go env -w GO111MODULE=auto
RUN go build -v -o /urlyzer


# Lean Container Image
FROM gcr.io/distroless/base-debian12 AS release-task

WORKDIR /

COPY --from=build-task /urlyzer /urlyzer

USER nonroot:nonroot

ENTRYPOINT ["/urlyzer"]
