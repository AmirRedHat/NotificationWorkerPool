FROM golang:1.26 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /notification ./notification/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=build /notification ./notification

ENTRYPOINT ["./notification"]