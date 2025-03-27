FROM golang:1.23 AS build

WORKDIR /gokafka

COPY . .

RUN go mod download
RUN go build -o build/api cmd/api/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /api

LABEL org.opencontainers.image.source = https://github.com/xoticdsign/gokafka
LABEL org.opencontainers.image.description = "API Microservice Image \
\
API Microservice is responsible for taking and handling hypothetical post requests and producing the necessary event to kafka \
\
API Microservice consist of: \
1. API itself (handlers, server, tests etc.); \
2. Kafka producer and consumer."
LABEL org.opencontainers.image.licenses = MIT

COPY --from=build gokafka/build/api .
COPY --from=build gokafka/config ./config

EXPOSE 8080

CMD ["./api"]