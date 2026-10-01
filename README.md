# PlayLedger

PlayLedger is an attendance and billing platform for recurring pickleball sessions. It records which players participated in each time slot, calculates individual charges from actual participation, and maintains a running ledger of billed amounts, payments and outstanding balances.

Court bookings typically span several hours, while players arrive and leave at different times. Rather than splitting the total cost evenly across everyone who attended, PlayLedger bills each player only for the time slots they participated in, at their proportional share of the courts booked for those slots. Unpaid balances carry forward automatically to subsequent sessions.

Built with **Go (Gin, GORM, PostgreSQL)** and **React (TypeScript, Vite)**.

## Features

- **Session and time slot management.** Sessions can be divided into hourly slots automatically, with the number of courts configured per slot.
- **Slot-level attendance.** An attendance grid (players × time slots) records participation for each slot.
- **Proportional billing.** Each slot's court cost is divided evenly among the players in that slot.
- **Payments and ledger.** Payments are recorded per player, with balances carried forward across sessions.
- **Player history.** Each player has a per-session breakdown of charges and payments, with a WhatsApp link for payment reminders.
- **Authentication.** JWT-based authentication with token expiry and configurable self-registration.
- **API documentation.** OpenAPI specification served through Swagger UI at `/docs`.

## Billing model

For every time slot:

```
slot cost       = courts booked × court price per hour
cost per player = slot cost ÷ players in that slot
```

A player's charge for a session is the sum of their charges across all slots they participated in. Slots with no players are not billed.

**Example** (court price ₹300 per hour):

| Time slot   | Courts booked | Players in slot | Charge per player |
| ----------- | ------------- | --------------- | ----------------- |
| 7pm – 8pm   | 2             | 3               | ₹200              |
| 8pm – 9pm   | 1             | 2               | ₹150              |

A player present in both slots is billed ₹350; a player present only in the first slot is billed ₹200.

Payments are recorded against a player and session rather than individual time slots. Each player's balance is calculated as:

```
balance = total paid − total billed      (negative = still owes)
```

A player billed ₹350 who pays ₹300 carries a balance of −₹50 into subsequent sessions.

## Project structure

```
playledger/
├── cmd/api/main.go            # entry point and routes
├── internal/
│   ├── config/                # env loading and database connection
│   ├── domain/                # models and request types
│   ├── services/              # business logic: billing, ledger, sessions, auth (+ tests)
│   ├── handlers/              # HTTP layer only
│   ├── middleware/            # JWT authentication
│   └── docs/                  # OpenAPI spec and Swagger UI
└── frontend/                  # React app
    └── src/
        ├── api/               # API client and typed endpoint functions
        ├── pages/             # Dashboard, Sessions, Session detail, Players, ...
        ├── components/
        └── lib/               # formatting and helpers
```

Handlers deal with HTTP; anything that decides money or business rules lives in `services`, where it can be unit tested without a database.

## Running locally

**Requirements:** Go 1.26+, Node 22+, PostgreSQL.

### 1. Database

```bash
createdb playledger
```

Tables are created automatically when the API starts.

### 2. API

```bash
cp .env.example .env      # then set DB_PASSWORD and JWT_SECRET
go run ./cmd/api
```

The API listens on `http://localhost:8080`. Check it with `curl localhost:8080/health`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173` and create an account. In development, Vite forwards `/api/*` to the Go server, so no CORS setup is needed.

### Tests

```bash
go test ./...
cd frontend && npm run lint && npm run build
```

## Configuration

All settings come from environment variables, or from a `.env` file in the directory you start the API from. See [`.env.example`](.env.example).

| Variable             | Default                  | Notes                                             |
| -------------------- | ------------------------ | ------------------------------------------------- |
| `DB_HOST`            | `localhost`              |                                                   |
| `DB_PORT`            | `5432`                   |                                                   |
| `DB_USER`            | `postgres`               |                                                   |
| `DB_PASSWORD`        | (empty)                  |                                                   |
| `DB_NAME`            | `playledger`             |                                                   |
| `DB_SSLMODE`         | `disable`                | Use `require` for hosted databases                |
| `SERVER_PORT`        | `8080`                   |                                                   |
| `JWT_SECRET`         | (required)               | The API refuses to start without it               |
| `JWT_TTL_HOURS`      | `168`                    | How long a login lasts (7 days)                   |
| `CORS_ORIGINS`       | `http://localhost:5173`  | Comma-separated                                   |
| `ALLOW_REGISTRATION` | `true`                   | Set to `false` once your own account exists       |

## API

Full interactive docs are served at **`http://localhost:8080/docs`**.
Everything except `/health`, `/register`, `/login` and `/docs` needs an `Authorization: Bearer <token>` header from `/login`.

| Method | Path                                  | Description                                    |
| ------ | ------------------------------------- | ---------------------------------------------- |
| POST   | `/register`                           | Create an account                              |
| POST   | `/login`                              | Get a JWT                                      |
| GET    | `/sessions`                           | List sessions with slots and players           |
| POST   | `/sessions`                           | Create a session (optionally with hourly slots)|
| GET    | `/sessions/:id`                       | One session                                    |
| DELETE | `/sessions/:id`                       | Delete a session (refused if it has payments)  |
| GET    | `/sessions/:id/billing`               | What each player owes for the session          |
| POST   | `/sessions/:id/timeslots`             | Add a time slot                                |
| PUT    | `/timeslots/:id`                      | Change courts booked                           |
| DELETE | `/timeslots/:id`                      | Delete a slot                                  |
| POST   | `/timeslots/:id/players`              | Mark a player as playing in a slot             |
| DELETE | `/timeslots/:id/players/:playerId`    | Remove a player from a slot                    |
| GET    | `/players`                            | List players                                   |
| POST   | `/players`                            | Add a player                                   |
| GET    | `/players/:id/ledger`                 | A player's balance, sessions and payments      |
| GET    | `/ledger`                             | Every player's balance                         |
| GET    | `/payments`                           | List payments (`?player_id=`, `?session_id=`)  |
| POST   | `/payments`                           | Record a payment                               |
| DELETE | `/payments/:id`                       | Delete a payment                               |

## Roadmap

- [x] Core backend: sessions, slots, attendance, billing, payments, ledger
- [x] Unit tests for billing, ledger, sessions and auth
- [x] Environment configuration
- [x] JWT authentication
- [x] React frontend
- [x] OpenAPI / Swagger docs
- [ ] Docker and Docker Compose
- [ ] CI (tests, lint, build on every push)
- [ ] Deployment
