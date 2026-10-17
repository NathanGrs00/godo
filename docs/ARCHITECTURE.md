# Architecture

## Overview

```
 Browser (laptop)            Browser (phone)
 ┌──────────────┐            ┌──────────────┐
 │ React app    │            │ React app    │
 │ DataProvider │            │ DataProvider │
 └──┬────────▲──┘            └──┬────────▲──┘
    │REST    │SSE /events       │REST    │SSE /events
 ┌──▼────────┴───────────────────▼────────┴──┐
 │ Go API (Gin)                              │
 │  handlers → services → repositories       │
 │                 │                         │
 │                 └─► realtime.Hub (per user)│
 └──────────────────────┬────────────────────┘
                        │
                    MongoDB
```

## Backend layers

Each layer only depends on the one below it, through interfaces, so every
layer can be tested on its own.

| Layer | Package | Responsibility |
|---|---|---|
| HTTP | `internal/api`, `api/handlers`, `api/middleware` | Routing, auth check, JSON in/out, mapping errors to status codes. No business logic. |
| Service | `internal/service` | Validation, ownership, linking tasks to tags/deadlines, cascading deletes, publishing change events. |
| Repository | `internal/repository` (+ `mongo`, `memory`) | Storing and loading documents. Every query is scoped to a user ID. |
| Realtime | `internal/realtime` | Fan-out of change events to a user's open streams. |
| Auth | `internal/auth` | Password hashing (bcrypt) and login tokens (JWT, HS256). |

`ResourceHandler` is a single generic CRUD handler used for tasks, deadlines
and tags, and `repository.Store[T]` is a single generic storage interface;
this keeps the three resources consistent and removes duplicated code.

### Data model

Every user-owned document embeds `models.Meta`:

| Field | Purpose |
|---|---|
| `id` | ObjectID |
| `user_id` | Owner (never sent to clients) |
| `version` | Starts at 1, +1 on every write |
| `created_at`, `updated_at` | Timestamps (UTC) |

Relationships:

- A **task** stores `deadline_ids` and `tag_ids`.
- A **deadline** stores `tag_ids`. Its tasks are *derived* (all tasks whose
  `deadline_ids` contain it), so there is one source of truth and the two
  sides can never disagree.
- `deadline.passed` is computed when read, not stored, so it is always correct.
- Deleting a tag or deadline removes its ID from every linked item.

## How sync between devices works

1. **Identity.** Each device logs in to the same account and gets its own
   token. All data is stored per user on the server.
2. **Live updates.** After logging in, the app opens `GET /events`, a
   Server-Sent Events stream. When any device changes something, the service
   publishes an event (`task.updated`, `tag.deleted`, …) to the user's hub,
   which forwards it to every open stream of that user. The receiving app
   applies it to its state straight away.
3. **Catching up.** Each (re)connection starts with an `event: ready`
   message. The app reloads everything on `ready`, so changes made while a
   device was offline or asleep are picked up. The stream reconnects
   automatically with exponential backoff (1 s → 30 s) and immediately when
   the browser comes back online. Load happens *after* subscribing, so no
   change can slip through between the two.
4. **Ordering.** The client only replaces an item with an equal or newer
   `version`, so a late or duplicate event can never roll back newer data.
5. **Conflicts.** Updates send the `version` the client is editing. The
   database only accepts the write if that is still the stored version
   (optimistic concurrency); otherwise the API returns **409 Conflict**.
   The app then reloads and tells the user the item was changed on another
   device; saving again keeps their version, cancelling keeps the other.
6. **Slow clients.** If a stream falls more than 64 events behind, the hub
   closes it rather than blocking or silently dropping events. The client
   reconnects and reloads, so it never ends up with a gap.

### Why Server-Sent Events?

Sync only needs server → client pushes (writes already go through REST), so
SSE is simpler than WebSockets: plain HTTP, works through proxies, and is
easy to test. The frontend reads it with `fetch` instead of `EventSource`
so it can send the `Authorization` header rather than putting the token in
the URL (where it would end up in logs).

### Limitations and next steps

- The hub lives in the server's memory, so all devices must reach the same
  backend instance. To run several instances, publish events through a
  shared broker instead (MongoDB change streams with a replica set, or
  Redis pub/sub) — only `realtime.Hub` needs replacing, since services
  depend on the small `service.Publisher` interface.
- There is no offline editing; when the server is unreachable, changes
  cannot be saved (the app shows "Offline – retrying").
- Login tokens cannot be revoked before they expire (`TOKEN_TTL`).

## Frontend structure

| Folder | Responsibility |
|---|---|
| `api/client.ts` | `fetch` wrapper: base URL, token header, JSON, `ApiError` with status and field |
| `services/` | One small module per resource (`taskService`, …) |
| `sync/eventStream.ts` | Incremental SSE parser and reconnecting stream |
| `sync/dataReducer.ts` | Pure reducer: applies every change (local or remote) by version |
| `context/AuthProvider.tsx` | Login state, token storage, auto-logout on 401 |
| `context/DataProvider.tsx` | Loads data, subscribes to the stream, exposes actions; reloads on 409 |
| `components/` | UI; components use `useData()` and never call `fetch` directly |
| `utils/` | Validation (mirrors the server rules), date helpers, search |

Local changes and remote events go through the same reducer path, so the
screen shows the same result whichever device made the change.
