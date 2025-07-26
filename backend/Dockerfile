# -------- Stage 1: Download dependencies --------
FROM golang:1.24.4-alpine as deps

WORKDIR /app
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

# -------- Stage 2: Build application --------
FROM golang:1.24.4-alpine as builder

# Enable CGO
ENV CGO_ENABLED=1
ENV GO111MODULE=on

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

# ✅ Compile the app
RUN go build -o bookcabin-voucher-app ./cmd/app

# -------- Stage 3: Run the binary --------
FROM alpine:latest

WORKDIR /app

# Install CA certs (for Viper or HTTP clients)
RUN apk add --no-cache sqlite-libs ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/bookcabin-voucher-app .

# Binary is already executable if built correctly, but add just in case
RUN chmod +x /app/bookcabin-voucher-app

# Copy env file and data directory (for layouts, etc.)
COPY --from=builder /app/app.env .
COPY --from=builder /app/config ./config
COPY --from=builder /app/data ./data



EXPOSE 8081

CMD ["./bookcabin-voucher-app"]