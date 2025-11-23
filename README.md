# TaskFlow - Learning Go, Docker & Clean Architecture

> A production-grade REST API built to master Go, Docker, and modern backend development practices.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://img.shields.io/badge/Tests-Passing-success)](./test)

TaskFlow is a robust task management API designed to demonstrate **Clean Architecture** principles in Go. It features secure authentication, containerized environments, and a comprehensive testing suite.

**This is a learning project.** It documents my journey from "Hello World" to a production-ready backend structure.

---

## Key Features

* **Clean Architecture:** Strict separation of concerns (Handlers → Services → Repositories).
* **Secure Auth:** JWT implementation with Bcrypt password hashing.
* **Containerization:** Optimized Multi-stage Docker builds for Dev and Prod.
* **Quality Assurance:** Unit & Integration tests with high coverage + Race detection.
* **Developer Experience:** Hot-reloading (Air), Swagger docs

---

## Documentation

I have documented the technical details and my learning process in the `docs/` folder:

* **[Architecture & Design](./docs/ARCHITECTURE.md)** - Breakdowns of the Clean Architecture layers and folder structure.
* **[API Reference](./docs/API.md)** - Endpoints, Request/Response examples, and Auth flow diagrams.
* **[Development Guide](./docs/DEVELOPMENT.md)** - Setup instructions, testing commands, and troubleshooting.
* **[What I Learned](./docs/LEARNING.md)** - A log of challenges faced and concepts mastered (Go routines, Interfaces, Docker networking).

---

## Quick Start

The easiest way to run the project is with Docker Compose.

### Prerequisites
* Docker & Docker Compose

### Run with One Command

```bash
# 1. Clone the repo
git clone https://github.com/joshua-sajeev/taskflow.git
cd taskflow

# 2. Setup Environment (Important!)
cp .env.example .env
# Open .env and set a secure JWT_SECRET

# 3. Start the App
docker-compose up --build
````

The API will be available at `http://localhost:8080`.

### Verify it works

Visit the Swagger UI to interact with the API:
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**
## Tech Stack

| Category      | Technology      | Usage                      |
| :------------ | :-------------- | :------------------------- |
| **Language**  | **Go (Golang)** | Core logic                 |
| **Framework** | **Gin**         | HTTP Routing & Middleware  |
| **Database**  | **MySQL 8.0**   | Persistent storage         |
| **ORM**       | **GORM**        | Data access & Migrations   |
| **DevOps**    | **Docker**      | Containerization & Compose |
| **Testing**   | **Testify**     | Assertions & Mocks         |
| **Docs**      | **Swagger**     | API Documentation          |

-----

## Contributing & Feedback

This project is open for code review\! If you see a non-idiomatic Go pattern or a security flaw, please open an issue.