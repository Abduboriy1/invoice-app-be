# Invoice Backend

A Go backend service for invoice management with Jira and Square integrations.

## Project Structure

```
invoice-app-be/
├── cmd/                    # Application entry points
│   ├── api/               # HTTP API server
│   └── worker/            # Background job processor
├── internal/              # Private application code
│   ├── domain/            # Business entities & rules
│   ├── infrastructure/    # External concerns (DB, APIs)
│   ├── interfaces/        # Delivery mechanisms (HTTP, jobs)
│   └── pkg/               # Internal shared packages
├── config/                # Configuration
├── migrations/            # Database migrations
└── scripts/               # Development scripts
```

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL 16+

## Getting Started

1. Clone the repository
2. Run setup script:
   ```bash
   make dev-setup
   ```

3. Start the API server:
   ```bash
   make run-api
   ```

Or to hot refresh run:

```bash
   air
   ```

## Development

### Running with Docker

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Running Tests

```bash
make test
```

### Running Migrations

```bash
make migrate
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login

### Invoices

- `GET /api/invoices` - List invoices
- `POST /api/invoices` - Create invoice
- `GET /api/invoices/{id}` - Get invoice
- `PUT /api/invoices/{id}` - Update invoice
- `DELETE /api/invoices/{id}` - Delete invoice

### Time Entries

- `GET /api/time-entries` - List time entries
- `POST /api/time-entries` - Create time entry
- `GET /api/time-entries/{id}` - Get time entry
- `PUT /api/time-entries/{id}` - Update time entry
- `DELETE /api/time-entries/{id}` - Delete time entry

## Environment Variables

| Variable            | Description         | Default   |
|---------------------|---------------------|-----------|
| SERVER_PORT         | API server port     | 8080      |
| DB_HOST             | PostgreSQL host     | localhost |
| DB_PORT             | PostgreSQL port     | 5432      |
| DB_USER             | PostgreSQL user     | postgres  |
| DB_PASSWORD         | PostgreSQL password | -         |
| DB_NAME             | PostgreSQL database | invoice   |
| JWT_SECRET          | JWT signing secret  | -         |
| JIRA_BASE_URL       | Jira instance URL   | -         |
| JIRA_API_TOKEN      | Jira API token      | -         |
| SQUARE_ACCESS_TOKEN | Square API token    | -         |

## License

MIT
