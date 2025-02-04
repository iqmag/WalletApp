FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o WalletApp ./cmd/main.go

FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/WalletApp /WalletApp

ENTRYPOINT ["/WalletApp"]