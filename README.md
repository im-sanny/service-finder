# Service Finder API

A robust, production-grade RESTful API built with **Go (Golang)** for managing service providers. This project demonstrates clean architecture, secure authentication, and advanced database interactions.

## Features

- **Clean Architecture:** Strictly layered design (Handler → Service → Repository) for maintainability and testability.
- **Secure Authentication:** JWT-based auth with bcrypt password hashing.
- **Advanced Data Handling:**
  - Pagination (`?page=1&limit=10`)
  - Dynamic Filtering (`?location=Dhaka&service_id=1`)
  - Batch Operations (Atomic Create/Delete)
- **Database Management:** Automated schema migrations using `goose`.
- **Production Ready:** Graceful shutdown, environment variable configuration, and connection pooling.

## Tech Stack

- **Language:** Go 1.26
- **Database:** PostgreSQL
- **Routing:** Standard Library `http.ServeMux` (Go 1.22+ patterns)
- **Authentication:** JSON Web Tokens (JWT) & `golang.org/x/crypto/bcrypt`
- **Migrations:** `pressly/goose`
- **Configuration:** `joho/godotenv`

## Project Structure

```text
service-finder/
├── cmd/
│   └── main.go              # Entry point, wiring & migration embedding
│
├── internal/                # Enforces private package boundaries
│   ├── config/
│   │   └── config.go        # Typed env loading & validation
│   ├── database/
│ 	│ 	└── migrations/      # SQL files (embedded relative to main.go)
│   │   └── postgres.go      # DB connection & migration runner
│   ├── handler/             # HTTP request/response handling
│   ├── middleware/          # Auth & security checks
│   ├── model/               # Data structures & DTOs
│   ├── repository/          # Database interactions & SQL queries
│   ├── routes/              # Route registration & mux setup
│   └── service/             # Business logic & validation
├── .air.toml                # Hot reload config (build cmd: ./cmd)
├── .env                     # Environment variables (gitignored)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL installed and running

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/im-sanny/service-finder.git
    cd service-finder
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Setup Environment:**
    Create a `.env` file in the root directory and add your credentials:

    ```env
    DB_URL=postgres://user:password@localhost:5432/serfin?sslmode=disable
    PORT=8080
    JWT_SECRET=your_super_secret_key
    API_KEY=your_api_key_for_testing
    ```

4.  **Run the application:**
    ```bash
    go run cmd/main.go
    ```
    _Note: The app will automatically apply database migrations on startup._

## 📡 API Endpoints

### Authentication

| Method | Endpoint         | Description             | Access |
| :----- | :--------------- | :---------------------- | :----- |
| `POST` | `/auth/register` | Register a new user     | Public |
| `POST` | `/auth/login`    | Login and get JWT token | Public |

### Services

| Method   | Endpoint         | Description                                      | Access    |
| :------- | :--------------- | :----------------------------------------------- | :-------- |
| `GET`    | `/services`      | Get all services (supports pagination/filtering) | Public    |
| `POST`   | `/services`      | Create a new service                             | Protected |
| `DELETE` | `/services/{id}` | Delete a service                                 | Protected |

### Providers

| Method | Endpoint           | Description             | Access    |
| :----- | :----------------- | :---------------------- | :-------- |
| `GET`  | `/providers`       | Get all providers       | Public    |
| `POST` | `/providers/batch` | Batch create providers  | Protected |
| `PUT`  | `/providers/{id}`  | Update provider details | Protected |

> **Protected Routes:** Require a valid JWT token in the `Authorization: Bearer <token>` header.

## 🧪 Testing with cURL

**Register a User:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Sanny","email":"sanny@test.com","password":"123"}'
```

**Login:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"sanny@test.com","password":"123"}'
```

**Get Paginated Providers:**

```bash
curl "http://localhost:8080/providers?page=1&limit=5&location=Dhaka"
```

## 🤝 Contributing

Feel free to fork this repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

---
