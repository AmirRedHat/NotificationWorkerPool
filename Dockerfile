FROM golang:1.25-alpine as Builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY worker .
ARG SERVICE

CMD ["sh", "-c", "go run ./${SERVICE}/cmd"]