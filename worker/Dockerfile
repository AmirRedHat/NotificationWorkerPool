FROM golang:1.26 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /app/worker-bin ./worker/cmd


FROM alpine:latest

WORKDIR /app

COPY --from=build /app/worker-bin /app/worker-bin
COPY --from=build /app/worker/config /app/worker/config

ENTRYPOINT ["/app/worker-bin"]