# Go + Convex Todo API

A simple Todo backend using:
- **Go + Fiber** for REST endpoints
- **Convex** for data storage and mutations

This project is a reference starter for integrating Convex into Go projects.

## Project Structure

- `main.go` — Fiber server and HTTP routes
- `convex/schema.ts` — Convex table schema (`todoList`)
- `convex/todoActions.ts` — Convex query/mutations used by the Go API

## Prerequisites

- Go (matching `go.mod`, currently `1.25`)
- Node.js + npm
- Convex account/project initialized

## Setup

1. Install Go dependencies:

```bash
go mod tidy
```

2. Install Node dependencies:

```bash
npm install
```

3. Create `.env` in the project root:

```env
PORT=8080
CONVEX_URL=https://<your-deployment>.convex.cloud
CONVEX_ADMIN_KEY=<your-convex-admin-key>
```

> `CONVEX_URL` is required.

## Run the Project

### 1) Run Convex dev

```bash
npx convex dev
```

### 2) Run the Go API

Option A:

```bash
go run main.go
```

Option B (auto-reload with Air):

```bash
air
```

Server starts at `http://localhost:$PORT`.

## API Endpoints

### Health
- `GET /` → `{ "msg": "hello world" }`

### Todos
- `GET /todos` — list todos
- `POST /add` — create todo
- `PATCH /update/:id` — mark todo as completed
- `PATCH /update/body/:id` — update todo body
- `DELETE /delete/:id` — delete todo

## Example Requests

Create todo:

```bash
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"body":"Buy milk"}'
```

Update todo body:

```bash
curl -X PATCH http://localhost:8080/update/body/<todo_id> \
  -H "Content-Type: application/json" \
  -d '{"body":"Buy milk and bread"}'
```

## Convex Function Mapping

Go route → Convex function path:

- `GET /todos` → `todoActions:list` (query)
- `POST /add` → `todoActions:add` (mutation)
- `PATCH /update/:id` → `todoActions:setCompleted` (mutation)
- `PATCH /update/body/:id` → `todoActions:setBody` (mutation)
- `DELETE /delete/:id` → `todoActions:remove` (mutation)

## Notes

- Table name is `todoList` in `convex/schema.ts` and must match usage in `convex/todoActions.ts`.
- If Convex returns `convex error`, verify function names and argument shapes match between `main.go` and `convex/todoActions.ts`.
