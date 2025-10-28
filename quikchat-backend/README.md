# QuikChat Backend

This is the backend service for QuikChat, a scalable, real-time messaging platform.

## Overview

- **Language:** Go
- **Architecture:** Clean Architecture
- **Database:** PostgreSQL
- **Cache/Broker:** Redis
- **API:** RESTful HTTP and WebSockets

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running Locally

1.  Create a `.env` file from the `.env.example`.
2.  Run the services using Docker Compose:
    ```sh
    docker-compose up --build
    ```

The API will be available at `http://localhost:8080`.

