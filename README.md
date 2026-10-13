# Gym Workout API 🏋️

**A RESTful backend API for managing gym workouts, exercises, and training logs**, built with **Go, Echo, PostgreSQL, and GORM**. It demonstrates JWT-based authentication, layered architecture, relational data modeling, and integration with a third-party BMI service.

> **Project status:** Portfolio / ongoing improvements. The source code and unit tests are available, but the previous Supabase database connection needs to be reconfigured. The API is **not currently verified as running end-to-end** and no live deployment is provided.

## Features

- **User accounts:** Register, log in, and retrieve an authenticated user profile.
- **Authentication:** Password hashing with bcrypt and protected endpoints using JWT.
- **Workout management:** Create, list, view, update, and delete workouts.
- **Exercise management:** Add exercises to workouts, list exercises by workout, and delete exercises.
- **Training logs:** Record sets, repetitions, and weight for each exercise, and retrieve log history.
- **BMI integration:** Retrieve BMI information through a third-party API (RapidAPI).
- **API documentation:** Swagger UI for exploring available endpoints.
- **Unit tests:** Mock-based tests for the workout use case.

## Tech Stack

| Layer / Purpose | Technology |
| --- | --- |
| Language | Go 1.25 |
| HTTP framework | Echo v4 |
| Database | PostgreSQL |
| ORM | GORM |
| Authentication | JWT, bcrypt |
| Input validation | go-playground/validator |
| Configuration | Viper, `.env` |
| API documentation | Swagger / swaggo |
| Testing | Go `testing`, Testify |
| Third-party integration | RapidAPI BMI Calculator |

## Project Structure

```text
.
├── config/           # PostgreSQL connection
├── delivery/         # HTTP handlers and route registration
├── docs/             # Generated Swagger documentation
├── domain/           # Entities, DTOs, and interfaces
├── middleware/       # JWT authentication middleware
├── repository/       # Database access using GORM
├── usecase/          # Business logic and unit tests
├── ddl.sql           # Database schema (see setup note)
├── ERD.png           # Entity relationship diagram
├── main.go           # Application entry point
├── .env.example      # Example environment configuration
└── go.mod
```

### Architecture

```text
HTTP Client
    |
    v
Echo Router + JWT Middleware
    |
    v
Delivery (HTTP Handler)
    |
    v
Usecase (Business Logic)
    |
    v
Repository (GORM)
    |
    v
PostgreSQL
```

## Entity Relationship Diagram

![Gym Workout API ERD](ERD.png)

## API Endpoints

All endpoints are prefixed with `/api`. A valid `Authorization: Bearer <token>` header is required for protected endpoints.

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| `POST` | `/api/users/register` | Register a new user | No |
| `POST` | `/api/users/login` | Log in and obtain JWT | No |
| `GET` | `/api/users` | Get authenticated user profile and BMI | Yes |
| `POST` | `/api/workouts` | Create a workout | Yes |
| `GET` | `/api/workouts` | List the user's workouts | Yes |
| `GET` | `/api/workouts/:id` | Get a workout by ID | Yes |
| `PUT` | `/api/workouts/:id` | Update a workout | Yes |
| `DELETE` | `/api/workouts/:id` | Delete a workout | Yes |
| `POST` | `/api/exercises` | Add an exercise to a workout | Yes |
| `GET` | `/api/exercises/workout/:workout_id` | List exercises in a workout | Yes |
| `DELETE` | `/api/exercises/:id` | Delete an exercise | Yes |
| `POST` | `/api/logs` | Create an exercise log | Yes |
| `GET` | `/api/logs/exercise/:exercise_id` | Get an exercise's log history | Yes |

### Example: Register a user

**Request:** `POST /api/users/register`

```json
{
  "username": "sample_user",
  "email": "sample@example.com",
  "full_name": "Sample User",
  "password": "replace_with_a_strong_password",
  "weight": 70,
  "height": 170
}
```

## Run Locally

### Prerequisites

- Go **1.25** or a compatible version
- Access to a PostgreSQL database (local or hosted)
- RapidAPI BMI Calculator key for the BMI-enabled profile endpoint

### 1. Clone the repository

```bash
git clone https://github.com/imam7710/go-gym-workout-api.git
cd go-gym-workout-api
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` with **your own** values:

```dotenv
DATABASE_URL=postgresql://postgres:YOUR_PASSWORD@localhost:5432/gym_app?sslmode=disable
JWT_SECRET=replace_with_a_long_random_secret
RAPIDAPI_KEY=your_rapidapi_key
PORT=8080
```

For a hosted Supabase database, replace `DATABASE_URL` with the connection string provided by your own active project (and enable TLS as required).

> Never commit `.env`, database passwords, JWT secrets, or API keys. `.env.example` should contain placeholders only.

### 3. Prepare the database

Create an empty PostgreSQL database. Review `ddl.sql` and execute the schema definitions using your SQL client or Supabase SQL Editor.

> **Known setup issue:** The original `ddl.sql` contains a duplicate `CREATE TABLE exercise_logs` statement. Remove the second duplicate definition before executing the script on a new database, unless this has already been corrected in your working copy.

### 4. Install dependencies and start the server

```bash
go mod download
go run main.go
```

With a valid database connection, the API runs at `http://localhost:8080` by default.

### 5. Open Swagger documentation

Visit **http://localhost:8080/swagger/index.html** after the server starts. The source-controlled OpenAPI files are in `docs/`.

## Tests and Code Quality

```bash
go test ./...
go vet ./...
```

Workout use-case unit tests are currently included in `usecase/workout_usecase_test.go`. Full end-to-end and security regression testing is still a work in progress.

## Current Limitations and Next Improvements

- Restore or replace the inactive Supabase connection, then verify the API end-to-end.
- Finish security regression checks for resource ownership and JWT validation.
- Remove the duplicate exercise log table definition from the schema if still present.
- Expand unit and integration tests across handlers, repositories, and exercise logs.
- Review Swagger documentation after subsequent endpoint changes.
- Set up deployment and CI when the application is ready.

## Project Background

Originally developed as part of the **HACKTIV8 Fulltime Golang Bootcamp** and refined as a personal backend development portfolio project. The repository showcases hands-on practice with RESTful API design, layered backend architecture, authentication, SQL database integration, and testing.

## Author

**Muhammad Imam Fadhilah**  
GitHub: [@imam7710](https://github.com/imam7710)  
Portfolio: [portofolio.fadhillah643.workers.dev](https://portofolio.fadhillah643.workers.dev/)
