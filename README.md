# Bank Transfers API

A production-ready Golang REST API to handle internal account transfers, built with:

-   **Go 1.21+**
-   **Gin** (web framework)
-   **GORM** (ORM)
-   **PostgreSQL**
-   **golang-migrate** (for schema migrations)
-   **Docker + Compose**
-   **Air** for hot reloading in dev

## Getting Started

### Requirements

-   [Docker](https://www.docker.com/)
-   Make (optional but recommended)

## Clone & Setup

### Run the following commands.

    git clone https://github.com/einnar82/money-transfer.git
    cd money-transfer
    cp .env.example .env

## Running the App (via Docker Compose)

### Start the services

```
docker compose up --build
```
-   App will run at: [http://localhost:8888](http://localhost:8888)
-   DB will run on: `localhost:5432` (PostgreSQL 16)

## Environment Variables

 `.env` used for both Go app and Makefile:

```
DB_HOST=db
DB_PORT=5432
DB_NAME=appdb
DB_USER=postgres
DB_PASSWORD=postgres
APP_PORT=8888
```

## Database Migrations

 I used [`golang-migrate`](https://github.com/golang-migrate/migrate) in a dedicated container.

**Create a new migration:** `make migrations` 

**Run migrations:** `make migrate` 

**Rollback last migration:**  `make migrate-rollback` 


## Developer Experience

-   **Air** is used for hot reload in development.
    
-   Changes in Go code automatically reload the app.