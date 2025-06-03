FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY order_service/go.mod order_service/go.sum ./
RUN go mod download

COPY order_service/cmd ./cmd
COPY order_service/internal ./internal
COPY order_service/pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest


COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]