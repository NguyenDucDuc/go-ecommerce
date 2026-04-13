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
├── api-gateway/            # API Gateway service
├── common/                 # Shared utilities, models, and constants
├── notification-service/   # Notification handling service
├── order-service/          # Order processing and management
├── product-service/        # Product catalog and inventory
├── proto/                  # Protocol Buffer definitions
├── user-service/           # Authentication and User management
├── .env                    # Environment variables
├── Makefile                # Build and run commands
└── README.md               # Project documentation
```

### Run a Service (e.g., Auth Service):
```go run cmd/auth/main.go```
