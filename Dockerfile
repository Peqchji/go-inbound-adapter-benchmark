FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ${SERVICE_PATH}

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /server .

# Install curl for health check (optional but good)
RUN apk --no-cache add curl

CMD ["./server"]
