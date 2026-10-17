# HTTP API

Base URL: `http://localhost:8080` (configurable with `PORT`).

All request and response bodies are JSON. Every endpoint except `/health`,
`/auth/register` and `/auth/login` requires:

```
Authorization: Bearer <token>
```

## Responses

- Success bodies wrap data in a named key: `{"task": {...}}`, `{"tasks": [...]}`.
- `DELETE` returns `204 No Content`.
- Errors are `{"error": "message"}`; validation errors add `"field"`.

| Status | Meaning |
|---|---|
| 400 | Invalid body, invalid ID, or validation error (`field` names the field) |
| 401 | Missing/invalid/expired token, or wrong email/password |
| 404 | Item does not exist or belongs to another user |
| 409 | Email already registered, or the item was changed since the given `version` |

## Auth

### `POST /auth/register`

```json
{ "email": "you@example.com", "password": "at least 8 chars" }
```

`201` → `{"token": "...", "user": {"id", "email", "created_at"}}`

### `POST /auth/login`

Same body. `200` → same response. Log in on each device separately.

### `GET /auth/me`

`200` → `{"user": {...}}`

## Tasks

| Method | Path | Body | Response |
|---|---|---|---|
| GET | `/tasks` | – | `{"tasks": [Task]}` (oldest first) |
| POST | `/tasks` | TaskInput | `201 {"task": Task}` |
| GET | `/tasks/:id` | – | `{"task": Task}` |
| PUT | `/tasks/:id` | TaskInput incl. `version` | `{"task": Task}` |
| DELETE | `/tasks/:id` | – | `204` |

**TaskInput**

| Field | Type | Rules |
|---|---|---|
| `title` | string | required, ≤ 200 chars (trimmed) |
| `description` | string | ≤ 2000 chars |
| `priority` | `"low"` \| `"medium"` \| `"high"` | required |
| `completed` | bool | |
| `deadline_ids` | string[] | must be your deadlines; duplicates removed |
| `tag_ids` | string[] | must be your tags; duplicates removed |
| `version` | number | **required for PUT**: the version you are editing |

**Task** = TaskInput fields + `id`, `version`, `created_at`, `updated_at`.

`PUT` replaces all fields, so send the full task.

## Deadlines

Same five routes under `/deadlines`, wrapped in `deadline` / `deadlines`.

**DeadlineInput**

| Field | Type | Rules |
|---|---|---|
| `title` | string | required, ≤ 200 chars |
| `description` | string | ≤ 2000 chars |
| `date` | string (RFC 3339) | required; stored in UTC |
| `tag_ids` | string[] | must be your tags |
| `version` | number | required for PUT |

**Deadline** also has `passed` (bool, computed: date is in the past).
Deleting a deadline unlinks it from all tasks.

## Tags

Same five routes under `/tags`, wrapped in `tag` / `tags`.

**TagInput**

| Field | Type | Rules |
|---|---|---|
| `title` | string | required, ≤ 50 chars |
| `color` | string | hex colour `#RRGGBB` (stored upper case) |
| `version` | number | required for PUT |

Deleting a tag unlinks it from all tasks and deadlines.

## Live events

### `GET /events`

A [Server-Sent Events](https://html.spec.whatwg.org/multipage/server-sent-events.html)
stream of changes to the user's data, from any device.

```
event: ready
data: {}

data: {"type":"task.created","id":"65f...","data":{...task...}}

data: {"type":"tag.deleted","id":"65f..."}

: ping
```

- `ready` is sent once the subscription is live. Reload your data then.
- Event types: `task.created`, `task.updated`, `task.deleted`, and the same
  for `deadline.*` and `tag.*`. `data` holds the item after the change and
  is absent for deletions.
- Deleting a tag or deadline is followed by `*.updated` events for every
  item it was unlinked from.
- `: ping` comments keep the connection open (every 25 s).
- If the server closes the stream, reconnect and reload.

## Health

### `GET /health`

`200` → `{"status": "ok"}` (no auth).
