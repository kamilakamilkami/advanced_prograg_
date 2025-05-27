# Final Exam

## Project Overview
This is a microservices-based e-commerce platform with three core services:
- **Auth Service**: Handles user authentication and authorization
- **Products Service**: Manages product inventory
- **Orders Service**: Processes customer orders 

## 🛠 Technologies Used
- **Backend**: 
  - gRPC (primary communication between services)
  - Golang
  - NATS
- **Database**: 
  - 
- **Caching**:
  - Redis
- **Other**:
  - Docker 
  - Postman

## How to Run Locally

1. **Clone the repository**:
   ```bash
   git clone https://github.com/kamilakamilkami/advanced_prograg_.git
   ```
3. **Start services**:
   ```bash
   # In separate terminals
   cd services/user-service && go run main.go
   cd services/inventory-service && go run main.go
   cd services/order-service && go run main.go
   cd api-gateway && go run main.go
   ```


## Endpoints

### Auth Service 
```go
func RegisterUserRoutes(r *gin.Engine) {
    r.POST("/register", services.RegisterUser)          // Create new account
    r.POST("/authenticate", services.AuthenticateUser) // Login with credentials
    r.GET("/profile/:user_id", services.GetUserProfile) // Get user details
}
```

### Products Service 
```go
func RegisterInventoryRoutes(r *gin.Engine) {
    r.GET("/products", services.GetProducts)            // List all products 
    r.GET("/products/:id", services.GetProductByID)     // Get specific product
    r.POST("/products", services.CreateProduct)         // Add new product 
    r.PATCH("/products/:id", services.UpdateProduct)   // Modify product 
    r.DELETE("/products/:id", services.DeleteProduct)   // Remove product 
}
```

### Orders Service
```go
func RegisterOrderRoutes(r *gin.Engine) {
    r.POST("/orders", services.CreateOrder)            // Submit new order
    r.GET("/orders/:id", services.GetOrderByID)        // Retrieve order details
    r.PATCH("/orders/:id", services.UpdateOrderStatus) // Update status 
    r.GET("/orders", services.GetUserOrders)           // List user's orders
}
```
