# Tasks API

A robust RESTful API service built with Go for managing tasks.

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.24 or higher
- **SQLite**: Database engine

### 🛠️ Installation

1.  **Clone the repository**

2.  **Setup Environment**
    Copy the example environment file to create your local configuration:
    ```bash
    cp .env.example .env
    ```

3.  **Install Dependencies**
    ```bash
    go mod download
    ```

### 🏃‍♂️ Running the Application

To start the API server locally:

```bash
go run cmd/api/main.go
```

## 📁 Project Structure

- **`cmd/api`**: Entry point of the application.
- **`internal`**: Private application code (business logic, data access, etc.).
- **`pkg`**: Public libraries that can be used by external applications.
- **`migrations`**: Database schema migrations.

## 🗄️ Database Migrations

This project uses [goose](https://github.com/pressly/goose) for database migrations.

### Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Run Migrations

To apply all available migrations:

```bash
goose -dir migrations sqlite3 ./database.db up
```

To create a new migration:

```bash
goose -dir migrations sqlite3 ./database.db create add_users_table sql
```

## 🧪 Testing

Run the test suite using standard Go tooling:

```bash
go test ./...
```
