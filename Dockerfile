# syntax=docker/dockerfile:1

# build-stage
FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tasks-api

# build-release-stage
FROM gcr.io/distroless/base-debian11 AS build-release-stage

COPY --from=build-stage /tasks-api /tasks-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/tasks-api" ]

