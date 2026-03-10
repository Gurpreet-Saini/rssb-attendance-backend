# Multi-Center Attendance Management Backend

Production-ready attendance management API built with Golang for multi-center operations.

## Stack

- Backend: Go, Gin, GORM, PostgreSQL, JWT, Excelize
- DevOps: Docker, Docker Compose, Makefile, environment-based config

## Features

- JWT authentication with `super_admin`, `center_admin`, and `operator` roles
- Multi-center employee and attendance management
- Attendance via name, employee ID, badge number, or barcode scan
- Center-restricted access controls for admins and operators
- Dashboard with attendance KPIs and charts
- Excel export for attendance reports
- Rate limiting, bcrypt hashing, and audit log persistence

## Project Structure

```text
backend-repo/
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
```

## Quick Start

### Docker

```bash
docker compose up --build
```

Backend API: [http://localhost:8081/api](http://localhost:8081/api)

Default super admin credentials:

- Email: `admin@example.com`
- Password: `Admin@123`

### Local Development

1. Backend

```bash
cp .env.example .env
go mod tidy
go run main.go
```

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

## Notes

- SQL migrations are provided in [001_init.sql](/Users/gurpreetsaini/Documents/Playground/backend-repo/migrations/001_init.sql).
- The backend also runs `AutoMigrate` at startup to simplify first boot in containerized environments.
- Set `CORS_ALLOWED_ORIGINS` to a comma-separated list of allowed frontend origins in environments where browsers will call this API directly.
- Operators can mark attendance and search employees, but cannot create or update employees.

## Make Targets

```bash
make install-backend
make build-backend
make up
make down
```
