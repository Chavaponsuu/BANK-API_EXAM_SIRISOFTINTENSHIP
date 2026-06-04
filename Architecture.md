# Project Architecture

This project follows a **Layered Architecture** (often referred to as the **Controller-Service-Repository** pattern). This architecture emphasizes the **Separation of Concerns**, making the codebase easier to maintain, test, and scale.

It is a standard and highly recommended approach for building RESTful APIs in Go (especially when using frameworks like Gin).

## Directory Structure & Responsibilities

The application is divided into several distinct layers, each with a specific responsibility in handling a request.

### 1. Presentation Layer (HTTP & Routing)
Responsible for receiving HTTP requests, basic validation, and returning HTTP responses.

*   **`routes/`**: Defines the API endpoints (URLs) and maps them to the appropriate handlers.
*   **`handlers/` (Controllers)**: The entry point for logic. Handlers extract data from the incoming HTTP request (URL parameters, JSON body), perform basic validation, call the appropriate Service, and format the final HTTP response (e.g., returning JSON with correct HTTP status codes).
*   **`middleware/`**: Intercepts requests before they reach the handlers. Used for cross-cutting concerns like Authentication (JWT), Logging, CORS, and Request Timeout handling.
*   **`dto/` (Data Transfer Objects)**: Structs specifically designed for incoming requests and outgoing responses. This prevents exposing internal database models directly to the client.

### 2. Business Logic Layer
The core of the application where all the rules and logic reside.

*   **`services/`**: Contains the business logic. A service orchestrates the flow of data, applies business rules, and communicates with repositories to fetch or save data. Services should *not* know anything about HTTP (Gin contexts) or the underlying database implementation.

### 3. Data Access Layer
Responsible for all interactions with the database.

*   **`repositories/`**: Abstracts the database operations (CRUD). It provides a clean interface for the services to interact with data. If the database technology changes (e.g., from MySQL to PostgreSQL), only the repository layer needs to be updated.

### 4. Domain Layer
The fundamental data structures of the application.

*   **`models/` (Entities)**: Go structs that represent the business entities and typically map directly to database tables (using an ORM like GORM).

### 5. Infrastructure & Configuration
Files related to setting up the application environment.

*   **`config/`**: Loads and manages environment variables and application configurations.
*   **`db/`**: Handles the initialization and connection pooling for the database.
*   **`migrations/`**: Contains SQL scripts or ORM logic for managing database schema changes over time.

## Request Flow Example

When a client makes an API request (e.g., `POST /users`), the flow typically looks like this:

1.  **Client** sends an HTTP Request.
2.  **Route** (`routes/`) catches `/users` and directs it to the User Handler.
3.  **Middleware** (`middleware/`) might check if the user is authenticated.
4.  **Handler** (`handlers/`) binds the JSON payload to a User DTO (`dto/`), validates it, and passes the data to the User Service.
5.  **Service** (`services/`) applies business rules (e.g., "hash the password", "check if email exists"). It then calls the User Repository.
6.  **Repository** (`repositories/`) executes the SQL query (via GORM) using the User Model (`models/`) to save the data to the Database.
7.  **Repository** returns success/failure to the Service.
8.  **Service** returns the result to the Handler.
9.  **Handler** constructs a JSON response (perhaps using another DTO) and sends it back to the **Client**.

## Advantages of this Architecture
*   **Maintainability**: Clear boundaries make it easy to find where a specific piece of code lives.
*   **Testability**: The clear separation allows for easy Unit Testing. You can easily mock the Repository to test the Service logic without needing a real database.
*   **Flexibility**: You can change the transport layer (e.g., move from REST to gRPC) or the database layer without affecting the core business logic in the services.
