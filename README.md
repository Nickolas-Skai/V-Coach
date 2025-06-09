# V-Coach

## Project Overview
V-Coach is a web-based coaching platform designed to facilitate interview preparation, user management, and interactive coaching sessions. The application is built using Go and PostgreSQL, and features a modular structure for scalability and maintainability. It provides tools for managing questions, user authentication, interview sessions, and more, with a focus on educational and coaching use cases.

## File and Folder Structure

- `cmd/web/` - Main web server application, including handlers, middleware, routes, and template rendering logic.
- `internal/data/` - Data models and logic for users, questions, responses, sessions, and validation.
- `ui/html/` - HTML templates for the web interface (dashboard, login, signup, interview, etc.).
- `ui/static/` - Static assets such as CSS and images.
- `migrations/` - Database migration scripts for PostgreSQL.
- `air/` - Contains the Air live-reload tool for Go development, including its configuration and documentation.
- `uploads/` - Directory for uploaded files and images.
- `Makefile` - Common development commands (run, test, database migrations, etc.).
- `go.mod` / `go.sum` - Go module dependencies.
- `.envrc` - Environment variables for local development (not committed, but referenced in code).

## Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd V-Coach
   ```

2. **Set up environment variables:**
   - Copy or create a `.envrc` file in the project root. Define variables such as `ADDRESS`, `VCOACH_DB_DSN`, and `SESSION_SECRET`.
   - Example `.envrc`:
     ```env
     ADDRESS=":8080"
     VCOACH_DB_DSN="postgres://vcoach:password@localhost:5432/vcoach?sslmode=disable"
     SESSION_SECRET="your-secret-key"
     ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Set up the database:**
   - Start PostgreSQL and create the `vcoach` database and user if not already present.
   - Run migrations:
     ```bash
     make db/migrations/up
     ```

5. **Run the application:**
   ```bash
   make run
   ```
   Or, for live reload during development (requires Air):
   ```bash
   cd air
   go install github.com/air-verse/air@latest
   air -c .air.toml
   ```

6. **Access the app:**
   Open your browser and go to `http://localhost:8080` (or the port you set in `.envrc`).

## Additional Notes

- Use the `Makefile` for common tasks like running tests (`make run/tests`), formatting code (`make fmt`), and managing migrations.
- The `air/` directory contains a local copy of the Air tool for live reloading Go apps. See `air/README.md` for usage.
- Uploaded files and images are stored in the `uploads/` directory.

---
For more details, see the documentation in each subdirectory and comments in the code.