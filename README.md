# ERP Microservices with Go + Kafka + React + PostgreSQL

This repository is a learning-oriented ERP example that demonstrates a **multi-service architecture** in Go where services communicate through Kafka events and persist data in **local PostgreSQL only**.

## Services

- **department-service** (`:8081`) manages departments.
- **employee-service** (`:8082`) manages employees and consumes department events.
- **salary-service** (`:8083`) manages salaries and consumes employee events.
- **project-service** (`:8084`) manages projects/timelines and consumes department + employee events.
- **frontend** (`:5173`) React management console for CRUD-style requests.
- **postgres** (`:5432`) local system of record for all ERP entities.

Each `POST` writes to PostgreSQL and publishes an event to Kafka (topic per domain), and dependent services consume relevant topics.

## Architecture

```text
React UI --> Go REST APIs --> PostgreSQL (persistence)
                      \--> Kafka Topics --> Other Go services (event consumers)
```

## Run with Docker

```bash
docker compose up --build
```

Then open:

- Frontend: http://localhost:5173
- Department API: http://localhost:8081/departments
- Employee API: http://localhost:8082/employees
- Salary API: http://localhost:8083/salaries
- Project API: http://localhost:8084/projects
- PostgreSQL: `postgres://postgres:postgres@localhost:5432/erp?sslmode=disable`

## Local development

### Prerequisites

- Go 1.22+
- PostgreSQL running locally
- `psql` client installed (services use it internally)

### Go services

Set these variables before running services:

```bash
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/erp?sslmode=disable'
export KAFKA_REST_URL='http://localhost:8082'
```

Run services:

```bash
go run ./cmd/department-service
go run ./cmd/employee-service
go run ./cmd/salary-service
go run ./cmd/project-service
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Example payloads

### Department

```json
{ "id": "dept-001", "name": "Engineering" }
```

### Employee

```json
{ "id": "emp-001", "name": "Alice", "departmentId": "dept-001", "role": "Backend Engineer" }
```

### Salary

```json
{ "id": "sal-001", "employeeId": "emp-001", "monthly": 5800, "currency": "USD" }
```

### Project with timeline

```json
{
  "id": "prj-001",
  "name": "ERP Core Upgrade",
  "department": "Engineering",
  "members": ["emp-001"],
  "timeline": "Q2-Q4 2026",
  "startDate": "2026-04-01",
  "expectedEnd": "2026-12-15"
}
```
