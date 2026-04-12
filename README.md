# gRPC Product Management System

A modern full-stack application with Go backend using gRPC and a React frontend.

## Table of Contents

1. [Technologies](#technologies)
2. [Architecture](#architecture)
3. [Data Flow](#data-flow)
4. [Getting Started](#getting-started)
5. [Development Workflow](#development-workflow)
6. [API Reference](#api-reference)
7. [Project Structure](#project-structure)

---

## Technologies

### Backend (Go)

| Technology           | Purpose                        |
| -------------------- | ------------------------------ |
| **Go 1.26**          | Backend language               |
| **gRPC**             | High-performance RPC framework |
| **Protocol Buffers** | Interface definition language  |
| **PostgreSQL**       | Relational database            |
| **pgx**              | PostgreSQL driver for Go       |
| **godotenv**         | Environment variable loading   |

### Frontend (React)

| Technology      | Purpose                    |
| --------------- | -------------------------- |
| **React 18**    | UI library                 |
| **TypeScript**  | Type-safe JavaScript       |
| **Vite**        | Build tool and dev server  |
| **pnpm**        | Package manager            |
| **Biome**       | Linting and formatting     |
| **protobuf-ts** | TypeScript code generation |

### Infrastructure

| Technology         | Purpose                              |
| ------------------ | ------------------------------------ |
| **Docker**         | Containerization                     |
| **Docker Compose** | Multi-container orchestration        |
| **nginx**          | Reverse proxy and static file server |

---

## Architecture

```
┌────────────────────────────────────────────────────────────────────────┐
│                              HOST MACHINE                              │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │                     Docker Network                             │    │
│  │                                                                │    │
│  │   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │    │
│  │   │   postgres   │    │     api      │    │     web      │     │    │
│  │   │  (Database)  │◄──►│ (Go/gRPC)    │◄──►│ (nginx)      │     │    │
│  │   │   :5432      │    │ :50051,8080  │    │   :80        │     │    │
│  │   └──────────────┘    └──────────────┘    └──────────────┘     │    │
│  │         │                   │                   │              │    │
│  └─────────┼───────────────────┼───────────────────┼──────────────┘    │
│            │                   │                   │                   │
│     host:5432         host:8080           host:80  │                   │
│            │                   │                   │                   │
│            ▼                   ▼                   ▼                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐              │
│  │Database UI   │    │  Postman/    │    │  Browser     │              │
│  │(pgAdmin)     │    │  grpcurl     │    │  (React)     │              │
│  └──────────────┘    └──────────────┘    └──────────────┘              │
└────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Create Product Flow

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Browser │    │  nginx  │    │    Go   │    │   pgx   │    │Postgres │
└────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
     │              │              │              │              │
     │ POST /api/   │              │              │              │
     │ products     │              │              │              │
     ├─────────────►│              │              │              │
     │              │ proxy_pass   │              │              │
     │              ├─────────────►│              │              │
     │              │              │ Decode       │              │
     │              │              ├─────────────►│              │
     │              │              │ SQL INSERT   │              │
     │              │              ├──────────────────────────────►│
     │              │              │              │ RETURNING    │
     │              │              │◄──────────────────────────────┤
     │              │              │ Product      │              │
     │              │ JSON         │              │              │
     │              │◄─────────────┤              │              │
     │ JSON         │              │              │              │
     │◄─────────────┤              │              │              │
     │              │              │              │              │
```

### List Products Flow

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Browser │    │  nginx  │    │    Go   │    │   pgx   │    │Postgres │
└────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
     │              │              │              │              │
     │ GET /api/    │              │              │              │
     │ products     │              │              │              │
     ├─────────────►│              │              │              │
     │              │ proxy_pass   │              │              │
     │              ├─────────────►│              │              │
     │              │              │ SELECT *     │              │
     │              │              ├──────────────────────────────►│
     │              │              │ Cursor       │              │
     │              │              │◄──────────────────────────────┤
     │              │              │ Products     │              │
     │              │ JSON array   │ array        │              │
     │              │◄─────────────┤              │              │
     │ JSON array   │              │              │              │
     │◄─────────────┤              │              │              │
```

### CRUD Operations Summary

| Operation  | HTTP Method | Endpoint             | gRPC Method | SQL                        |
| ---------- | ----------- | -------------------- | ----------- | -------------------------- |
| **Create** | POST        | `/api/products`      | `Create`    | `INSERT ... RETURNING`     |
| **Read**   | GET         | `/api/products/{id}` | `Get`       | `SELECT ... WHERE id = $1` |
| **Update** | PUT         | `/api/products/{id}` | `Update`    | `UPDATE ... WHERE id = $1` |
| **Delete** | DELETE      | `/api/products/{id}` | `Delete`    | `DELETE ... WHERE id = $1` |
| **List**   | GET         | `/api/products`      | `List`      | `SELECT * FROM products`   |

---

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose**
- **pnpm** (optional, for local development)
- **Node.js 22+** (optional, for local development)
- **Go 1.26+** (optional, for local development)

### Quick Start (Docker)

```bash
# Clone and navigate to project
cd /gRPC

# Start all services
docker compose up --build -d

# Wait for containers to start (~10 seconds)
sleep 10

# Open browser
# http://localhost
```

### Access Services

| Service        | URL                   | Description            |
| -------------- | --------------------- | ---------------------- |
| **Web App**    | http://localhost      | React CRUD interface   |
| **API (HTTP)** | http://localhost:8080 | REST endpoints         |
| **API (gRPC)** | localhost:50051       | gRPC service (grpcurl) |
| **Database**   | localhost:5432        | PostgreSQL             |

### Useful Commands

```bash
# View running containers
docker compose ps

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Rebuild a specific service
docker compose build api
docker compose up -d api
```

---

## Development Workflow

### Option 1: Container Development (Production-like)

```bash
# Make code changes in web/ or api/
# Rebuild to see changes
docker compose build web
docker compose up -d web
```

### Option 2: Local Development (Hot Reload)

```bash
# Terminal 1: Start services except web
docker compose up -d postgres api

# Terminal 2: Run React locally
cd web
pnpm dev

# Open http://localhost:5173
```

### Option 3: Full Local Development

```bash
# Terminal 1: Start PostgreSQL
docker run -d --name grpc-postgres \
  -e POSTGRES_USER=grpcuser \
  -e POSTGRES_PASSWORD=grpcpass \
  -e POSTGRES_DB=grpcdb \
  -p 5432:5432 postgres:15-alpine

# Terminal 2: Run Go API
cd api
source .env  # or set DATABASE_URL manually
./server

# Terminal 3: Run React
cd web
pnpm dev
```

---

## API Reference

### REST Endpoints (via nginx proxy)

```bash
# List all products
curl http://localhost/api/products

# Get single product
curl http://localhost/api/products/1

# Create product
curl -X POST http://localhost/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99,"quantity":100}'

# Update product
curl -X PUT http://localhost/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated","price":19.99,"quantity":50}'

# Delete product
curl -X DELETE http://localhost/api/products/1
```

### gRPC Endpoints

```bash
# Using grpcurl
grpcurl -plaintext -d '{"name":"Widget","price":9.99,"quantity":100}' \
  localhost:50051 pb.ProductService/Create

grpcurl -plaintext -d '{}' \
  localhost:50051 pb.ProductService/List

grpcurl -plaintext -d '{"id":"1"}' \
  localhost:50051 pb.ProductService/Get
```

### gRPC Web UI

```bash
# Start gRPC UI for interactive testing
grpcui -plaintext localhost:50051

# Open browser at http://localhost:60551
```

---

## Project Structure

```
/gRPC/
│
├── README.md                 # This file
├── docker-compose.yml        # Container orchestration
├── .env                   # Environment variables
│
├── proto/
│   └── product.proto       # Protocol Buffers definition
│
├── api/                    # Go backend
│   ├── Dockerfile         # Multi-stage build
│   ├── .env             # DATABASE_URL
│   ├── go.mod           # Go dependencies
│   ├── go.sum
│   ├── main.go           # Entry point + HTTP handlers
│   ├── server           # Compiled binary (after build)
│   │
│   ├── db/
│   │   └── db.go       # PostgreSQL connection pool
│   │
│   ├── handlers/
│   │   └── product.go # gRPC CRUD handlers
│   │
│   └── pb/
│       └── product.pb.go # Generated Go from proto
│
└── web/                    # React frontend
    ├── Dockerfile         # Build + nginx
    ├── nginx.conf         # Reverse proxy config
    ├── package.json
    ├── pnpm-lock.yaml
    ├── tsconfig.json
    │   # vite.config.ts
    │   # biome.json
    │
    └── src/
        ├── main.tsx      # React entry point
        ├── App.tsx      # Main component (CRUD UI)
        ├── index.css     # Global styles
        │
        ├── proto/       # Generated TS from proto
        │   ├── product.ts
        │   ├── product.grpc-client.ts
        │   └── product.proto
        │
        └── vite-env.d.ts # Vite type definitions
```

---

## Protocol Buffers Definition

The core of this project is the `product.proto` file:

```protobuf
syntax = "proto3";
option go_package = "github.com/TiagoBehencks/gRPC/api/pb";

package pb;

// Message definitions
message Product {
    string id = 1;
    string name = 2;
    double price = 3;
    int32 quantity = 4;
}

message CreateProductRequest {
    string name = 1;
    double price = 2;
    int32 quantity = 3;
}

// ... more messages

// Service definition
service ProductService {
    rpc Create(CreateProductRequest) returns (Product);
    rpc Get(GetProductRequest) returns (Product);
    rpc Update(UpdateProductRequest) returns (Product);
    rpc Delete(DeleteProductRequest) returns (Empty);
    rpc List(ListProductsRequest) returns (ListProductsResponse);
}
```

### Generating Code

```bash
# Generate Go code (in api/ folder)
cd api
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/product.proto

# Generate TypeScript code (in web/ folder)
cd web
protoc --ts_out=client_grpc1:src \
  --proto_path src src/proto/product.proto
```

---

## Troubleshooting

### Container won't start

```bash
# Check logs
docker compose logs postgres
docker compose logs api
docker compose logs web

# Check if ports are in use
netstat -tulpn | grep -E '5432|8080|50051|80'
```

### Database connection errors

```bash
# Verify PostgreSQL is running
docker compose ps

# Check DATABASE_URL format
# Should be: postgres://user:pass@hostname:port/database
```

### Frontend shows "Failed to fetch"

```bash
# Test API directly
curl http://localhost:8080/api/products

# Test via nginx
curl http://localhost/api/products

# Check nginx logs
docker compose logs web
```

---

## Credits

Built with:

- [Go](https://go.dev/)
- [gRPC](https://grpc.io/)
- [Protocol Buffers](https://protobuf.dev/)
- [React](https://react.dev/)
- [Vite](https://vitejs.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
