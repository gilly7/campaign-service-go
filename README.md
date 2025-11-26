## Sample Run Output

```
PS C:\Users\ngeno\MyDev\campaign-service> go run cmd/server/main.go
2025/11/26 12:57:11 Looking for migrations in: C:\Users\ngeno\MyDev\campaign-service\migrations
2025/11/26 12:57:11 goose: no migrations to run. current version: 20251126
2025/11/26 12:57:11 Migrations applied successfully!
2025/11/26 12:57:11 Connected to PostgreSQL + migrations applied
2025/11/26 12:57:11 Service running. Press Ctrl+C to stop.
2025/11/26 12:57:11 Server starting on http://localhost:8080
2025/11/26 12:57:20 Queued message for alice@gmail.com
2025/11/26 12:57:20 Sending to alice@gmail.com: Hey UserA! Get 20% off – https://example.com/offer/alice
2025/11/26 12:57:20 Queued message for bob@gmail.com
2025/11/26 12:57:20 Queued message for ngenogilbert07@gmail.com
2025/11/26 12:57:20 Sent to alice@gmail.com
2025/11/26 12:57:20 Sending to bob@gmail.com: Hey UserB! Get 30% off – https://example.com/offer/bob
```

# Campaign Service – Go Backend Exercise

A production-grade marketing campaign service with personalization, async delivery, and clean architecture.

## Tech Stack

- **Go** 1.24+
- **PostgreSQL** (with automatic migrations via goose v3)
- **Redis** (Streams for async message processing)
- **Chi** (HTTP REST API router)
- **pgx v5** (PostgreSQL driver + pool)
- **text/template** (personalization engine)

## Features

- Clean hexagonal/clean architecture (`cmd/`, `internal/`)
- HTTP REST API with Chi router
- PostgreSQL with automatic migrations (goose v3)
- Redis Streams for reliable async message processing
- Background worker with consumer group
- Go text/template personalization engine (future AI-ready)
- Graceful shutdown
- Zero-config local development (Docker Compose)

## Project Structure

```
campaign-service/
├── cmd/server/main.go                # Application entrypoint
├── internal/
│   ├── api/          → HTTP handlers
│   ├── campaign/     → Business logic & repository
│   ├── database/     → DB connection + migrations
│   ├── message/      → Templating + Redis queue
│   └── models/       → Domain models
├── migrations/       → SQL migration files
├── docker-compose.yml
├── go.mod
└── README.md
```

## Quick Start (Windows 11 / PowerShell)

```powershell
# 1. Clone & enter directory
git clone https://github.com/gilly7/campaign-service-go.git
cd campaign-service

# 2. Start PostgreSQL & Redis
docker compose up -d

# 3. Run the service
go run cmd/server/main.go
```

You should see:

- Migrations applied successfully!
- Connected to PostgreSQL + migrations applied
- Server starting on http://localhost:8080
- Background worker started – processing queue
- Service running. Press Ctrl+C to stop.

## Test the API

```powershell
# Health check
Invoke-RestMethod http://localhost:8080/health

# Create a campaign
Invoke-RestMethod http://localhost:8080/campaigns -Method POST -ContentType "application/json" -Body '{
  "name": "Black Friday Blast",
  "template": "Hey {{.FirstName}}! Get {{.Discount}}% off – claim here: {{.OfferURL}}",
  "user_ids": ["alice", "bob", "ngenogilbert07"]
}'
```

Watch the server console – you’ll see real-time personalized messages being processed!

## API Endpoints

| Method | Endpoint         | Description                              |
|--------|------------------|------------------------------------------|
| GET    | /health          | Health check                             |
| POST   | /campaigns       | Create campaign + enqueue messages       |

### Request body example

```json
{
  "name": "Summer Sale",
  "template": "Hi {{.FirstName}}! Enjoy {{.Discount}}% off!",
  "user_ids": ["u123", "u456"]
}
```

### Response example

```json
{
  "id": "a1b2c3d4-...",
  "name": "Summer Sale",
  "template": "...",
  "status": "active",
  "created_at": "2025-11-26T12:00:00Z"
}
```

## Personalization

The template engine supports:

```gotemplate
Hi {{.FirstName}}!
You're getting {{.Discount}}% off.
Claim now: {{.OfferURL}}
```

Custom functions:

- `{{upper "hello"}}` → `HELLO`
- `{{formatDate .JoinDate "Jan 2006"}}`

Future-ready for AI: just add `{{aiRewrite "Make this exciting"}}`

## Development

```powershell
# Re-run after code changes
go run cmd/server/main.go

# Run tests (add your own in the future)
go test ./...
```

## Shutdown

Press Ctrl+C – the service shuts down gracefully:

- Stops accepting new HTTP requests
- Finishes in-flight messages
- Closes DB and Redis connections

Project Structure

campaign-service/
├── cmd/server/main.go                # Application entrypoint
├── internal/
│   ├── api/          → HTTP handlers
│   ├── campaign/     → Business logic & repository
│   ├── database/     → DB connection + migrations
│   ├── message/      → Templating + Redis queue
│   └── models/       → Domain models
├── migrations/       → SQL migration files
├── docker-compose.yml
├── go.mod
└── README.md

Quick Start (Windows 11 / PowerShell)powershell

# 1. Clone & enter directory
git clone https://github.com/yourname/campaign-service.git
cd campaign-service

# 2. Start PostgreSQL & Redis
docker compose up -d

# 3. Run the service
go run cmd/server/main.go

You should see:

Migrations applied successfully!
Connected to PostgreSQL + migrations applied
Server starting on http://localhost:8080
Background worker started – processing queue
Service running. Press Ctrl+C to stop.

Test the APIpowershell

# Health check
Invoke-RestMethod http://localhost:8080/health

# Create a campaign
Invoke-RestMethod http://localhost:8080/campaigns -Method POST -ContentType "application/json" -Body '{
  "name": "Black Friday Blast",
  "template": "Hey {{.FirstName}}! Get {{.Discount}}% off – claim here: {{.OfferURL}}",
  "user_ids": ["alice", "bob", "ngenogilbert254"]
}'

Watch the server console – you’ll see real-time personalized messages being processed!API EndpointsMethod
Endpoint
Description
GET
/health
Health check
POST
/campaigns
Create campaign + enqueue messages

Request body example:json

{
  "name": "Summer Sale",
  "template": "Hi {{.FirstName}}! Enjoy {{.Discount}}% off!",
  "user_ids": ["u123", "u456"]
}

Response:json

{
  "id": "a1b2c3d4-...",
  "name": "Summer Sale",
  "template": "...",
  "status": "active",
  "created_at": "2025-11-26T12:00:00Z"
}

PersonalizationThe template engine supports:gotemplate

Hi {{.FirstName}}!
You're getting {{.Discount}}% off.
Claim now: {{.OfferURL}}

Custom functions:{{upper "hello"}} → HELLO
{{formatDate .JoinDate "Jan 2006"}}

Future-ready for AI: just add {{aiRewrite "Make this exciting"}}Tech StackGo 1.24+
Chi – lightweight router
pgx v5 – PostgreSQL driver + connection pool
goose v3 – database migrations
Redis Streams – durable async queue
text/template – safe, powerful personalization

Developmentpowershell

# Re-run after code changes
go run cmd/server/main.go

# Run tests (add your own in the future)
go test ./...

ShutdownPress Ctrl+C – the service shuts down gracefully:Stops accepting new HTTP requests
Finishes in-flight messages
Closes DB and Redis connections

