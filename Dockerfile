FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app
COPY go.mod .
COPY cmd/fwdr/main.go ./cmd/fwdr/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o fwdr ./cmd/fwdr

FROM alpine:3.23
WORKDIR /app
RUN apk add --update --no-cache yq && \
    adduser -D -H fwdr
COPY --chmod=755 scripts/entrypoint.sh .
COPY --from=builder /app/fwdr .
USER fwdr
CMD ["./entrypoint.sh"]
