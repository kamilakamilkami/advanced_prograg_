# E-Commerce Microservices Platform

---

## 1. Project Overview  
A production-ready, Go-based microservices architecture for an e-commerce platform. It consists of:

- **API Gateway** (`/api-gateway`, HTTP → port 8080)  
  A unified HTTP façade (Gin) forwarding client requests to underlying gRPC services.

- **Auth Service** (`/services/user-service`, gRPC → port 50053)  
  Manages user registration, JWT-based authentication, and profile retrieval.

- **Products Service** (`/services/inventory-service`, gRPC → port 5001)  
  Implements full CRUD for product catalog, with pagination, filtering, Redis caching, and stock updates via NATS events.

- **Orders Service** (`/services/order-service`, gRPC → port 5002)  
  Handles order placement, status transitions, per-user order history, and publishes inventory events on NATS.

**Key focus**: gRPC API design, event-driven messaging (NATS), per-service caching (Redis), and isolated persistence (MongoDB).

---

## 2. Technologies Used  
- **Language & Frameworks**  
  - Go 1.19+  
  - Gin (HTTP API Gateway)  
  - Protocol Buffers & gRPC  
- **Messaging & Events**  
  - NATS (publish/subscribe)  
- **Databases & Caching**  
  - MongoDB (each service uses its own database/collection)  
  - Redis (in-service cache)  
- **Infrastructure & Tooling**  
  - Docker (MongoDB, Redis, NATS containers)  
  - Go Modules  
  - Postman or curl for HTTP testing  

---

## 3. Prerequisites  
- Go 1.19 or newer  
- Docker Engine  

---

## 4. Setup & Installation  

```bash
# 1) Clone the repo  
git clone https://github.com/kamilakamilkami/advanced_prograg_  
cd advanced_prograg_

# 2) Start infra services  
docker run -d --name mongo  -p 27017:27017 mongo:latest  
docker run -d --name redis  -p 6379:6379   redis:latest  
docker run -d --name nats   -p 4222:4222   nats:latest  

# 3) Fetch Go dependencies  
for svc in services/user-service \
           services/inventory-service \
           services/order-service \
           api-gateway; do
  cd "$svc" && go mod download && cd - || exit
done
```
## 5. Running the Platform

### Auth Service (gRPC → 50053)
cd services/user-service && go run main.go

### Products Service (gRPC → 5001)
cd services/inventory-service && go run main.go

### Orders Service (gRPC → 5002)
cd services/order-service && go run main.go

### API Gateway (HTTP → 8080)
cd api-gateway && go run main.go

## 6. Running Tests
```bash
cd services/user-service      && go test ./...
cd services/inventory-service && go test ./...
cd services/order-service     && go test ./...
```
 - **Unit tests cover business logic (use-cases).**
- **Integration tests verify MongoDB, Redis, and NATS interactions.**

## 7.HTTP Endpoints
### Auth Service 
```go
func RegisterUserRoutes(r *gin.Engine) {
  r.POST("/register",       services.RegisterUser)      // Create account
  r.POST("/authenticate",   services.AuthenticateUser)  // Login & issue JWT
  r.GET("/profile/:user_id",services.GetUserProfile)    // Fetch profile
}

```

### Products Service 
```go
func RegisterInventoryRoutes(r *gin.Engine) {
  r.GET("/products",        services.GetProducts)       // List products
  r.GET("/products/:id",    services.GetProductByID)    // Get product by ID
  r.POST("/products",       services.CreateProduct)     // Add product
  r.PATCH("/products/:id",  services.UpdateProduct)     // Update product
  r.DELETE("/products/:id", services.DeleteProduct)     // Delete product
}

```

### Orders Service
```go
func RegisterOrderRoutes(r *gin.Engine) {
  r.POST("/orders",         services.CreateOrder)       // Place order
  r.GET("/orders/:id",      services.GetOrderByID)      // Get order details
  r.PATCH("/orders/:id",    services.UpdateOrderStatus) // Change order status
  r.GET("/orders",          services.GetUserOrders)     // List user’s orders
}
```
## 8. gRPC API Definitions
### Auth Service (port 50053)
```protobuf
service UserService {
  rpc Register       (RegisterRequest)   returns (RegisterResponse);
  rpc Authenticate   (AuthRequest)       returns (AuthResponse);
  rpc GetUserProfile (UserID)            returns (UserProfile);
}

```

### Products Service (port 5001)
```protobuf
service InventoryService {
  rpc CreateProduct      (Product)             returns (ProductID);
  rpc ListProducts       (ListProductsRequest) returns (ListProductsResponse);
  rpc GetProductByID     (ProductID)           returns (Product);
  rpc UpdateProduct      (UpdateProductRequest) returns (google.protobuf.Empty);
  rpc DeleteProduct      (ProductID)           returns (google.protobuf.Empty);
}

```

### Orders Service (port 5002)
```protobuf
service OrderService {
  rpc CreateOrder        (OrderRequest)           returns (OrderResponse);
  rpc GetOrderByID       (OrderID)                returns (Order);
  rpc UpdateOrderStatus  (UpdateStatusRequest)    returns (google.protobuf.Empty);
  rpc GetOrdersByUser    (GetOrdersByUserRequest) returns (GetOrdersByUserResponse);
}
```
## 🔑 9.Key Features

### 👤 User Management
- Secure user registration with password hashing
- JWT-based authentication for stateless session management
- User profile retrieval via unique user ID

### 📦 Product Catalog
- Full CRUD operations implemented over gRPC
- Advanced filtering and pagination for efficient product listings
- Redis caching for optimized access to frequently requested products
- Real-time stock adjustment via NATS event-driven architecture

### 🛒 Order Processing
- Order creation and validation using structured protobuf messages
- Event-based inventory synchronization using NATS
- Real-time order status updates with gRPC PATCH
- Authenticated user order history with filtering support

### ⚙️ Infrastructure & Operations
- Containerized architecture using Docker for MongoDB, Redis, and NATS
- Per-service MongoDB database with isolated models
- Redis for low-latency caching and performance boost
- Centralized API Gateway using Gin to route all HTTP requests

### ✅ Testing & Quality Assurance
- Unit and integration tests for all major services and data flows
- End-to-end flow validation using Postman and curl
- Stable, testable microservice environment with green build status


