FROM golang:1.23-alpine AS base

RUN --mount=type=cache,target=/root/.cache/go-build go install github.com/go-task/task/v3/cmd/task@latest

FROM base AS build

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=$GOPATH/pkg go mod download
RUN --mount=type=cache,target=/root/.cache/go-build task build BIN=./bin/app

FROM alpine:3.20
WORKDIR /server
COPY --from=build /app/bin/app .
CMD ["/server/app"]

