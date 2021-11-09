FROM golang:1.17-alpine

WORKDIR /usr/local/go/src/short-links

COPY src/ ./

RUN go mod download
RUN go build -o /short-links

EXPOSE 8081
ENTRYPOINT [ "/short-links" ]
