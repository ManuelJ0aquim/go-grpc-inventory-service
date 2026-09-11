# Go gRPC Inventory Service

Inventory microservice built in Go, using gRPC, Protocol Buffers, PostgreSQL, and Docker.

The service is responsible for checking product availability and performing stock reservations in a transactional way, safe against concurrency.

---

## Objective

This project simulates an **Inventory Service** within a microservices architecture.

The service provides two main operations:

- **`CheckStock`**: checks whether there is enough stock for a given product.
- **`ReserveStock`**: reserves a quantity of the product, safely decrementing the stock.

### Communication Flow

![Go gRPC Inventory Service Architecture](docs/screenshots/image.png)

---

## Architecture

The project follows a well-defined layered architecture:

```
┌─────────────────────────────┐
│       gRPC Handler          │
│   infrastructure/grpc       │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│       Service Layer         │
│       business rules        │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│      Repository Layer       │
│         PostgreSQL          │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│         PostgreSQL          │
└─────────────────────────────┘
```

### Layers

- **Domain**: contains business entities and domain errors.
- **Service**: contains the business rules related to stock control.
- **Repository**: responsible for data persistence in PostgreSQL.
- **gRPC Handler**: exposes the service operations through the gRPC API.

---

## Technologies

- **Go** 1.25
- **gRPC**
- **Protocol Buffers**
- **PostgreSQL** 15
- **Docker** & **Docker Compose**
- **database/sql** + **lib/pq**

---

## Project Structure

```text
.
├── api/
│   └── proto/
│       └── inventory/
│           └── v1/
│               └── inventory.proto
├── cmd/
│   └── server/
│       └── main.go
├── infra/
│   ├── Dockerfile
│   └── docker-compose.yml
├── internal/
│   ├── domain/
│   │   └── inventory.go
│   ├── service/
│   │   └── inventory_service.go
│   └── infrastructure/
│       ├── grpc/
│       │   ├── handler.go
│       │   └── pb/
│       ├── migrations/
│       │   └── init.sql
│       └── repository/
│           └── postgres.go
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Prerequisites

Before getting started, make sure you have the following tools installed:

- [Go](https://go.dev/)
- [Docker](https://docs.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- `protoc` (Protocol Buffers Compiler)
- [`grpcurl`](https://github.com/fullstorydev/grpcurl) (for testing the gRPC API via CLI)

---

## Running the Project

Clone the repository:

```bash
git clone https://github.com/ManuelJ0aquim/go-grpc-inventory-service.git
cd go-grpc-inventory-service
```

Start the infrastructure with Docker:

```bash
make docker-up
```

This command will start PostgreSQL on port `5432` and the Inventory Service on port `50051`.

Check the running containers:

```bash
docker ps
```

---

## Database

The application uses PostgreSQL. The initial migration creates the following table structure:

```sql
CREATE TABLE inventory (
    product_id VARCHAR(50) PRIMARY KEY,
    quantity INT NOT NULL CHECK (quantity >= 0),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Initial Seed Data

| Product | Initial Stock |
| ------- | ------------- |
| prod-1  | 100           |
| prod-2  | 50            |
| prod-3  | 0             |

To query the data directly from the database:

```bash
docker exec -it inventory_postgres \
  psql -U postgres -d inventory_db \
  -c "SELECT * FROM inventory;"
```

---

## gRPC API

The gRPC contract is defined in `api/proto/inventory/v1/inventory.proto`.

### CheckStock

Checks whether there is enough stock for the requested product.

```protobuf
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
```

### ReserveStock

Performs the reservation, decrementing the requested quantity from stock.

```protobuf
rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
```

---

## Testing the Application

You can use `grpcurl` to interact with the gRPC methods.

### 1. CheckStock — stock available

```bash
grpcurl -plaintext \
  -proto api/proto/inventory/v1/inventory.proto \
  -d '{"product_id":"prod-1","quantity":10}' \
  localhost:50051 \
  inventory.v1.InventoryService/CheckStock
```

Response:

```json
{
  "available": true,
  "currentStock": 100
}
```

### 2. CheckStock — insufficient stock

```bash
grpcurl -plaintext \
  -proto api/proto/inventory/v1/inventory.proto \
  -d '{"product_id":"prod-1","quantity":150}' \
  localhost:50051 \
  inventory.v1.InventoryService/CheckStock
```

Response:

```json
{
  "currentStock": 100
}
```

> **Note:** in Protobuf v3, boolean fields with a `false` value use the default value and can be omitted from JSON serialization — that's why `available` does not appear in the response above.

### 3. ReserveStock — reservation completed

```bash
grpcurl -plaintext \
  -proto api/proto/inventory/v1/inventory.proto \
  -d '{"order_id":"order-001","product_id":"prod-1","quantity":20}' \
  localhost:50051 \
  inventory.v1.InventoryService/ReserveStock
```

Response:

```json
{
  "success": true,
  "message": "Stock reserved successfully"
}
```

---

## Testing via Postman

In addition to testing via `grpcurl`, the service was also tested using **Postman** (with gRPC support).

### 1. CheckStock — `prod-3` with zero stock

![CheckStock prod-3](docs/screenshots/Screenshot%202026-09-11%20114502.png)

### 2. ReserveStock — `prod-3`, quantity 20 → insufficient stock

![ReserveStock prod-3 insufficient](docs/screenshots/Screenshot%202026-09-11%20114633.png)

### 3. ReserveStock — `prod-2`, quantity 50 → reservation completed successfully

![ReserveStock prod-2 success](docs/screenshots/Screenshot%202026-09-11%20114710.png)

### 4. CheckStock — `prod-1` with available stock (40 units)

![CheckStock prod-1](docs/screenshots/Screenshot%202026-09-11%20114756.png)

### 5. CheckStock — `prod-2` with updated stock after the reservation (0 units)

![CheckStock prod-2 updated](docs/screenshots/Screenshot%202026-09-11%20114816.png)

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
