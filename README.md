# 🚀 Starter Kit REST API GoFiber

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)
![Fiber Version](https://img.shields.io/badge/Fiber-v2-black?style=for-the-badge)
![GORM](https://img.shields.io/badge/GORM-v2-red?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

A production-ready, feature-rich RESTful API Starter Kit built with **Golang**, **Fiber v2**, and **GORM**.

This project is a port of a popular Node.js/Express boilerplate, adapted for Go best practices. It supports **SQLite** (for quick dev) and **PostgreSQL** (for production), includes fully automated Python API tests, and is Docker-ready.

---

## 📑 Table of Contents
- [Features](#-features)
- [Project Structure](#-project-structure)
- [Getting Started (Local)](#-getting-started-local-recommended)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [🧪 Automated API Testing](#-automated-api-testing-python)
- [🐳 Docker Deployment](#-docker-deployment)
- [License](#-license)

---

## ✨ Features

*   **⚡ Framework**: Built on [Fiber](https://gofiber.io), the fastest HTTP engine for Go.
*   **💾 Database Agnostic**: Uses [GORM](https://gorm.io) to support **SQLite** and **PostgreSQL** effortlessly.
*   **🔐 Authentication**: robust JWT authentication (Access & Refresh Tokens).
*   **🛡️ Security**: Helmet headers, CORS, Rate Limiting, and Password Hashing (Bcrypt).
*   **📩 Email Service**: SMTP integration for Forgot Password & Email Verification flows.
*   **📝 Validation**: Request body/query validation using struct tags.
*   **📄 Pagination**: Built-in helper for easy data pagination.
*   **📑 Swagger Docs**: Auto-generated API documentation.
*   **🧪 API Tests**: Automated Python scripts to replace manual Postman work.
*   **🐳 Docker**: Production-ready `Dockerfile` with volume persistence.

---

## 📂 Project Structure

A clean, standard Go project layout:

```
starter-kit-restapi-gofiber/
├── cmd/server/          # Application entry point
├── docs/                # Swagger generated files
├── internal/
│   ├── config/          # Environment configuration
│   ├── database/        # Database connection & migration
│   ├── handlers/        # HTTP Controllers
│   ├── middleware/      # Auth, Logging, Limiter
│   ├── models/          # GORM Database Models
│   ├── dto/             # Data Transfer Objects (Validation)
│   ├── routes/          # Route definitions
│   └── services/        # Business Logic
├── pkg/                 # Public utilities (JWT, Roles, Utils)
├── scripts/             # (Optional) API Test scripts
├── Dockerfile           # Docker build configuration
├── entrypoint.sh        # Docker entry script
└── README.md
```

---

## 🚀 Getting Started (Local) [Recommended]

We recommend running the project locally first to understand the flow before containerizing it.

### Prerequisites
*   [Go](https://golang.org/dl/) (version 1.20+)
*   [Python 3](https://www.python.org/) (for testing scripts)

### 1. Clone & Install
```bash
git clone https://github.com/yourusername/starter-kit-restapi-gofiber.git
cd starter-kit-restapi-gofiber

# Install Go dependencies
go mod tidy
```

### 2. Configure Environment
Copy the example file to `.env`:
```bash
cp .env.example .env
```

**For Quick Start (SQLite):**
Ensure your `.env` has these values. No installation required; it creates a file.
```ini
DB_DRIVER=sqlite
DB_NAME=gofiber_app.db
```

### 3. Run the Server
```bash
go run cmd/server/main.go
```
The server will start at `http://localhost:3000`.

---

## ⚙️ Configuration

Manage your configuration in `.env` (Local) or `.env.docker` (Docker).

| Key | Description | Example |
| :--- | :--- | :--- |
| `PORT` | Server Port | `3000` or `5005` |
| `DB_DRIVER` | Database Type | `sqlite` or `postgres` |
| `DB_HOST` | Database Host | `localhost` or `restapi-gofiber-postgres` |
| `JWT_SECRET` | Secret key for tokens | `s0m3s3cr3tk3y` |
| `SMTP_HOST` | Email Server Host | `smtp.mailtrap.io` |

---

## 📖 API Documentation

Once the server is running, you can view the full interactive API documentation via Swagger UI.

👉 **URL:** `http://localhost:3000/swagger/`

*(Note: If you modify code comments, regenerate docs using `swag init -g cmd/server/main.go -o ./docs`)*

---

## 🧪 Automated API Testing (Python)

Forget Postman! This project comes with a suite of Python scripts that automatically handle tokens, headers, and chaining requests.

### How to use
1.  Navigate to the project root (where the `*.py` scripts are).
2.  Ensure `utils.py` and `secrets.json` (auto-created) are present.
3.  Run the scripts in logical order:

```bash
# 1. Register a user (Saves tokens to secrets.json automatically)
python A1.auth_register.py

# 2. Login (Updates tokens in secrets.json)
python A2.auth_login.py
```

> **Note:** These scripts mimic real-world usage. You don't need to copy-paste tokens manually; the scripts manage them via `secrets.json`.

---

## 🐳 Docker Deployment

If you prefer using Docker, follow these manual steps to ensure persistence and networking.

### 1. Prepare Environment
Create a specific env file for Docker:
```bash
cp .env.example .env.docker
```
**Edit `.env.docker`** to use the PostgreSQL container we will create:
```ini
PORT=5005
DB_DRIVER=postgres
DB_HOST=restapi-gofiber-postgres
DB_USER=postgres
DB_PASSWORD=supersecretpassword
DB_NAME=gofiber_db
```

### 2. Create Network & Volumes
We use a custom network so containers can talk to each other, and volumes so data isn't lost when containers stop.

```bash
# Create Network
docker network create restapi_gofiber_network

# Create Volumes
docker volume create restapi_gofiber_db_volume
docker volume create restapi_gofiber_media_volume
```

### 3. Run Database Container (PostgreSQL)
Start the database first.
```bash
docker run -d \
  --network restapi_gofiber_network \
  --name restapi-gofiber-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=supersecretpassword \
  -e POSTGRES_DB=gofiber_db \
  -v restapi_gofiber_db_volume:/var/lib/postgresql/data \
  postgres:15-alpine
```

### 4. Build & Run Application Container
Build the Go image and run it.

```bash
# Build Image
docker build -t restapi-gofiber-app .

# Run Container
docker run -d -p 5005:5005 \
  --env-file .env.docker \
  --network restapi_gofiber_network \
  -v restapi_gofiber_media_volume:/app/media \
  --name restapi-gofiber-container \
  restapi-gofiber-app
```
Your API is now accessible at `http://localhost:5005`.

---

### 🛠️ Docker Management Commands

Useful commands to manage your containers:

#### View Logs
See what's happening inside the app.
```bash
docker logs -f restapi-gofiber-container
```

#### Stop Container
Safely stop the application.
```bash
docker stop restapi-gofiber-container
```

#### Start Container
Start the container again (data remains intact).
```bash
docker start restapi-gofiber-container
```

#### Remove Container
Remove the container instance (requires stopping first).
```bash
docker rm restapi-gofiber-container
```

#### List Volumes
See your persistent storage.
```bash
docker volume ls
```

#### ⚠️ Delete Volumes
**WARNING:** This deletes your Database and Media files **permanently**.
```bash
docker volume rm restapi_gofiber_db_volume
docker volume rm restapi_gofiber_media_volume
```

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.