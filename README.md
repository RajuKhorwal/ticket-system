# Ticket System (Golang)

Backend assignment for EVA Bharat. A user can register, log in, create tickets,
view only their own tickets, and move a ticket through its status
(open -> in_progress -> closed).

No database - data is kept in memory, as allowed in the assignment brief.
This kept the project small and avoided extra setup for a two-day timeline.

Dependencies used (besides the Go standard library):
- github.com/golang-jwt/jwt/v5 - creating and verifying login tokens
- golang.org/x/crypto/bcrypt - hashing passwords

## Project Structure

```
ticket-system/
├── main.go                # server setup and route wiring
├── models.go               # User and Ticket structs
├── store.go                 # in-memory storage (thread-safe)
├── auth.go                   # password hashing + JWT
├── middleware.go               # auth check for protected routes
├── handlers_auth.go            # register/login
├── handlers_tickets.go        # ticket endpoints
├── respond.go                # JSON response helpers
├── errors.go                # shared error values
├── Dockerfile               # build + run the service
├── .env.example             # sample env vars
└── go.mod / go.sum          # dependencies
```

## Running Locally

```bash
go mod download
go run .
```

Server starts on port 8080. Health check:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Running with Docker

```bash
docker build -t ticket-system .
docker run -p 8080:8080 ticket-system
curl http://localhost:8080/health
```

By default the app uses a built-in dev JWT secret, so the commands above
work with no extra setup. To use your own secret:

```bash
docker run -p 8080:8080 -e JWT_SECRET=your-long-random-secret ticket-system
```

## API

| Method | Endpoint | Auth | Purpose |
|---|---|---|---|
| GET | /health | No | Health check |
| POST | /auth/register | No | Register a new user |
| POST | /auth/login | No | Log in, get a JWT |
| POST | /tickets | Yes | Create a ticket |
| GET | /tickets | Yes | List your own tickets |
| GET | /tickets/{id} | Yes | Get one of your own tickets |
| PATCH | /tickets/{id}/status | Yes | Update your own ticket's status |

Protected routes need:
```
Authorization: Bearer <token>
```

### Example

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
# returns a token

curl -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Printer broken","description":"Office printer not working"}'

curl http://localhost:8080/tickets -H "Authorization: Bearer <TOKEN>"

curl http://localhost:8080/tickets/t1 -H "Authorization: Bearer <TOKEN>"

curl -X PATCH http://localhost:8080/tickets/t1/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```

## Assumptions

- In-memory storage, so data resets on server restart. Allowed by the brief,
  keeps the project simple and needs no DB setup on a free hosting tier.
- Status flow is enforced strictly step by step: open -> in_progress ->
  closed. A ticket can't skip from open straight to closed, and closed
  tickets can't move anywhere else.
- If a user requests a ticket owned by someone else, the API returns 404
  (not 403), so it doesn't confirm whether the ticket ID even exists.
- Minimum password length of 6 characters at registration - not specified
  in the brief, added as a basic sanity check.
- Server listens on 8080 by default, but reads PORT from the environment
  if set (some hosting platforms assign their own port).

## Deployment

Deployed on Render (free tier), built directly from the Dockerfile in this
repo.

- Deployed URL: <ADD_YOUR_URL_HERE>
- Health check: <ADD_YOUR_URL_HERE>/health

Render's free tier sleeps after inactivity and takes a few seconds to wake
on the next request - expected behavior, not a bug.

## Submission Checklist

- [ ] GitHub repository link
- [ ] Deployed application URL
- [ ] Public /health URL
- [ ] README (local run, Docker run, deployment URL, assumptions)
- [ ] .env.example
