# Multi-Tenant SaaS Event Platform (Backend)

> A scalable, multi-tenant Event Management System written in **Go** with **PostgreSQL**.

## Key Features

* **Level 1: Core Foundation**
    * REST API with **Fiber (Go)**.
    * **JWT Authentication** with secure password hashing (**Bcrypt**).
    * Input Validation & Error Handling.

* **Level 2: Business Logic & Concurrency**
    * **ACID Transactions:** Uses **Pessimistic Locking** to prevent ticket overbooking.
    * **Fest Support:** Complex event hierarchies (Events & Sub-Events).
    * **Team Registration:** Validation logic for team constraints.

* **Level 3: SaaS Architecture**
    * **Multi-Tenancy:** Data isolation using `organization_id`. User A cannot access User B's data.
    * **DevOps:** Fully containerized with **Docker** & **Docker Compose**.
    * **Scalable DB:** **PostgreSQL** with optimized schema for tenant scoping.

---

## Tech Stack

* **Language:** Go (Golang) 1.25.5
* **Framework:** Fiber v3
* **Database:** PostgreSQL 15
* **ORM:** GORM (with constraints & associations)
* **Auth:** JWT (JSON Web Tokens)
* **DevOps:** Docker & Docker Compose

---

## Quick Start (Docker)

The easiest way to run the project.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/bit2swaz/ieee-cs-webdev.git
    cd backend
    ```

2.  **Create .env file:**
    ```bash
    cp .env.example .env
    # Ensure keys are set (DB_HOST=db for Docker)
    ```

3.  **Run with Docker Compose:**
    ```bash
    docker-compose up --build
    ```

    * API will be running at: `http://localhost:3000`
    * Database will be running on port `5433` (mapped to 5432).

---

## Manual Setup (Local)

1.  Ensure **PostgreSQL** is running locally.
2.  Update `.env` with `DB_HOST=localhost`.
3.  Run the application:
    ```bash
    go mod tidy
    go run cmd/api/main.go
    ```

---

## API Endpoints

### Authentication
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Register a new Tenant (Org) & Admin |
| `POST` | `/api/auth/login` | Login and receive JWT |

### Events (Protected)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/events` | Create a new Event (Scoped to Org) |
| `GET` | `/api/events` | List all Events for your Org |
| `POST` | `/api/events/:id/book` | Book a ticket (Thread-safe) |

---

## Architecture Highlights

### Concurrency Control (Ticket Booking)
To prevent race conditions where multiple users buy the last ticket simultaneously, I implemented **Database Transactions** with **Row-Level Locking**.

```go
// Code from internal/handlers/ticket.go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&event, "id = ?", eventID)

```

This ensures strict serializability for booking requests.

### Multi-Tenancy

Data isolation is enforced at the Application level using Middleware (`middleware/auth.go`). The `org_id` is extracted from the JWT and injected into every database query, ensuring tenants are strictly separated.

---

*Made with ❤️ by [bit2swaz](https://x.com/bit2swaz)*