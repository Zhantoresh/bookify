# Booking System API

REST API for appointment booking with JWT authentication, RBAC, PostgreSQL, migrations, background workers, and Docker-based local setup.

## Structure

The project follows the `cmd/internal/pkg` layout and includes:

- `cmd/api` for application bootstrap and graceful shutdown
- `internal/domain` for `users`, `services`, and `appointments`
- `internal/repository/postgres` for PostgreSQL persistence
- `internal/service` for auth, services, and appointments business logic
- `internal/transport/http` for handlers, middleware, and routing
- `internal/worker` for reminders and async task execution
- `migrations`, `seeds`, `docs`, and `tests`

## Quick start

```bash
cp .env.example .env
docker compose up --build
```

API is available at `http://localhost:8080`.

## Main endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/services`
- `GET /api/v1/services`
- `POST /api/v1/appointments`
- `GET /api/v1/appointments/my`
- `GET /api/v1/appointments/available-slots`
- `GET /health`

## Local development

```bash
make run
make test
make migrate-up
```

## Postman testing

Import:

- `docs/postman_collection.json`
- `docs/postman_environment.json`

Recommended order:

1. `00 Setup -> Health Check`
2. `01 Auth -> Register Provider`
3. `01 Auth -> Register Client`
4. `01 Auth -> Login Provider`
5. `00 Setup -> Set Provider Token`
6. `01 Auth -> Validate Active Token`
7. `02 Provider Flow -> Create Service`
8. `02 Provider Flow -> List My Services`
9. `03 Public Checks -> List Services`
10. `03 Public Checks -> Get Available Slots`
11. `01 Auth -> Login Client`
12. `00 Setup -> Set Client Token`
13. `04 Client Flow -> Create Appointment`
14. `04 Client Flow -> List My Appointments`
15. `04 Client Flow -> Get Appointment By ID`
16. `00 Setup -> Set Provider Token`
17. `05 Provider Appointment Actions -> Provider List My Appointments`
18. `05 Provider Appointment Actions -> Confirm Appointment`

Optional:

19. `06 Negative Cases -> Client Cannot Confirm Appointment`
20. `06 Negative Cases -> Booking In The Past Fails`

What the collection does automatically:

- saves `provider_token` and `client_token`
- switches active bearer token through `Set Provider Token` and `Set Client Token`
- saves `service_id` after service creation
- saves `appointment_id` after appointment creation

Important:

- `Create Service` must be called with provider token active
- `Create Appointment` must be called with client token active
- `Confirm Appointment` must be called with provider token active
- `appointment_start_time` must be in the future relative to current date
