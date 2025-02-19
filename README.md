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
  Loads environment variables from `.env` files and parses them into strongly typed configuration structs. 

- **Logging**  
  Structured logging using Uber's Zap, supporting both development and production modes.  

- **Database Initialization**  
  Supports in-memory SQLite for unit tests and PostgreSQL for production, with auto-migration of domain models.

- **Distributed Locking**  
  Implements distributed locking using Redis (via [redislock](https://github.com/bsm/redislock)) to ensure data consistency in concurrent scenarios.

- **Testing / 測試**  
  Comprehensive unit and integration tests using Testify, miniredis, and Testcontainers.  

---

## Project Structure

```plaintext
D:.
├─api                    # API layer code for handling HTTP requests
│  ├─test                # Test code for the Ping API
│  └─trade               # for handling trade operations (Transfer, Confirm, Cancel)
│      └─mock            # Dummy implementations for the TCC module, used in unit tests
├─cmd                    # Command-line tools (e.g., code generation, etc.)
│  └─codegen             # Code generation tool
├─configs                # Configuration files (e.g., .env files and other config files)
├─docker                 # Docker configuration and deployment scripts
│  └─sql                 # Docker configurations related to SQL databases
├─errors                 # Error handling module that centrally defines error types and error wrapping logic
├─pkg                    # Core business logic and modules
│  ├─dao                 # Data Access Object (DAO) layer code
│  ├─middleware          # Middleware implementations (e.g., logging, error handling, authentication, etc.)
│  ├─models              # Domain models (ORM definitions, enums, event payloads, etc.)
│  │  ├─enum             # Enumerated types definitions (e.g., errcode, tcc, etc.)
│  │  │  ├─errcode       # Error code enumerations
│  │  │  └─tcc           # TCC-related enumerations
│  │  ├─eventpayload     # Event payload definitions
│  │  ├─orm              # ORM model definitions (database entities)
│  │  └─trade            # Trade module model definitions
│  ├─module              # Module implementations (e.g., cache, config, database, distributed lock, logging, and test utilities)
│  │  ├─cache            # Cache module (e.g., Redis, etc.)
│  │  ├─config           # Configuration management module
│  │  ├─database         # Database initialization and related utilities
│  │  ├─distributedlock  # Distributed lock implementation (e.g., based on Redis)
│  │  ├─logger           # Logging module (based on Zap)
│  │  └─test             # Module-level test helper utilities
│  └─router              # Routing layer code (e.g., API routes using Gin)
├─repository             # Data repository layer, handling data persistence logic
└─service                # Service layer code
    └─tcc                # TCC-related business logic
        └─mock           # Mock implementations or test cases for the TCC module
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
    go build -o points-system .
    ```
2. **Run the application**
    ```bash
    ./points-system
    ```

## Docker Deployment
The project can be containerized using Docker Compose. Below is an example docker-compose.yml file:

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
The project uses environment variables for configuration, with files named .env.<environment>.

Place your configuration files in the configs/ directory. For example:

```plaintext
.env.example (default)
.env.production
.env.development
```
The config module loads and parses these configuration files.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

For major changes, please open an issue first to discuss what you would like to change.


## License
This project is licensed under the MIT License - see the LICENSE file for details.