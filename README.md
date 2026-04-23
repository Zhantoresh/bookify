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

Seed data is loaded automatically by Docker Compose after migrations.

Default seeded users:

- admin: `admin@booking.com` / `Admin123!`
- provider: `doctor@booking.com` / `Provider123!`
- client: `client@booking.com` / `Client123!`

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
2. `07 Admin -> Login Admin`
3. `07 Admin -> Set Admin Token`
4. `07 Admin -> Admin Dashboard`
5. `01 Auth -> Register Provider`
6. `01 Auth -> Register Client`
7. `01 Auth -> Login Provider`
8. `00 Setup -> Set Provider Token`
9. `01 Auth -> Validate Active Token`
10. `02 Provider Flow -> Create Service`
11. `02 Provider Flow -> List My Services`
12. `03 Public Checks -> List Services`
13. `03 Public Checks -> Get Available Slots`
14. `01 Auth -> Login Client`
15. `00 Setup -> Set Client Token`
16. `04 Client Flow -> Create Appointment`
17. `04 Client Flow -> List My Appointments`
18. `04 Client Flow -> Get Appointment By ID`
19. `00 Setup -> Set Provider Token`
20. `05 Provider Appointment Actions -> Provider List My Appointments`
21. `05 Provider Appointment Actions -> Confirm Appointment`

Optional:

22. `06 Negative Cases -> Client Cannot Confirm Appointment`
23. `06 Negative Cases -> Booking In The Past Fails`

What the collection does automatically:

- saves `provider_token` and `client_token`
- saves `admin_token`
- switches active bearer token through `Set Provider Token` and `Set Client Token`
- saves `service_id` after service creation
- saves `appointment_id` after appointment creation

Important:

- `Create Service` must be called with provider token active
- `Create Appointment` must be called with client token active
- `Confirm Appointment` must be called with provider token active
- `appointment_start_time` must be in the future relative to current date
- seeded admin credentials: `admin@booking.com` / `Admin123!`

Admin endpoints:

- `GET /api/v1/admin/dashboard`
- `GET /api/v1/admin/users`
- `GET /api/v1/admin/users/{id}`
- `PATCH /api/v1/admin/users/{id}` with body `{"role":"provider|client|admin"}`
- `DELETE /api/v1/admin/users/{id}`
