# Build stage
FROM golang:1.20-alpine as BuildStage

WORKDIR /app

COPY src/ ./

RUN go build -o ./short-links ./cmd/server


# Deploy stage
FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /app/ ./

EXPOSE 8081

ENTRYPOINT [ "/short-links" ]
