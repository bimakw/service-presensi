# Service Presensi

[![CI](https://img.shields.io/github/actions/workflow/status/bimakw/service-presensi/ci.yml?branch=main)](https://github.com/bimakw/service-presensi/actions)

Attendance management API with JWT auth, RBAC, geofencing, and audit trails. Built with hexagonal architecture in Go + MongoDB.

## Running

```bash
cp .env.example .env  # configure MongoDB URI, JWT secret, etc.
go run cmd/api/main.go
```

Or with Docker:

```bash
docker-compose up -d
```

### Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection |
| `MONGO_DATABASE` | `presensi` | Database name |
| `JWT_SECRET` | - | Token signing key |
| `JWT_EXPIRE_MINUTES` | `1440` | Token TTL |
| `GEOFENCE_ENABLED` | `false` | Enable location-based check-in |
| `DEFAULT_RADIUS_METERS` | `100` | Geofence radius |

## API

Auth endpoints (`/api/auth/*`) — register, login, profile, change password.

Attendance (`/api/presensi/*`) — CRUD + check-in/check-out with optional geofencing.

Locations (`/api/locations/*`) — manage allowed check-in locations (admin).

Audit (`/api/audit/*`) — query audit trail by entity, user, or action (admin).

Analytics (`/api/analytics/*`) — daily/monthly summaries, per-user stats, status breakdown.

## License

MIT — see [LICENSE](LICENSE).
