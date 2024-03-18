FROM golang:1.22.0-alpine AS builder

WORKDIR /go/src/github.com/plusik10/order-info-service
COPY . .

RUN go mod download
RUN go build -o ./bin/sender cmd/sender/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/github.com/plusik10/order-info-service/bin/sender .
COPY --from=builder /go/src/github.com/plusik10/order-info-service/config/config.yml /root/config.yml
COPY --from=builder /go/src/github.com/plusik10/order-info-service/.env .

CMD ["./sender"]