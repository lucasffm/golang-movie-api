# Movie API

A RESTful API built with Go and Gin framework for managing movies and reviews.

## Requirements

- Go 1.21+
- PostgreSQL

## Setup

1. Clone and navigate to the project:
```bash
cd movie-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create PostgreSQL database:
```sql
CREATE DATABASE movies_db;
```

4. Configure environment (copy `.env.example` to `.env` and update values):
```bash
cp .env.example .env
```

5. Run the server:
```bash
go run main.go
```

## API Endpoints

### Movies

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/movies` | List all movies |
| GET | `/api/v1/movies/:id` | Get a movie by ID |
| POST | `/api/v1/movies` | Create a new movie |
| PUT | `/api/v1/movies/:id` | Update a movie |
| DELETE | `/api/v1/movies/:id` | Delete a movie |
| GET | `/api/v1/movies/:id/reviews` | Get reviews for a movie |

### Reviews

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/reviews` | List all reviews |
| GET | `/api/v1/reviews/:id` | Get a review by ID |
| POST | `/api/v1/reviews` | Create a new review |
| PUT | `/api/v1/reviews/:id` | Update a review |
| DELETE | `/api/v1/reviews/:id` | Delete a review |

## Example Requests

### Movies

#### List All Movies
```bash
curl http://localhost:8080/api/v1/movies
```

#### Get Movie by ID
```bash
curl http://localhost:8080/api/v1/movies/1
```

#### Create Movie
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"The Matrix","director":"The Wachowskis","year":1999,"genre":"Sci-Fi","description":"A computer hacker learns about the true nature of reality"}'
```

#### Update Movie
```bash
curl -X PUT http://localhost:8080/api/v1/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"The Matrix","director":"The Wachowskis","year":1999,"genre":"Sci-Fi, Action","description":"Updated description"}'
```

#### Delete Movie
```bash
curl -X DELETE http://localhost:8080/api/v1/movies/1
```

#### Get Reviews for a Movie
```bash
curl http://localhost:8080/api/v1/movies/1/reviews
```

### Reviews

#### List All Reviews
```bash
curl http://localhost:8080/api/v1/reviews
```

#### Get Review by ID
```bash
curl http://localhost:8080/api/v1/reviews/1
```

#### Create Review
```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{"movie_id":1,"rating":9,"comment":"Mind-bending classic!","watched_at":"2024-01-15T20:00:00Z"}'
```

#### Update Review
```bash
curl -X PUT http://localhost:8080/api/v1/reviews/1 \
  -H "Content-Type: application/json" \
  -d '{"rating":10,"comment":"Absolutely masterpiece!"}'
```

#### Delete Review
```bash
curl -X DELETE http://localhost:8080/api/v1/reviews/1
```
