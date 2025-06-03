FROM golang:1.24.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app /app/

CMD ["/app/app"]