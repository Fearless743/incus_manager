FROM golang:1.21-alpine AS builder

WORKDIR /app

# Configure Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GONOSUMCHECK=*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
