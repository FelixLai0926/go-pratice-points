# go-pratice-points
 go-pratice
# Points System

A distributed points system built in Go that manages user accounts, transactions, and events.  

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Docker Deployment / Docker](#docker-deployment)
- [Running Tests](#running-tests)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Configuration Management**  
  Loads environment variables from YAML configuration files `(configs/*.yaml)` and parses them into strongly typed configuration structs.

- **Logging**  
  Structured logging using Uber's Zap, supporting both development and production modes.  

- **Database Initialization & Migrations**  
  Supports PostgreSQL as the primary database, with automatic migration execution on startup using `migrate`.

- **Transaction Management**  
  Implements Unit of Work (UoW) to handle database transactions efficiently, ensuring consistency and rollback on failures.

- **Distributed Locking**
  Implements distributed locking using Redis (via redislock) to prevent race conditions and ensure data consistency.
- **RESTful API**
  Provides well-structured API endpoints with request validation and response DTOs.
- **Testing**
  Comprehensive unit and integration tests using Testify, miniredis, and Testcontainers for mocking dependencies.
---

## Project Structure

```plaintext
├─cmd
│  ├─genmodel        # CLI tool to generate models
│  └─points          # Main entry point for the application
├─configs            # Configuration files (e.g., YAML, JSON, ENV)
├─docker             # Docker-related files and configurations
├─internal           # Core application logic (follows Clean Architecture)
│  ├─adapter         # Adapters for external interfaces (e.g., HTTP, gRPC)
│  │  └─http
│  │      ├─controller   # HTTP request handlers
│  │      ├─dto          # Data Transfer Objects (DTOs) for API requests/responses
│  │      ├─middleware   # Middleware for request processing
│  │      └─router       # HTTP route definitions
│  ├─domain          # Business domain layer (core logic)
│  │  ├─command        # Application commands (CQRS pattern)
│  │  ├─entity         # Business entities (domain models)
│  │  ├─event          # Domain events
│  │  ├─port           # Interfaces for dependency inversion (input/output ports)
│  │  ├─repository     # Interfaces for data persistence
│  │  └─valueobject    # Value objects (immutable domain concepts)
│  ├─infrastructure   # Implementation details (external dependencies)
│  │  ├─dbconnection   # Database connection handling
│  │  ├─distributedlock # Distributed locking mechanisms
│  │  └─persistence    # Data persistence layer (ORM, repositories)
│  │      ├─gorm
│  │      │  ├─dao      # Data Access Objects (DAOs) using GORM
│  │      │  └─model    # ORM models
│  │      └─repository  # Repository implementations
│  ├─shared           # Shared utilities and error handling
│  │  ├─apperror      # Custom application errors
│  │  ├─errcode       # Error codes
│  │  └─mapper        # Object mappers and conversions
│  └─usecase          # Application use cases (business logic)
│      ├─locking       # Use cases related to distributed locking
│      └─transaction   # Use cases related to transactions
├─migrations         # Database migration files
└─test
    └─mock           # Mock implementations for testing
```

## Prerequisites

- **Go**: Version 1.18 or higher (supports generics)  
- **Docker**: Required for integration tests using Testcontainers
- **Redis** and **PostgreSQL**: Required in production for caching, database, and distributed locks

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/points-system.git
   cd points-system
   ```
2. **Install dependencies:**
   ```bash
   go mod download
   ```
## Running the Application

1. **Build the application**
    ```bash
    go build -o points-system.exe
    ```
2. **Run the application**
    ```bash
    ./points-system.exe
    ```

## Docker Deployment
The project can be containerized using Docker Compose. 

You can also specify the environment by setting an environment variable (e.g., APP_ENV) in the Compose file.

1. **Build the application**
    ```bash
    ./docker-compose up -d
    ```

## Running Tests
To run all tests (unit and integration), execute:

```bash
go test -v ./...
```
Note: Ensure Docker is running if integration tests are enabled.

## Configuration
The project uses environment variables for configuration, with files named <environment>.yaml

Place your configuration files in the configs/ directory. For example:

```plaintext
example.yaml (default)
production.yaml
development.yaml
```
The config module loads and parses these configuration files.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

For major changes, please open an issue first to discuss what you would like to change.


## License
This project is licensed under the MIT License - see the LICENSE file for details.