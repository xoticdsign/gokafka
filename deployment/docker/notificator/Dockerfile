FROM golang:1.23 AS build

WORKDIR /gokafka

COPY . .

RUN go mod download
RUN go build -o build/notificator cmd/notificator/main.go

FROM gcr.io/distroless/base-debian12

LABEL org.opencontainers.image.source https://github.com/xoticdsign/gokafka
LABEL org.opencontainers.image.description Notificator Microservice Image
LABEL org.opencontainers.image.licenses MIT

WORKDIR /notificator

COPY --from=build gokafka/build/notificator .
COPY --from=build gokafka/config ./config

CMD ["./notificator"]