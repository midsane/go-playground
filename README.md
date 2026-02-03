# Go Backend Fundamentals – Execution Manual

This repository is a **deliberate, modular learning system** for mastering Go backend development from first principles to production-level systems.

Each folder:
- Is **independent**
- Has its **own `go.mod`**
- Focuses on **one backend concept**
- Avoids magic and abstractions until fundamentals are solid

---

## 00-go-basics

### Tasks
- Implement:
  - value vs pointer receivers
  - interface satisfaction (implicit)
  - struct embedding
  - zero values
- Write examples showing:
  - when heap allocation happens
  - escape analysis (`go build -gcflags="-m"`)

### Advice
If you don’t deeply understand this folder, everything later becomes fragile.  
Most Go bugs are not framework bugs — they are **mental model bugs**.

---

## 01-project-structure

### Tasks
- Create:
  - `cmd/`
  - `internal/`
  - `pkg/`
- Build a small service using clean architecture:
  - handlers → services → repositories
- Enforce **dependency direction**

### Advice
Folder structure is **architecture made visible**.  
If imports go in random directions, your system will rot fast.

---

## 02-cli-app

### Tasks
- Build a CLI with:
  - `flag`
  - `cobra`
- Handle:
  - exit codes
  - invalid args
  - help commands

### Advice
CLIs teach **explicit control flow** and **error discipline** — skills most backend devs lack.

---

## 03-config-management

### Tasks
- Load config from:
  - env vars
  - `.env`
- Validate config on startup
- Fail fast if config is invalid

### Advice
A service that starts with bad config is already broken.  
Never let config errors leak into runtime.

---

## 04-logging

### Tasks
- Implement:
  - standard `log`
  - structured logging (zap)
- Add:
  - request ID
  - log levels
- Avoid logging inside libraries

### Advice
Logs are for **operators**, not developers.  
If logs don’t explain what happened and why, they’re useless.

---

## 05-error-handling

### Tasks
- Use:
  - sentinel errors
  - wrapped errors
- Demonstrate:
  - `errors.Is`
  - `errors.As`
- Create domain-specific errors

### Advice
Errors are **control flow**, not strings.  
If your errors are vague, debugging becomes gambling.

---

## 06-context

### Tasks
- Pass `context.Context` through layers
- Implement:
  - cancellation
  - timeouts
- Show what **NOT** to store in context

### Advice
Context misuse causes leaks, stuck goroutines, and outages.  
Treat context as a first-class API.

---

## 07-concurrency

### Tasks
- Implement:
  - worker pool
  - fan-in / fan-out
- Use:
  - mutex
  - channels
- Detect race conditions (`-race`)

### Advice
Concurrency doesn’t make code faster — it makes it harder.  
If you can’t explain why concurrency is needed, don’t use it.

---

## 08-http-server

### Tasks
- Build HTTP server using `net/http`
- Add:
  - graceful shutdown
  - read/write timeouts
- Handle SIGTERM correctly

### Advice
Your server lifecycle matters more than your handlers.  
Crashes during deploys are amateur mistakes.

---

## 09-routing-middleware

### Tasks
- Use `chi`
- Write custom middleware:
  - logging
  - auth
  - rate limiting
- Control middleware order

### Advice
Middleware order = **request behavior**.  
Get it wrong and you introduce invisible bugs.

---

## 10-auth

### Tasks
- Implement:
  - password hashing
  - JWT auth
  - refresh tokens
- Protect routes using middleware

### Advice
Auth bugs are security incidents, not “bugs”.  
Never roll crypto yourself.

---

## 11-database

### Tasks
- Connect to Postgres & MySQL
- Implement:
  - CRUD
  - transactions
  - migrations
- Compare:
  - raw SQL
  - GORM
  - sqlc

### Advice
If you don’t understand SQL deeply, ORMs will betray you.

---

## 12-cache

### Tasks
- Use Redis
- Implement:
  - cache-aside
  - TTL
- Handle cache misses and invalidation

### Advice
Caching hides performance problems — don’t use it blindly.

---

## 13-messaging

### Tasks
- Implement:
  - pub/sub
  - async consumers
- Simulate:
  - retries
  - failures
  - dead letters

### Advice
Distributed systems fail by default.  
Your job is to **design for failure**, not hope.

---

## 14-background-jobs

### Tasks
- Create:
  - worker queue
  - cron job
- Implement:
  - retries
  - idempotency

### Advice
Background jobs are where data corruption happens silently.  
Design defensively.

---

## 15-grpc

### Tasks
- Define proto files
- Implement:
  - unary RPC
  - streaming RPC
- Add interceptors

### Advice
gRPC is not “faster REST”.  
It’s a **different system design tradeoff**.

---

## 16-websockets

### Tasks
- Implement:
  - WebSocket server
  - connection lifecycle
- Handle:
  - disconnects
  - scaling concerns

### Advice
WebSockets introduce stateful complexity.  
If REST works, prefer REST.

---

## 17-search

### Tasks
- Index data in Elasticsearch
- Implement:
  - full-text search
  - pagination
- Compare DB search vs ES

### Advice
Search engines are not databases.  
Use them only when justified.

---

## 18-rate-limiting

### Tasks
- Implement:
  - token bucket
  - Redis-based limiter
- Apply per-user and per-IP limits

### Advice
Rate limiting is about **protecting systems**, not users.

---

## 19-testing

### Tasks
- Write:
  - unit tests
  - table-driven tests
- Mock dependencies
- Add integration tests

### Advice
If code is hard to test, it’s badly designed.

---

## 20-observability

### Tasks
- Add:
  - metrics
  - tracing
- Use OpenTelemetry concepts

### Advice
If you can’t observe it, you can’t operate it.

---

## 21-deployment

### Tasks
- Dockerize services
- Use multi-stage builds
- Understand env separation
- Basic Kubernetes manifests

---

## Final Rules

1. One folder at a time  
2. No skipping boring topics  
3. Write explanations in READMEs  
4. If you can’t explain it → redo it  
5. Depth > speed

---
