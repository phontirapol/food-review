# =====================  Build Stage =====================
FROM golang:1.18-alpine3.16 AS builder

# Important:
#   Because this is a CGO enabled package, you are required to set it as 1.
ENV CGO_ENABLED=1 GOOS=linux
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

# =====================  Main Stage =====================
FROM alpine:3.16

COPY /cmd/ .
COPY /pkg/ .
COPY --from=builder /app/ .

EXPOSE 5555
CMD ["./main"]