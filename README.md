# Bookify - Appointment Booking System

MVP implementation of a REST API service for managing schedules and booking appointments with specialists.

## Architecture

Follows Clean Architecture principles with clear separation of concerns:
- **Handlers**: HTTP layer, handles requests and responses
- **Services**: Business logic layer
- **Repositories**: Data access layer
- **Domain**: Core entities

## Features

- Browse available specialists (barber, dentist, psychologist)
- View available time slots for each specialist
- Book available time slots
- View your bookings
- Prevents double booking of the same time slot

## Tech Stack

- Go 1.21
- PostgreSQL 15
- docker-compose
- golang-migrate

## Quick Start

### Prerequisites
- Docker
- Docker Compose

### Running the Application

1. Clone the repository
```bash
cd bookify
```

2. Build and start the application with Docker Compose
```bash
docker-compose up --build
```

This will:
- Start a PostgreSQL database
- Run migrations automatically
- Start the API server on `http://localhost:8080`

### Stopping the Application

```bash
docker-compose down
```

To remove data as well:
```bash
docker-compose down -v
```

## API Endpoints

### 1. Get All Specialists
```
GET /specialists
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Smith",
    "type": "barber"
  },
  {
    "id": 2,
    "name": "Dr. Sarah Johnson",
    "type": "dentist"
  },
  {
    "id": 3,
    "name": "Dr. Michael Lee",
    "type": "psychologist"
  }
]
```

### 2. Get Specialist with Available Time Slots
```
GET /specialistsWithSlots/{id}
```

**Example Request:**
```
GET /specialistsWithSlots/1
```

**Response:**
```json
{
  "specialist": {
    "id": 1,
    "name": "John Smith",
    "type": "barber"
  },
  "time_slots": [
    {
      "id": 1,
      "specialist_id": 1,
      "time": "2026-03-27T09:00:00Z",
      "is_booked": false
    },
    {
      "id": 2,
      "specialist_id": 1,
      "time": "2026-03-27T10:00:00Z",
      "is_booked": false
    }
  ]
}
```

### 3. Create a Booking
```
POST /bookings
Content-Type: application/json
```

**Request Body:**
```json
{
  "time_slot_id": 1
}
```

**Success Response (201 Created):**
```json
{
  "id": 1,
  "user_id": 1,
  "time_slot_id": 1,
  "created_at": "2026-03-26T12:00:00Z"
}
```

**Error Responses:**
- `409 Conflict` - Slot already booked
```json
{
  "error": "this slot is already booked"
}
```

- `404 Not Found` - Time slot not found
```json
{
  "error": "time slot not found"
}
```

### 4. Get Your Bookings
```
GET /bookings
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "time_slot_id": 1,
    "created_at": "2026-03-26T12:00:00Z",
    "specialist": "John Smith",
    "slot_time": "2026-03-27T09:00:00Z"
  }
]
```

## Testing the API

### Using curl

1. **Get all specialists:**
```bash
curl http://localhost:8080/specialists
```

2. **Get specialist with slots:**
```bash
curl http://localhost:8080/specialistsWithSlots/1
```

3. **Book a time slot:**
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{"time_slot_id": 1}'
```

4. **Get your bookings:**
```bash
curl http://localhost:8080/bookings
```

### Using Postman

Import the endpoints above into Postman for easier testing.

## Project Structure

```
bookify/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   └── postgres.go             # Database connection
│   ├── domain/
│   │   ├── specialist.go           # Specialist entity
│   │   ├── time_slot.go            # TimeSlot entity
│   │   └── booking.go              # Booking entity
│   ├── handlers/
│   │   └── handler.go              # HTTP handlers
│   ├── middleware/
│   │   └── middleware.go           # HTTP middleware
│   ├── repository/
│   │   ├── specialist_repository.go
│   │   ├── time_slot_repository.go
│   │   └── booking_repository.go
│   └── service/
│       ├── specialist_service.go   # Business logic for specialists
│       └── booking_service.go      # Business logic for bookings
├── migrations/
│   ├── 000001_create_specialists_table.up.sql
│   ├── 000001_create_specialists_table.down.sql
│   ├── 000002_create_time_slots_table.up.sql
│   ├── 000002_create_time_slots_table.down.sql
│   ├── 000003_create_bookings_table.up.sql
│   └── 000003_create_bookings_table.down.sql
├── docker-compose.yml              # Docker services configuration
├── Dockerfile                       # Application container
├── go.mod                          # Go module definition
└── README.md                       # This file
```

## Development

### Local Development (without Docker)

1. Install Go 1.21 or later
2. Install PostgreSQL
3. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=bookify
export DB_SSLMODE=disable
```

4. Create the database:
```bash
createdb bookify
```

5. Run migrations:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/bookify?sslmode=disable" up
```

6. Build and run:
```bash
go build -o api ./cmd/api
./api
```

## Design Decisions

### MVP Simplifications
- **No Authentication**: Using hardcoded `userId = 1` for all requests
- **Pre-created Slots**: Time slots are pre-created in migrations
- **Simple Validation**: Basic checks for slot availability

### Architecture Notes
- Clean Architecture with clear layer separation
- Repository pattern for data access
- Service layer contains all business logic
- Handlers are thin, delegating to services
- No external frameworks, using Go standard library

## Future Enhancements
- User authentication with JWT
- Role-based access control (CLIENT, PROVIDER)
- Dynamic time slot generation
- Booking cancellation
- Email notifications
- Payment processing
- Advanced scheduling (recurring slots, exceptions)