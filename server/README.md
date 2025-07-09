# Online Auction System – Server Directory Documentation

## Overview

The `server` directory contains the backend code for the Online Web-based Auction System (OWAS). It is responsible for handling API requests, business logic, database interactions, authentication, real-time updates, and integrations (e.g., Stripe payments, image uploads).

---

## Directory Structure

```
server/
  ├── cmd/
  │   └── api/
  │       └── main.go
  ├── internal/
  │   ├── auth/
  │   ├── config/
  │   ├── db/
  │   ├── handlers/
  │   ├── imagesuploader/
  │   ├── middleware/
  │   ├── models/
  │   ├── routes/
  │   ├── services/
  │   └── ws/
  ├── migrate/
  ├── pkg/
  ├── .env
  ├── .air.toml
  ├── docker-compose.yml
  ├── go.mod
  ├── go.sum
  └── Makefile
```

---

## Key Components

### 1. `cmd/api/main.go`
- **Entry point** for the backend server.
- Loads configuration, initializes logger, database, services, and starts the HTTP server.
- Sets up graceful shutdown and Swagger documentation.

### 2. `internal/`
Contains the main application logic, organized by concern:

#### a. `internal/auth/`
- JWT authentication, token generation, and validation.

#### b. `internal/config/`
- Application configuration loading (from environment variables).
- Rate limiter, database, Stripe, and other settings.

#### c. `internal/db/`
- Database connection and setup logic.

#### d. `internal/handlers/`
- HTTP request handlers for users, auctions, payments, images, etc.

#### e. `internal/imagesuploader/`
- Handles image upload, validation, and storage.

#### f. `internal/middleware/`
- Middleware for authentication, authorization, rate limiting, etc.

#### g. `internal/models/`
- Data models and request/response structs.

#### h. `internal/routes/`
- API route definitions and grouping.

#### i. `internal/services/`
- Business logic and service layer for users, auctions, notifications, etc.

#### j. `internal/ws/`
- WebSocket hub for real-time auction and notification updates.

---

### 3. `migrate/`
- Database migration files (SQL).

### 4. `pkg/`
- Utility packages (e.g., environment variable helpers, unique filename generator).

### 5. `Makefile`
- Common development tasks: migrations, Swagger docs, Docker commands.

### 6. `docker-compose.yml`
- Docker configuration for local development (database, server, etc.).

### 7. `.env`
- Environment variables for local development.

---

## API

- All endpoints are prefixed with `/api/v1`.
- Swagger docs available at `/api/v1/swagger/index.html`.
- Authentication uses JWT (via cookies or Authorization header).
- Admin and user endpoints are separated and protected by middleware.

---

## Notable Features

- **User Management:** Signup, login, profile, password change, admin user management.
- **Auction Management:** Create, view, bid, delete auctions (with admin controls).
- **Image Uploads:** Secure image upload and storage for auction items.
- **Payments:** Stripe integration for auction payments.
- **Notifications:** Real-time notifications via WebSockets.
- **Rate Limiting:** Configurable rate limiters for sensitive and general operations.
- **Swagger Documentation:** Auto-generated API docs for easy exploration.

---

## Getting Started

1. **Install dependencies:**  
   ```sh
   go mod download
   ```

2. **Set up environment:**  
   Copy `.env.example` to `.env` and adjust values as needed.

3. **Run database migrations:**  
   ```sh
   make mup
   ```

4. **Start the server:**  
   ```sh
   go run cmd/api/main.go
   ```
   Or use Docker Compose:
   ```sh
   make dkup
   ```

5. **Access API docs:**  
   Visit `http://localhost:<PORT>/api/v1/swagger/index.html`

---

## References

- `cmd/api/main.go`
- `internal/routes/routes.go`
- `internal/handlers/user.go`
- `pkg/env.go`
- `pkg/uniquefilename.go`
- `Makefile`

---

## Testing
- To run the tests, navigate to the server directory and execute the following command:

```go test ./...```
- This will run all the tests in the server directory.

## Contributing
- Contributions are welcome! If you'd like to contribute to this project, please fork the repository and submit a pull request.

# License
- This project is licensed under the MIT License.

For more details, see the inline comments and Swagger documentation in the codebase.