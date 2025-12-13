# Service Presensi

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=flat&logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Attendance management system built with **Clean Architecture** principles. Features JWT authentication, role-based access control (RBAC), audit logging, and geofencing capabilities.

## Features

- **Authentication & Authorization**
  - JWT-based authentication with secure token handling
  - Role-based access control (Admin/User)
  - Password hashing with bcrypt
  - Login rate limiting protection

- **Attendance Management**
  - Check-in / Check-out functionality
  - CRUD operations for attendance records
  - User-specific attendance history

- **Geofencing**
  - Location-based check-in validation
  - Configurable allowed locations with radius
  - Haversine formula for accurate distance calculation

- **Audit Logging**
  - Tracks all data modifications
  - Records user actions with timestamps
  - Query by entity, user, or action type

- **Attendance Analytics**
  - Daily/monthly attendance summary
  - Per-user attendance statistics
  - Status breakdown (hadir, terlambat, izin, sakit, alpha)
  - Percentage calculations

- **Security & Performance**
  - Global rate limiting (Token Bucket algorithm)
  - Request logging and recovery middleware
  - CORS configuration

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.24 |
| Database | MongoDB |
| Auth | JWT (golang-jwt/jwt/v5) |
| Validation | go-playground/validator |
| Architecture | Hexagonal / Clean Architecture |

## Project Structure

```
service-presensi/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapter/
│   │   ├── inbound/
│   │   │   └── http/            # HTTP handlers & middleware
│   │   └── outbound/
│   │       └── mongodb/         # Database repositories
│   ├── application/
│   │   └── usecase/             # Business logic
│   ├── domain/
│   │   ├── entity/              # Domain entities
│   │   ├── repository/          # Repository interfaces
│   │   └── service/             # Domain services
│   └── infrastructure/          # Config & external services
├── pkg/                         # Shared packages
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Register new user | - |
| POST | `/api/auth/login` | User login | - |
| GET | `/api/auth/profile` | Get user profile | Required |

### Attendance (Presensi)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/presensi` | Create attendance record | Required |
| GET | `/api/presensi` | Get all attendance records | Required |
| GET | `/api/presensi/{id}` | Get attendance by ID | Required |
| PUT | `/api/presensi/{id}` | Update attendance | Required |
| DELETE | `/api/presensi/{id}` | Delete attendance | Admin |
| POST | `/api/presensi/{id}/checkin` | Check-in with location | Required |
| POST | `/api/presensi/{id}/checkout` | Check-out | Required |

### Locations (Geofencing)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/locations` | Create allowed location | Admin |
| GET | `/api/locations` | Get all locations | Admin |
| GET | `/api/locations/{id}` | Get location by ID | Admin |
| PUT | `/api/locations/{id}` | Update location | Admin |
| DELETE | `/api/locations/{id}` | Delete location | Admin |

### Audit Logs
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/audit` | Get all audit logs | Admin |
| GET | `/api/audit/{id}` | Get audit log by ID | Admin |
| GET | `/api/audit/entity` | Get logs by entity | Admin |
| GET | `/api/audit/user/{user_id}` | Get logs by user | Admin |

### Analytics
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/analytics/summary` | Overall attendance summary | Required |
| GET | `/api/analytics/daily?date=YYYY-MM-DD` | Daily summary | Required |
| GET | `/api/analytics/monthly?month=YYYY-MM` | Monthly summary | Required |
| GET | `/api/analytics/user/{user_id}` | User attendance statistics | Required |
| GET | `/api/analytics/status-breakdown` | Status distribution | Required |

## Quick Start

### Prerequisites
- Go 1.24+
- MongoDB 6.0+
- Docker (optional)

### Environment Variables

Create `.env` file:

```env
# Server
PORT=8080

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=presensi

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRY=24h

# Geofencing (optional)
GEOFENCING_ENABLED=true
DEFAULT_RADIUS_METERS=100
```

### Run Locally

```bash
# Clone repository
git clone https://github.com/bimakw/service-presensi.git
cd service-presensi

# Install dependencies
go mod download

# Run application
go run cmd/api/main.go
```

### Run with Docker

```bash
# Build and run
docker-compose up -d

# Check logs
docker-compose logs -f
```

## Usage Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Check-in with Location
```bash
curl -X POST http://localhost:8080/api/presensi/{id}/checkin \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": -6.2088,
    "longitude": 106.8456
  }'
```

### Get Monthly Analytics
```bash
curl -X GET "http://localhost:8080/api/analytics/monthly?month=2024-01" \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "success": true,
  "message": "Berhasil mengambil summary bulanan",
  "data": {
    "month": "2024-01",
    "summary": {
      "total_records": 150,
      "total_hadir": 100,
      "total_terlambat": 20,
      "total_izin": 10,
      "total_sakit": 15,
      "total_alpha": 5,
      "percentage_hadir": 80.0
    },
    "daily_stats": [
      {"date": "2024-01-01", "count": 5},
      {"date": "2024-01-02", "count": 8}
    ]
  }
}
```

## Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters):

```
┌─────────────────────────────────────────────────────────────┐
│                      Adapters (Inbound)                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ HTTP Handler│  │ Middleware  │  │ Rate Limiter        │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
└─────────┼────────────────┼────────────────────┼─────────────┘
          │                │                    │
          ▼                ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    Use Cases                         │    │
│  │  • AuthUseCase    • PresensiUseCase                  │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
          │                │                    │
          ▼                ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────────────┐     │
│  │  Entities  │  │ Repository │  │  Domain Services   │     │
│  │            │  │ Interfaces │  │  (LocationService) │     │
│  └────────────┘  └────────────┘  └────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
          │                │                    │
          ▼                ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                    Adapters (Outbound)                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              MongoDB Repositories                    │    │
│  │  • UserRepo  • PresensiRepo  • AuditRepo  • Location │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## License

This project is licensed under the MIT License with Attribution Requirement - see the [LICENSE](LICENSE) file for details.

## Author

**Bima Kharisma Wicaksana**
- GitHub: [@bimakw](https://github.com/bimakw)
- LinkedIn: [Bima Kharisma Wicaksana](https://www.linkedin.com/in/bima-kharisma-wicaksana-aa3981153/)
