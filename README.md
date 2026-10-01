# PlayLedger

PlayLedger tracks who played, who owes what, and who has paid for a recurring pickleball group.

A host books courts for a few hours, say 6pm to 10pm, but players come and go at different times. If Karan plays from 6 to 8, he should only pay for those two hours, and only his share of the courts that were booked in those hours. PlayLedger records attendance per time slot, works out each player's bill from that, and keeps a running balance so anything unpaid carries forward to the next session.

Built with **Go (Gin, GORM, PostgreSQL)** and **React (TypeScript, Vite)**.

## Features

- **Sessions split into time slots.** Create a session and have it cut into one-hour slots automatically, each with its own number of courts.
- **Attendance grid.** Players down the side, slots across the top; tap a cell to mark who played when.
- **Fair billing.** Each slot's court cost is split evenly between the players in that slot.
- **Payments and ledger.** Record payments, see every player's balance, and let unpaid amounts carry forward.
- **Player history.** Per-session bills and payments for each player, plus a WhatsApp reminder link for anyone who owes.
- **JWT authentication.** Tokens expire, and sign-ups can be switched off.
- **API docs.** Swagger UI at `/docs`.

## How billing works

For every time slot:

```
slot cost       = courts booked × court price per hour
cost per player = slot cost ÷ players in that slot
```

A player's bill for a session is the sum over the slots they played. Empty slots aren't billed.

Example, with courts at ₹300/hour:

| Slot        | Courts | Players              | Each pays |
| ----------- | ------ | -------------------- | --------- |
| 7pm – 8pm   | 2      | Rahul, Akhil, Anand  | ₹200      |
| 8pm – 9pm   | 1      | Rahul, Anand         | ₹150      |

So Akhil owes ₹200, while Rahul and Anand owe ₹350 each.

Payments are recorded per player and session, not per slot. A player's balance is:

```
balance = total paid − total billed      (negative = still owes)
```

If Rahul pays ₹300 of his ₹350, he carries −₹50 into his next session.

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
