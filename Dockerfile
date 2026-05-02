FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /vidarr .

FROM alpine:3.19

RUN apk add --no-cache ffmpeg ca-certificates tzdata

COPY --from=builder /vidarr /usr/local/bin/vidarr

EXPOSE 8080

ENTRYPOINT ["vidarr"]
CMD ["--config", "/etc/vidarr/vidarr.yaml", "--port", "8080"]
