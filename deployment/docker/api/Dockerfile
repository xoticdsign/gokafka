FROM golang:1.23 AS build

WORKDIR /gokafka

COPY . .

RUN go mod download
RUN go build -o build/api cmd/api/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /api

COPY --from=build gokafka/build/api .
COPY --from=build gokafka/config ./config

EXPOSE 8080

CMD ["./api"]