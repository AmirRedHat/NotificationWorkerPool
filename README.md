# Notification Service

A simple and scalable notification system implemented in Go.

The system receives notification requests through an HTTP BFF, forwards them to the Notification Service using gRPC, and processes notifications asynchronously using Redis and a configurable worker pool.

## Architecture
A BFF API Gateway to receive client requests with HTTP and send the requests to a notification service via G-RPC.
Notification service push a notification job in redis (each provider have an specific queue). 
Worker service will pull the jobs and process them.

BFF --HTTP--> notification --GRPC--> REDIS <---- workers 

## Features

- Go-based backend services
- HTTP API through BFF
- gRPC communication between BFF and Notification Service
- Redis-based asynchronous job queue
- Configurable worker pool
- Graceful shutdown
- Retry mechanism with exponential backoff
- Structured logging
- Support for multiple notification queue:
    - SMS
    - Email
    - Push
- Dockerized services
- Docker Compose for local development

## Requirements

- Docker
- Docker Compose

The application does not require Go or Redis to be installed locally when running through Docker Compose.

## Running the Project

Clone the repository:

```bash
git clone <repository-url>
cd <project-directory>
```

Build and start all services:

```bash
docker compose up --build
```

Run in detached mode:

```bash
docker compose up --build -d
```

Check running containers:

```bash
docker compose ps
```

View logs:

```bash
docker compose logs -f <service-name>
```

View logs for a specific service:

```bash
docker compose logs -f worker
```

Stop the application:

```bash
docker compose down
```

When running inside Docker Compose, services should communicate using their Compose service names rather than `localhost`.

For example:

```text
REDIS_ADDR=redis:6379
```

instead of:

```text
REDIS_ADDR=localhost:6379
```

## Worker Pool

The worker service uses a configurable worker pool to process notifications concurrently.
The number of workers can be configured without changing the application code.
Each worker continuously consumes notifications from Redis and processes them independently.

## Retry

Failed notification deliveries are retried automatically.
The retry mechanism uses exponential backoff:
After the maximum number of attempts is reached, the notification is considered failed and can be handled separately, for example by a Dead Letter Queue.

## Graceful Shutdown

The services listen for termination signals and shut down gracefully.

When a shutdown signal is received:

1. The service stops accepting new work.
2. Workers stop consuming new jobs.
3. Currently processing jobs are allowed to finish.
4. The application waits for all workers to exit.
5. The process terminates cleanly.