# Multi-Center Attendance Management Backend

Production-ready attendance management API built with Golang for multi-center operations.

## Stack

- Backend: Go, Gin, GORM, PostgreSQL, JWT, Excelize
- DevOps: Docker, Docker Compose, GitHub Actions, Railway, Makefile

## Features

- JWT authentication with `super_admin`, `center_admin`, and `operator` roles
- Multi-center employee and attendance management
- Attendance via name, employee ID, badge number, or barcode scan
- Center-restricted access controls for admins and operators
- Dashboard with attendance KPIs and charts
- Excel export for attendance reports
- Rate limiting, bcrypt hashing, and audit log persistence
- Automated CI/CD with GitHub Actions

## Project Structure

```text
backend-repo/
  .github/
    workflows/
      ci-cd.yml       # GitHub Actions CI/CD pipeline
  config/
  controllers/
  middleware/
  migrations/
  models/
  repositories/
  routes/
  services/
  utils/
  main.go
  Dockerfile
  docker-compose.yml
```

## Quick Start

### Local Development

1. Copy environment template

```bash
cp .env.example .env
# Edit .env with your database credentials
```

2. Run the application

```bash
go mod tidy
go run main.go
```

### Docker

```bash
docker compose up --build
```

Backend API: [http://localhost:8081/api](http://localhost:8081/api)

Default super admin credentials:

- Email: `admin@example.com`
- Password: `Admin@123`

## Deployment

### Railway (Recommended)

The app is configured for automatic deployment to Railway via GitHub Actions.

1. Create a Railway project
2. Add these GitHub Secrets:
   - `RAILWAY_TOKEN` - Railway account token
   - `RAILWAY_PROJECT_ID` - Your project ID
   - `RAILWAY_SERVICE_ID` - Your service ID
   - `DATABASE_URL` - PostgreSQL connection string
   - `JWT_SECRET` - Secure random string
3. Push to `main` branch - deployment starts automatically

### Manual Docker Deployment

```bash
# Build image
docker build -t attendance-api .

# Run container
docker run -p 8080:8080 -e DATABASE_URL="..." -e JWT_SECRET="..." attendance-api
```

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/ci-cd.yml`) automatically:

1. **On Push to Main:**
   - Builds Docker image
   - Pushes to GitHub Container Registry (ghcr.io)
   - Deploys to Railway

2. **On Pull Request:**
   - Builds Docker image for testing
   - Does not push or deploy

## API Summary

- `POST /api/auth/login`
- `GET /api/dashboard/summary`
- `GET /api/centers`
- `POST /api/centers`
- `GET /api/users`
- `POST /api/users`
- `GET /api/employees`
- `POST /api/employees`
- `PUT /api/employees/:id`
- `DELETE /api/employees/:id`
- `POST /api/attendance/checkin`
- `POST /api/attendance/checkout`
- `POST /api/attendance/scan`
- `GET /api/attendance`
- `GET /api/reports/excel`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `APP_ENV` | Environment (development/production) |
| `PORT` | Server port (default: 8080) |
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for JWT signing |
| `JWT_DURATION_HOURS` | Token validity in hours |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins |

## Notes

- SQL migrations are provided in `migrations/001_init.sql`.
- The backend also runs `AutoMigrate` at startup to simplify first boot in containerized environments.
- Set `CORS_ALLOWED_ORIGINS` to a comma-separated list of allowed frontend origins in environments where browsers will call this API directly.
- Operators can mark attendance and search employees, but cannot create or update employees.
- **Never commit `.env` file** - it's already in `.gitignore` to prevent accidental commits of secrets.

## Make Targets

```bash
make install-backend
make build-backend
make up
make down
```
