# 📅 Bookify - Booking System API

> **Complete, production-ready REST API** for appointment booking with JWT authentication, role-based access control, PostgreSQL, background workers, and comprehensive documentation.

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![API](https://img.shields.io/badge/API-OpenAPI%203.0.3-brightgreen.svg)
![Language](https://img.shields.io/badge/language-Go%201.23-00ADD8.svg)

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Postman (optional, for API testing)
- curl or similar HTTP client

### Start the API

```bash
# Clone and navigate
cd /Users/tsoy/Golang/bookify

# Copy environment (if needed)
cp .env.example .env

# Start everything with Docker Compose
docker compose up --build

# API will be available at http://localhost:8080
```

### Verify it's Running

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status": "ok"}
```

### First API Call - Register & Login

```bash
# 1. Register a new user (client)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe",
    "role": "client",
    "phone": "+1234567890"
  }'

# Response:
# {
#   "id": "usr_550e8400e29b41d4a716446655440000",
#   "email": "user@example.com",
#   "full_name": "John Doe",
#   "role": "client",
#   "created_at": "2026-04-20T10:00:00Z"
# }

# 2. Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'

# Response:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "user": {
#     "id": "usr_550e8400e29b41d4a716446655440000",
#     "email": "user@example.com",
#     "full_name": "John Doe",
#     "role": "client"
#   }
# }

# 3. Use token for authenticated requests
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📚 Documentation Files

| File | Purpose | Format |
|------|---------|--------|
| **swagger.yaml** | Complete OpenAPI 3.0.3 specification | YAML |
| **postman_collection.json** | Ready-to-import Postman collection with all endpoints | JSON |
| **postman_environment.json** | Pre-configured variables for Postman | JSON |

### How to Use Documentation

**Swagger UI (Web-based)**
```
1. Go to https://editor.swagger.io/
2. File → Import URL
3. Paste: raw GitHub URL to swagger.yaml
4. Or: Copy-paste entire yaml file content
```

**Postman (Desktop/Web)**
```
1. Open Postman
2. Click "Import" → "Upload Files"
3. Select both postman_collection.json and postman_environment.json
4. Click "Import"
5. Select environment from dropdown
6. Start testing!
```

---

## 🏗️ API Architecture

### Core Components

```
Bookify API
├── Authentication (JWT, OAuth-ready)
├── Services (Provider-managed)
├── Appointments (Client booking, Provider confirmation)
├── Users (Client, Provider, Admin roles)
├── Admin Dashboard (System monitoring)
├── Background Workers (Reminders, Notifications)
└── Database (PostgreSQL with migrations)
```

### User Roles & Permissions

| Feature | Client | Provider | Admin |
|---------|--------|----------|-------|
| Register/Login | ✅ | ✅ | ✅ |
| Create Services | ❌ | ✅ | ❌ |
| Book Appointments | ✅ | ❌ | ❌ |
| Confirm Appointments | ❌ | ✅ | ❌ |
| View Own Appointments | ✅ | ✅ | ❌ |
| View All Appointments | ❌ | ❌ | ✅ |
| View All Services | ✅ | ✅ | ✅ |
| Admin Dashboard | ❌ | ❌ | ✅ |
| List All Users | ❌ | ❌ | ✅ |

---

## 🔌 Complete API Reference

### Authentication Endpoints

#### Register New User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "provider@example.com",
  "password": "SecurePass123!",
  "full_name": "Dr. Jane Smith",
  "role": "provider",
  "phone": "+1234567890"
}
```

**Response (201 Created)**
```json
{
  "id": "usr_550e8400e29b41d4a716446655440000",
  "email": "provider@example.com",
  "full_name": "Dr. Jane Smith",
  "role": "provider",
  "phone": "+1234567890",
  "created_at": "2026-04-20T10:00:00Z",
  "updated_at": "2026-04-20T10:00:00Z"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "provider@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK)**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsInJvbGUiOiJwcm92aWRlciIsImV4cCI6MTcxMzYzMDQwMH0.signature",
  "user": {
    "id": "usr_550e8400e29b41d4a716446655440000",
    "email": "provider@example.com",
    "full_name": "Dr. Jane Smith",
    "role": "provider"
  }
}
```

#### Validate Token
```http
POST /api/v1/auth/validate
Authorization: Bearer {token}
```

**Response (200 OK)** - Token is valid

#### Get Current User
```http
GET /api/v1/users/me
Authorization: Bearer {token}
```

**Response (200 OK)**
```json
{
  "id": "usr_550e8400e29b41d4a716446655440000",
  "email": "provider@example.com",
  "full_name": "Dr. Jane Smith",
  "role": "provider",
  "phone": "+1234567890",
  "created_at": "2026-04-20T10:00:00Z"
}
```

---

### Services Endpoints (Provider)

#### Create Service
```http
POST /api/v1/services
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Dental Checkup",
  "description": "Professional dental examination and cleaning",
  "price": 150.00,
  "duration_minutes": 45
}
```

**Response (201 Created)**
```json
{
  "id": "svc_660e8400e29b41d4a716446655440000",
  "provider_id": "usr_550e8400e29b41d4a716446655440000",
  "name": "Dental Checkup",
  "description": "Professional dental examination and cleaning",
  "price": 150.00,
  "duration_minutes": 45,
  "is_active": true,
  "created_at": "2026-04-20T10:15:00Z"
}
```

#### List My Services (Provider)
```http
GET /api/v1/services/my
Authorization: Bearer {token}
```

**Response (200 OK)**
```json
[
  {
    "id": "svc_660e8400e29b41d4a716446655440000",
    "name": "Dental Checkup",
    "price": 150.00,
    "duration_minutes": 45,
    "is_active": true
  }
]
```

#### Get Service Details (Public)
```http
GET /api/v1/services/{id}
```

#### Update Service (Full)
```http
PUT /api/v1/services/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Premium Dental Checkup",
  "price": 200.00,
  "duration_minutes": 60
}
```

#### Partial Update Service
```http
PATCH /api/v1/services/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "price": 175.00,
  "is_active": true
}
```

#### Delete Service
```http
DELETE /api/v1/services/{id}
Authorization: Bearer {token}
```

**Response (204 No Content)**

---

### Appointments Endpoints

#### Create Appointment (Client)
```http
POST /api/v1/appointments
Authorization: Bearer {token}
Content-Type: application/json

{
  "service_id": "svc_660e8400e29b41d4a716446655440000",
  "start_time": "2026-04-25T14:00:00Z",
  "notes": "Please arrive 10 minutes early"
}
```

**Response (201 Created)**
```json
{
  "id": "apt_770e8400e29b41d4a716446655440000",
  "client_id": "usr_550e8400e29b41d4a716446655440001",
  "service_id": "svc_660e8400e29b41d4a716446655440000",
  "provider_id": "usr_550e8400e29b41d4a716446655440000",
  "start_time": "2026-04-25T14:00:00Z",
  "end_time": "2026-04-25T14:45:00Z",
  "status": "pending",
  "notes": "Please arrive 10 minutes early",
  "created_at": "2026-04-20T10:20:00Z"
}
```

#### List My Appointments
```http
GET /api/v1/appointments/my
Authorization: Bearer {token}
```

**Query Parameters:**
- `status` - Filter by status: `pending`, `confirmed`, `cancelled`, `completed`
- `page` - Pagination (default: 1)
- `limit` - Results per page (default: 20)

**Example:**
```http
GET /api/v1/appointments/my?status=confirmed&page=1&limit=10
```

#### Get Appointment Details
```http
GET /api/v1/appointments/{id}
Authorization: Bearer {token}
```

#### Confirm Appointment (Provider)
```http
PATCH /api/v1/appointments/{id}/confirm
Authorization: Bearer {token}
```

**Response (200 OK)**
```json
{
  "id": "apt_770e8400e29b41d4a716446655440000",
  "status": "confirmed",
  "confirmed_at": "2026-04-20T10:25:00Z"
}
```

#### Cancel Appointment
```http
PATCH /api/v1/appointments/{id}/cancel
Authorization: Bearer {token}
Content-Type: application/json

{
  "cancellation_reason": "Personal emergency"
}
```

#### Complete Appointment (Provider)
```http
PATCH /api/v1/appointments/{id}/complete
Authorization: Bearer {token}
```

#### Get Available Slots (Public)
```http
GET /api/v1/appointments/available-slots?service_id={id}&date=2026-04-25
```

**Response (200 OK)**
```json
{
  "service_id": "svc_660e8400e29b41d4a716446655440000",
  "date": "2026-04-25",
  "available_slots": [
    "2026-04-25T09:00:00Z",
    "2026-04-25T10:00:00Z",
    "2026-04-25T14:00:00Z",
    "2026-04-25T15:00:00Z"
  ]
}
```

---

### Admin Endpoints

#### Admin Dashboard
```http
GET /api/v1/admin/dashboard
Authorization: Bearer {admin_token}
```

**Response (200 OK)**
```json
{
  "status": "success",
  "message": "Welcome to admin panel",
  "role": "admin",
  "statistics": {
    "total_users": 42,
    "total_services": 15,
    "total_appointments": 128,
    "pending_appointments": 12
  }
}
```

#### List All Users
```http
GET /api/v1/admin/users
Authorization: Bearer {admin_token}
```

**Query Parameters:**
- `role` - Filter by role: `client`, `provider`, `admin`
- `page` - Pagination (default: 1)
- `limit` - Results per page (default: 20)

**Example:**
```http
GET /api/v1/admin/users?role=provider&page=1&limit=10
```

#### List All Appointments
```http
GET /api/v1/admin/appointments
Authorization: Bearer {admin_token}
```

**Query Parameters:**
- `status` - Filter by status
- `page` - Pagination
- `limit` - Results per page

---

## ⚙️ Background Workers & Goroutines

### Architecture

The API uses concurrent goroutines for non-blocking operations:

- **Worker Pool**: 5 concurrent workers with buffer size 100
  - Handles async tasks and background jobs
  - Graceful shutdown on termination
  - Reliable task processing

- **Reminder Worker**: Appointment reminder notifications
  - Monitors upcoming appointments
  - Sends notifications before scheduled time
  - Integrated with notification system

- **Async Notifier**: Non-blocking message queue
  - Background notification processing
  - Prevents I/O blocking
  - Configurable worker count and buffer

### Key Features

✅ **Graceful Shutdown** - Workers complete tasks before exit  
✅ **Context Propagation** - Proper cancellation throughout  
✅ **Error Logging** - Failed tasks logged without stopping pool  
✅ **Concurrency Control** - WaitGroup ensures clean shutdown  

---

## 🔐 Security

### JWT Authentication

All protected endpoints require a valid JWT token:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Token Lifecycle:**
- Issued on `/api/v1/auth/login` or `/api/v1/auth/register`
- Default expiration: 24 hours
- Validated on every protected endpoint
- Refresh strategy: Re-login to get new token

### Default Admin Credentials

```
Email: admin@bookify.kz
Password: admin123
```

> ⚠️ **Change these credentials in production!**

### Role-Based Access Control (RBAC)

```
Client:      Book appointments, view own data
Provider:    Manage services, confirm appointments
Admin:       Full system access, monitoring, user management
```

---

## 📋 Testing Workflow (Postman)

### Recommended Order

#### 1. Setup & Verification
```
00 Setup → Health Check
```

#### 2. Authentication
```
01 Auth → Register Provider
01 Auth → Register Client
01 Auth → Login Provider
01 Auth → Login Client
01 Auth → Validate Token
```

#### 3. Service Management
```
02 Services → Create Service (provider token)
02 Services → List My Services
02 Services → Get Service (public)
02 Services → Update Service
```

#### 4. Client Booking
```
04 Appointments → Create Appointment (client token)
04 Appointments → List My Appointments
04 Appointments → Get Appointment
```

#### 5. Provider Actions
```
05 Appointments → Confirm Appointment (provider token)
05 Appointments → Complete Appointment
```

#### 6. Admin Operations
```
07 Admin → Login Admin
07 Admin → Admin Dashboard
07 Admin → List All Users
07 Admin → List All Appointments
```

#### 7. Error Cases (Optional)
```
06 Error Cases → Various negative tests
```

---

## 🚨 Troubleshooting

### Common Issues

#### ❌ "Connection refused" on http://localhost:8080

**Problem:** API is not running

**Solution:**
```bash
# Check if containers are running
docker ps

# Restart containers
docker compose restart

# View logs
docker compose logs -f api
```

#### ❌ "Invalid token" error (401)

**Problem:** Token expired or malformed

**Solution:**
```bash
1. Re-login to get fresh token
2. Verify token is in Authorization header
3. Format: "Authorization: Bearer {token}"
4. Check token has no extra spaces
```

#### ❌ "Unauthorized" error (403)

**Problem:** User role lacks permission

**Solution:**
- Client trying provider action? Use provider token instead
- Check user roles in database
- Verify role in JWT token matches

#### ❌ "Service not found" (404)

**Problem:** Invalid service/appointment ID

**Solution:**
```bash
1. Verify ID was saved from previous request
2. Check in Postman environment variables
3. Make sure IDs match database records
```

#### ❌ "Slot already booked" (409)

**Problem:** Appointment time slot conflict

**Solution:**
```bash
1. Check available slots first:
   GET /api/v1/appointments/available-slots?service_id={id}&date={date}
2. Use slot from available list
3. Ensure time is in future
```

#### ❌ Database connection error

**Problem:** PostgreSQL not running

**Solution:**
```bash
# Check Docker logs
docker compose logs postgres

# Restart database
docker compose restart postgres

# Verify migrations ran
docker compose exec api make migrate-up
```

---

## 📊 Example Response Formats

### Error Response (4xx/5xx)

```json
{
  "error": "Validation failed",
  "details": {
    "email": "invalid email format",
    "password": "must be at least 8 characters"
  },
  "timestamp": "2026-04-20T10:30:00Z"
}
```

### Paginated Response

```json
{
  "data": [
    { "id": "svc_1", "name": "Service 1" },
    { "id": "svc_2", "name": "Service 2" }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### List Appointments with Filters

```http
GET /api/v1/appointments/my?status=confirmed&page=1&limit=10

Response (200 OK):
{
  "data": [
    {
      "id": "apt_1",
      "status": "confirmed",
      "start_time": "2026-04-25T14:00:00Z",
      "end_time": "2026-04-25T14:45:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 23,
    "total_pages": 3
  }
}
```

---

## 🔧 HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| **200** | OK | Request successful |
| **201** | Created | Resource created |
| **204** | No Content | Deletion successful |
| **400** | Bad Request | Invalid input |
| **401** | Unauthorized | Missing/invalid token |
| **403** | Forbidden | Insufficient permissions |
| **404** | Not Found | Resource not found |
| **409** | Conflict | Slot already booked |
| **500** | Server Error | Internal error |

---

## 🔗 Useful Links

- 📖 [API Swagger Editor](https://editor.swagger.io/)
- 📮 [Postman Documentation](https://learning.postman.com/)
- 🔑 [JWT.io - Token Decoder](https://jwt.io/)
- 🗄️ [PostgreSQL Docs](https://www.postgresql.org/docs/)
- 🐹 [Go Documentation](https://golang.org/doc/)

---

## 📝 Development Setup

### Local Development

```bash
# Install Go dependencies
go mod download

# Run migrations
make migrate-up

# Seed database
make seed

# Start API (without Docker)
go run cmd/api/main.go

# Run tests
make test

# Build binary
make build
```

### Environment Variables

```env
# Server
PORT=8080
LOG_LEVEL=info
JWT_SECRET=your_secret_key_here
JWT_EXPIRATION_HOURS=24

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=bookify
DB_PASSWORD=password
DB_NAME=bookify
DB_SSL_MODE=disable

# Workers
WORKER_POOL_SIZE=5
WORKER_BUFFER_SIZE=100
```

---

## 📄 License

MIT License - see LICENSE file for details

---

## 🤝 Support

Need help?

- 📧 Check existing issues on GitHub
- 🐛 Report bugs with clear reproduction steps
- 💬 Ask questions in discussions
- 📚 Read inline code documentation

---

**Happy booking! 🎉**
