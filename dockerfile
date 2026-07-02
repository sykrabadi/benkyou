FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o benkyou ./cmd

FROM alpine:3.19
WORKDIR /app

EXPOSE 8230

COPY --from=builder /app/benkyou .
COPY --from=builder /app/data ./data
CMD ["./benkyou"]