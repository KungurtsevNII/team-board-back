FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o teamboard ./cmd/teamboard
FROM alpine
WORKDIR /build
COPY --from=builder /build/teamboard /build/teamboard
COPY --from=builder /build/config /build/config
COPY --from=builder /build/migrations /build/migrations
CMD ["./teamboard", "--config", "config/prod.yaml"]