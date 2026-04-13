# E-commerce Microservices (Go)

A high-performance E-commerce system built with **Golang**, utilizing a **Microservices** architecture, **Clean Architecture**, and **Modular Monolith** principles to ensure scalability and maintainability.

## 🚀 Technology Stack
- **Language:** Go (Golang)
- **Communication:** gRPC (Internal Services), REST API (Gateway)
- **Database:** MongoDB (Document Store)
- **Caching/Session:** Redis (OTP, Caching, Rate Limiting)
- **Message Broker:** RabbitMQ (Event-Driven Architecture)
- **Infrastructure:** Docker, Docker Compose

## 🏗️ Project Structure
```text
.
├── cmd/                # Entry points for microservices
├── internal/           # Domain logic, UseCases, Repositories
│   ├── auth/           # Auth Service (JWT, OTP)
│   ├── order/          # Order Service
│   ├── product/        # Product Service
│   └── platform/       # Shared utilities (Redis, MongoDB, gRPC)
├── pkg/                # Public/Shared packages
├── proto/              # Protocol Buffer definitions
├── deploy/             # Infrastructure (Docker-Compose, K8s)
└── .env.example        # Environment variables template
```

### Run a Service (e.g., Auth Service):
```go run cmd/auth/main.go```
