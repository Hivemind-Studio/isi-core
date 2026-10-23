# Golang Backend Project

## Overview
This project is a modular and scalable backend application built with Golang. It is designed to support a service platform that includes features like role-based access control (RBAC), user management, coach booking, and payment processing. The project follows a clean architecture and is structured for maintainability and extensibility.

## Project Structure

```plaintext
.
├── cmd
│   └── api
│       ├── main.go         # Entry point of the application
│       └── server.go       # Server initialization logic
├── configs
│   └── config.go           # Application configuration handling
├── db
│   ├── migrations          # SQL scripts for database schema migrations
│   └── migrations.go       # Migration management code
├── internal
│   ├── constant            # Static constants (e.g., roles, URLs)
│   ├── dto                 # Data Transfer Objects for API payloads
│   ├── enum                # Enumerations for static values (e.g., gender)
│   ├── handler             # HTTP handlers for API endpoints
│   ├── repository          # Database interaction logic
│   ├── service             # Business logic implementation
│   └── usecase             # Use case implementations (e.g., creating a coach)
├── pkg
│   ├── dbtx                # Database transaction management
│   ├── logger              # Logging utilities
│   ├── mail                # Email handling
│   ├── middleware          # Middleware (e.g., JWT, request logging)
│   ├── mysqlconn           # MySQL connection setup
│   ├── redis               # Redis integration
│   ├── s3                  # S3 client for file storage
│   ├── translator          # Validation translators
│   └── validator           # Input validation utilities
├── template                # HTML templates for emails
├── utils                   # General-purpose utility functions
└── go.mod / go.sum         # Dependency management files
```

## Features

1. **Authentication and Authorization**
    - JWT-based authentication
    - Role-based access control (RBAC)

2. **Database**
    - MySQL integration
    - Migration scripts for schema management

3. **Domain-Specific Modules**
    - Users, roles, coaches, bookings, and payments

4. **Email Services**
    - Email templates for verification and password recovery

5. **Third-Party Integrations**
    - AWS S3 for file storage
    - Redis for caching

6. **Error Handling and Validation**
    - Custom HTTP error handling
    - Input validation utilities

## Setup Instructions

### Prerequisites
- Go 1.20+
- MySQL
- Redis
- AWS S3 credentials (if using file storage)

### Steps
1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd <repository_name>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database:
    - Update `configs/config.go` with your MySQL credentials.
    - Run migrations:
      ```bash
      go run db/migrations.go
      ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

5. Access the API at `http://localhost:<port>` (default: `8080`).

## Contribution Guidelines

1. Fork the repository and create a new branch for your feature/bugfix.
2. Write clear commit messages and ensure your code follows the project style.
3. Submit a pull request with a detailed description of your changes.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contact
For any questions or feedback, please contact the project maintainers at [email@example.com].
