# backend

Go backend for note-thing.

## Running in dev with hot reload

Install Air:

```bash
go install github.com/air-verse/air@latest
```

From `backend/`:

```bash
air
```

The server listens on port `18611` by default. Change with `PORT`.

## Running migrations

Set `DATABASE_URL` to your Postgres connection string, then from `backend/`:

```bash
go run ./cmd/migrate -direction up
go run ./cmd/migrate -direction down -steps 1
```

This uses `golang-migrate/migrate` ([repo](https://github.com/golang-migrate/migrate)).
