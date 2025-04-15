# Билд стадия
FROM golang:1.24.2-alpine AS builder

WORKDIR /lead

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /lead/mylead cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /lead/mylead .

COPY --from=builder /lead/config ./config/

COPY --from=builder /lead/.env .

EXPOSE 8080

CMD ["./mylead"]