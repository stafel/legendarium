# two stage build

FROM docker.io/golang:1.21.3 AS build

WORKDIR /app

COPY app .

# cgo is needed for sqlite
RUN go mod download && CGO_ENABLED=1 GOOS=linux go build -o backend

# we need a base image for dynamic linked libs and can't work from scratch
FROM docker.io/debian:bookworm-slim

EXPOSE 9000

WORKDIR /app

COPY --from=build /app/backend backend

ENTRYPOINT ["/app/backend"]