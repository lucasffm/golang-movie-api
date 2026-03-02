# Movie API

A RESTful API for managing movies and reviews, built with Go, Gin, and GORM.
The project follows **Clean Architecture** with dependency injection to keep business logic decoupled from frameworks and databases.

---

## Table of Contents

- [Architecture](#architecture)
  - [Layers](#layers)
  - [Dependency Rule](#dependency-rule)
  - [Data Flow](#data-flow)
  - [Adding a New Feature](#adding-a-new-feature)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Setup](#setup)
- [Running the Server](#running-the-server)
- [Environment Variables](#environment-variables)
- [API Reference](#api-reference)
- [Example Requests](#example-requests)

---

## Architecture

This project uses Clean Architecture, which organises code into concentric layers where **inner layers never depend on outer ones**. The goal is to keep business rules independent of the web framework, database driver, or any other infrastructure detail.

```
┌─────────────────────────────────────────┐
│              HTTP (Gin)                 │  ← handler layer
│  ┌───────────────────────────────────┐  │
│  │           Use Cases               │  │  ← usecase layer
│  │  ┌─────────────────────────────┐  │  │
│  │  │          Domain             │  │  │  ← domain layer (core)
│  │  │  entities · interfaces      │  │  │
│  │  └─────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
│              Repository (GORM)          │  ← repository layer
└─────────────────────────────────────────┘
```

### Layers

| Layer | Package | Responsibility |
|---|---|---|
| **Domain** | `internal/domain` | Entities (`Movie`, `Review`), repository and use-case **interfaces**, DTOs (input structs), domain errors, and database migration. Has **no dependencies** on other internal packages. |
| **Repository** | `internal/repository` | Implements the `domain.MovieRepository` and `domain.ReviewRepository` interfaces using GORM. The only place in the project where SQL/GORM code lives. Translates framework errors (e.g. `gorm.ErrRecordNotFound`) into domain errors. |
| **Use Case** | `internal/usecase` | Implements the `domain.MovieUseCase` and `domain.ReviewUseCase` interfaces. Contains all business logic (validation, orchestration). Depends only on domain interfaces — it never imports GORM or Gin. |
| **Handler** | `internal/handler` | Receives HTTP requests via Gin, delegates work to a use-case, and writes the HTTP response. Depends only on domain interfaces and types. Maps domain errors to HTTP status codes. |
| **Routes** | `routes` | Registers all Gin routes, attaches middleware, and maps each route to the correct handler method. |
| **Config** | `config` | Reads environment variables and opens the database connection, returning a `*gorm.DB`. |
| **Main** | `main.go` | The composition root. Instantiates every dependency in order and wires them together (dependency injection by hand). |

### Dependency Rule

Dependencies always point **inward**:

```
handler → usecase → domain ← repository
                              ↑
                           (implements domain interfaces)
```

- `handler` depends on `domain` interfaces (never on `repository` or GORM).
- `usecase` depends on `domain` interfaces (never on `handler` or Gin).
- `repository` depends on `domain` entities and GORM.
- `domain` depends on nothing internal.

This means you can swap the database (e.g. replace GORM+Postgres with a different ORM or an in-memory store) without touching a single line of business logic, and you can test use cases without a real database.

### Data Flow

A typical request walks through the layers in this order:

```
HTTP Request
    │
    ▼
[routes/routes.go]          ← selects handler method
    │
    ▼
[internal/handler]          ← parses & validates HTTP input, calls use case
    │
    ▼
[internal/usecase]          ← applies business rules, calls repository
    │
    ▼
[internal/repository]       ← executes SQL via GORM, returns domain entity
    │
    ▼
[internal/usecase]          ← receives entity, may enrich or validate further
    │
    ▼
[internal/handler]          ← serialises entity to JSON, writes HTTP response
    │
    ▼
HTTP Response
```

### Adding a New Feature

Say you need to add a `Watch List` resource. The steps are always the same:

1. **Domain** — define the `WatchList` entity, `CreateWatchListInput` / `UpdateWatchListInput` DTOs, `ErrWatchListNotFound`, and the `WatchListRepository` + `WatchListUseCase` interfaces in `internal/domain/watchlist.go`.
2. **Repository** — implement `WatchListRepository` with GORM in `internal/repository/watchlist_repository.go`.
3. **Use Case** — implement `WatchListUseCase` in `internal/usecase/watchlist_usecase.go`, injecting the repository interface.
4. **Handler** — implement `WatchListHandler` in `internal/handler/watchlist_handler.go`, injecting the use-case interface.
5. **Routes** — register the new routes in `routes/routes.go`.
6. **Wire** — instantiate and inject in `main.go`.

At no step does a lower-numbered layer need to know about a higher-numbered one.

---

## Project Structure

```
movie-api/
├── config/
│   └── database.go                 # Opens the Postgres connection
├── internal/
│   ├── domain/
│   │   ├── migration.go            # AutoMigrate for all entities
│   │   ├── movie.go                # Movie entity, interfaces & DTOs
│   │   └── review.go               # Review entity, interfaces & DTOs
│   ├── handler/
│   │   ├── helpers.go              # Shared utilities (parseID)
│   │   ├── movie_handler.go        # HTTP handlers for /movies
│   │   └── review_handler.go       # HTTP handlers for /reviews
│   ├── repository/
│   │   ├── movie_repository.go     # GORM implementation of MovieRepository
│   │   └── review_repository.go    # GORM implementation of ReviewRepository
│   └── usecase/
│       ├── movie_usecase.go        # Business logic for movies
│       └── review_usecase.go       # Business logic for reviews
├── routes/
│   └── routes.go                   # Route registration & middleware
├── .env.example                    # Environment variable template
├── go.mod
├── go.sum
└── main.go                         # Composition root (dependency injection)
```

---

## Requirements

- [Go](https://go.dev/dl/) 1.21+
- [PostgreSQL](https://www.postgresql.org/) 13+

Optional, for live reload during development:

- [`gin`](https://github.com/codegangsta/gin) — `go install github.com/codegangsta/gin@latest`

---

## Setup

**1. Clone the repository**

```bash
git clone <repository-url>
cd movie-api
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Create the database**

```sql
CREATE DATABASE movies_db;
```

**4. Configure environment variables**

```bash
cp .env.example .env
```

Edit `.env` with your database credentials (see [Environment Variables](#environment-variables)).

The application runs migrations automatically on startup — no manual SQL scripts needed.

---

## Running the Server

**Standard**

```bash
go run main.go
```

**Live reload** (requires the `gin` CLI tool)

```bash
gin run main.go
```

The server starts on port **8080**.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `postgres` | Postgres user |
| `DB_PASSWORD` | `postgres` | Postgres password |
| `DB_NAME` | `movies_db` | Database name |

---

## API Reference

All routes are prefixed with `/api/v1`.

### Movies

> `GET /api/v1/movies` requires HTTP Basic Auth (`admin` / `admin`).

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/v1/movies` | Basic Auth | List all movies (includes reviews) |
| `GET` | `/api/v1/movies/:id` | — | Get a movie by ID |
| `POST` | `/api/v1/movies` | — | Create a movie |
| `PUT` | `/api/v1/movies/:id` | — | Update a movie (partial update) |
| `DELETE` | `/api/v1/movies/:id` | — | Soft-delete a movie |
| `GET` | `/api/v1/movies/:id/reviews` | — | List reviews for a movie |

**Movie object**

```json
{
  "id": 1,
  "title": "The Matrix",
  "director": "The Wachowskis",
  "year": 1999,
  "genre": "Sci-Fi",
  "description": "A computer hacker learns about the true nature of reality.",
  "created_at": "2024-01-15T20:00:00Z",
  "updated_at": "2024-01-15T20:00:00Z",
  "reviews": []
}
```

### Reviews

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/reviews` | List all reviews (includes movie) |
| `GET` | `/api/v1/reviews/:id` | Get a review by ID |
| `POST` | `/api/v1/reviews` | Create a review |
| `PUT` | `/api/v1/reviews/:id` | Update a review (partial update) |
| `DELETE` | `/api/v1/reviews/:id` | Soft-delete a review |

**Review object**

```json
{
  "id": 1,
  "movie_id": 1,
  "rating": 9,
  "comment": "Mind-bending classic!",
  "watched_at": "2024-01-15T20:00:00Z",
  "created_at": "2024-01-15T20:00:00Z",
  "updated_at": "2024-01-15T20:00:00Z",
  "movie": {}
}
```

---

## Example Requests

### Movies

#### List all movies
```bash
curl -u admin:admin http://localhost:8080/api/v1/movies/
```

#### Get movie by ID
```bash
curl http://localhost:8080/api/v1/movies/1
```

#### Create a movie
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "director": "The Wachowskis",
    "year": 1999,
    "genre": "Sci-Fi",
    "description": "A computer hacker learns about the true nature of reality."
  }'
```

#### Update a movie
```bash
curl -X PUT http://localhost:8080/api/v1/movies/1 \
  -H "Content-Type: application/json" \
  -d '{
    "genre": "Sci-Fi, Action"
  }'
```

#### Delete a movie
```bash
curl -X DELETE http://localhost:8080/api/v1/movies/1
```

#### Get reviews for a movie
```bash
curl http://localhost:8080/api/v1/movies/1/reviews
```

### Reviews

#### List all reviews
```bash
curl http://localhost:8080/api/v1/reviews/
```

#### Get review by ID
```bash
curl http://localhost:8080/api/v1/reviews/1
```

#### Create a review
```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "rating": 9,
    "comment": "Mind-bending classic!",
    "watched_at": "2024-01-15T20:00:00Z"
  }'
```

#### Update a review
```bash
curl -X PUT http://localhost:8080/api/v1/reviews/1 \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 10,
    "comment": "An absolute masterpiece."
  }'
```

#### Delete a review
```bash
curl -X DELETE http://localhost:8080/api/v1/reviews/1
```
