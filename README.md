# Movies API

A REST API for managing a movie database — movies, genres, and actors, with many-to-many
relationships between them. Built in Go using the standard `net/http` package and raw SQL
against SQLite (no ORM).

## Project overview

The API is organized in three layers:

- **`handlers/`** — parses HTTP requests, decodes/validates JSON bodies, and writes responses.
- **`service/`** — business logic: orchestrates repository calls, translates database errors
  into typed API errors, enforces rules like "don't delete a genre that's still linked to movies
  unless `force=true`".
- **`repository/`** — data access. Raw SQL via `database/sql`, one file per entity
  (`movieRepo.go`, `actorRepo.go`, `genreRepo.go`).

Supporting packages:

- **`models/`** — request/response structs, including separate `*Patch` structs (e.g.
  `MoviePatch`) with pointer fields so PATCH requests can distinguish "field omitted" from
  "field explicitly cleared."
- **`errors/`** — a single `AppError` type (`code`, `message`, HTTP status) with constructors
  for the app's error cases, plus a `WriteError` helper that all handlers funnel through.
- **`database/`** — SQLite connection setup and schema creation (`CREATE TABLE IF NOT EXISTS`).
- **`routes/`** — maps HTTP method + path patterns to handler functions using Go's
  `http.ServeMux` (Go 1.22+ pattern matching, e.g. `GET /api/movies/{id}`).

### Data model

- **Genre**: `id`, `name`
- **Movie**: `id`, `title`, `releaseYear`, `duration`
- **Actor**: `id`, `name`, `birthDate` (ISO 8601, `YYYY-MM-DD`)
- **MovieGenre** / **MovieActor**: join tables implementing the many-to-many relationships
  (a movie can have multiple genres and multiple actors; a genre or actor can belong to
  multiple movies).

All `id` fields are auto-incrementing primary keys and are never modified after creation.

## Setup and installation

**Requirements:** Go 1.25+, a C toolchain (needed by `go-sqlite3`, which uses cgo).

\`\`\`bash
git clone <repo-url>
cd movies-api
go mod tidy
go run main.go
\`\`\`

The server starts on `http://localhost:8080` and opens (creating if necessary)
a SQLite database file at `./movies.db` in the project root. Tables are created
automatically on startup if they don't already exist — no separate migration step.

### Dependencies

- [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) — SQLite driver
- [`github.com/go-playground/validator/v10`](https://github.com/go-playground/validator) —
  struct-tag based input validation

## Usage guide

All request/response bodies are JSON. All endpoints are prefixed with `/api`.

### Movies

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/movies` | List all movies. Supports filters (see below). |
| `GET` | `/api/movies/{id}` | Get one movie by id. |
| `POST` | `/api/movies` | Create a movie, optionally with nested `actors`/`genres`. |
| `PATCH` | `/api/movies/{id}` | Partially update a movie (title, releaseYear, duration, and/or its genre/actor associations by id). |
| `DELETE` | `/api/movies/{id}` | Delete a movie. Blocked if it has linked genres/actors unless `?force=true`. |
| `GET` | `/api/movies/{id}/actors` | List all actors in a movie. |
| `GET` | `/api/movies/{id}/genres` | List all genres of a movie. |

**Movie filters** (query params on `GET /api/movies`, mutually exclusive):

- `?genre={genreId}` — movies in a given genre
- `?year={releaseYear}` — movies released in a given year
- `?actor={actorId}` — movies a given actor starred in

**Create a movie:**
\`\`\`json
POST /api/movies
{
  "title": "Inception",
  "releaseYear": 2010,
  "duration": 148,
  "genres": [{ "name": "Sci-Fi" }, { "name": "Action" }],
  "actors": [{ "name": "Leonardo DiCaprio", "birthDate": "1974-11-11" }]
}
\`\`\`
Genres and actors are created if they don't already exist (matched by name, or name +
birthDate for actors) and reused if they do — you can freely add an existing actor to a
new movie.

**Patch a movie** (any subset of fields; omitted fields are left unchanged):
\`\`\`json
PATCH /api/movies/1
{
  "title": "Inception (Director's Cut)",
  "genres": [1, 3],
  "actors": [2, 5]
}
\`\`\`
`genres`/`actors` in a PATCH body are **ids** of existing genres/actors (unlike the nested
objects used at creation time), and replace the movie's full set of associations — sending
`"genres": []` clears all of a movie's genres. Referencing a genre/actor id that doesn't
exist returns `404`.

### Actors

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/actors` | List all actors. Supports `?name=` filter (partial, case-insensitive). |
| `GET` | `/api/actors/{id}` | Get one actor by id, plus the movies they're in. |
| `POST` | `/api/actors` | Create an actor. |
| `PATCH` | `/api/actors/{id}` | Partially update an actor's name and/or birthDate. |
| `DELETE` | `/api/actors/{id}` | Delete an actor. Blocked if cast in any movie unless `?force=true`. |

### Genres

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/genres` | List all genres. |
| `GET` | `/api/genres/{id}` | Get one genre by id. |
| `POST` | `/api/genres` | Create a genre. |
| `PATCH` | `/api/genres/{id}` | Partially update a genre's name. |
| `DELETE` | `/api/genres/{id}` | Delete a genre. Blocked if linked to any movie unless `?force=true`. |

### Deleting entities with existing relationships

By default, deleting a genre or actor that's still linked to one or more movies is refused:

\`\`\`
DELETE /api/genres/1

HTTP/1.1 400 Bad Request
{
  "code": "GENRE_IN_USE",
  "message": "Cannot delete genre 'Sci-Fi' because it has 30 associated movies"
}
\`\`\`

Pass `force=true` to delete the entity and remove its associations in the same operation:

\`\`\`
DELETE /api/genres/1?force=true

HTTP/1.1 204 No Content
\`\`\`

The same pattern applies to actors (`ACTOR_IN_USE`) and movies (`MOVIE_IN_USE`, which
reports both linked actor and genre counts). This is enforced atomically — the link-count
check, unlinking, and deletion all happen inside a single database transaction, so a failure
partway through leaves nothing partially deleted.

### Error responses

All errors share one JSON shape:

\`\`\`json
{ "code": "MOVIE_NOT_FOUND", "message": "movie not found" }
\`\`\`

with a matching HTTP status code (`400` for bad input/validation/in-use conflicts, `404` for
not-found, `409` for duplicate/unique-constraint conflicts, `500` for unexpected server
errors). Request bodies are validated using struct tags
(`github.com/go-playground/validator`) before touching the database — e.g. a movie's
`title` is required, `duration` must be positive, and an actor's `birthDate` must match
`YYYY-MM-DD`.

## Not yet implemented

Called out explicitly rather than left silently missing:

- **Pagination** (`?page=`/`?size=` on list endpoints) — not implemented; list endpoints
  currently return the full result set.
- **Dedicated title search** (`GET /api/movies/search?title=`) — not implemented as its own
  endpoint. Filtering by genre, release year, and actor id is implemented on
  `GET /api/movies`.
- **Seed/sample data** — the database starts empty; there's no bundled seed script yet, so
  sample genres/movies/actors need to be created manually (e.g. via Postman) after first run.

## Testing

The API was tested manually with Postman — a collection covering CRUD for all three
entities, the filter query parameters, PATCH partial-update behavior, and the
delete-with-relationships / `force=true` error cases.
