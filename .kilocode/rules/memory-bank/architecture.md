# Architecture

## System Architecture
The project follows Hexagonal Architecture (Ports and Adapters) with clean separation of concerns:

- **Domain Layer**: Core business entities and logic (User, Account, Transaction, InterestHistory, Profile)
- **Application Layer**: Use cases and business logic (UserService, AccountService)
- **Ports Layer**: Interfaces defining contracts (repositories, services)
- **Adapters Layer**: Implementations of ports (HTTP handlers, PostgreSQL repositories, password service)

## Source Code Paths
- `cmd/api/main.go`: Application entry point
- `cmd/worker/main.go`: Background worker for interest calculation
- `internal/domain/`: Business entities and domain logic
- `internal/application/`: Use cases and services
- `internal/ports/`: Interface contracts
- `internal/adapters/`: Port implementations
- `internal/config/`: Configuration management
- `migrations/`: Database schema migrations

## Key Technical Decisions
- **Hexagonal Architecture**: Ensures testability and framework independence
- **GORM ORM**: For database operations with PostgreSQL
- **Echo Framework**: For REST API implementation
- **JWT Authentication**: For secure user sessions
- **Daily Interest Calculation**: Background worker processes interest for flexible savings accounts
- **Tiered Interest Rates**: Based on account age and balance thresholds

## Design Patterns
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic encapsulation
- **Dependency Injection**: Loose coupling between components
- **Factory Pattern**: Entity creation methods

## Component Relationships
- UserService depends on UserRepository, ProfileRepository, PasswordService
- AccountService depends on AccountRepository, UserRepository, ProfileRepository, TransactionRepository, InterestHistoryRepository
- HTTP handlers depend on services and use DTOs for request/response
- Repositories implement GORM-based data access with schema mapping