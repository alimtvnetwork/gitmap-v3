# 11 — Database Patterns

Schema design, migrations, query optimization, and connection management.

## Schema Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Table names | PascalCase | `UserRoles` |
| Column names | PascalCase | `CreatedAt` |
| Primary key | `INTEGER PRIMARY KEY AUTOINCREMENT` | `Id` |
| Strings | `TEXT DEFAULT ''` | |
| Booleans | `INTEGER DEFAULT 0` | |
| Timestamps | `TEXT DEFAULT CURRENT_TIMESTAMP` | |

Every table includes `Id` and `CreatedAt`.

## Migrations

Forward-only, append-only. Prefer additive changes (add columns with defaults). For breaking changes: detect legacy schema, provide recovery instructions.

## Query Rules

- Batch over loop (no N+1).
- Upsert with `INSERT ... ON CONFLICT ... DO UPDATE`.
- No `SELECT *` — list columns explicitly.
- Parameterized queries only — no string interpolation.
- Wrap multi-statement writes in transactions.

## Connection Management

CLI tools: single connection, `SetMaxOpenConns(1)`. Enable WAL mode and foreign keys. Always `defer rows.Close()`.

## Constraints

All SQL strings in the `constants` package. SQLite driver: `modernc.org/sqlite` (CGo-free). No UUID primary keys.

---

Source: `spec/05-coding-guidelines/11-database-patterns.md`
