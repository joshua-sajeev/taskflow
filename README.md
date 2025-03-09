# Taskflow - Job Scheduler

Taskflow is a **job scheduler** built using Golang. It allows users to enqueue jobs, process them concurrently using worker pools, and manage job execution efficiently. The project is designed to handle background task execution using a queue system.

## Features
- **Job Queue**: Enqueue and process jobs asynchronously.
- **Worker Pool**: Multiple workers process jobs concurrently.
- **Database Storage**: Jobs are stored in PostgreSQL using GORM.
- **REST API**: Built with Gin to handle job creation and monitoring.
- **Dockerized Setup**: PostgreSQL runs in a Docker container for easy deployment.
- **Graceful Shutdown**: Ensures workers stop properly when the application exits.
- **Concurrent Job Submission**: Allows multiple jobs to be submitted and processed simultaneously.

## Tech Stack
- **Golang**: Backend logic and worker execution.
- **Gin**: Web framework for handling API requests.
- **GORM**: ORM for PostgreSQL database interactions.
- **PostgreSQL**: Database for storing jobs.
- **Docker**: Containerized database setup.

## Installation & Setup
### Prerequisites
Ensure you have the following installed:
- [Golang](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Postman](https://www.postman.com/) (optional for testing API)

### 1. Clone the Repository
```bash
git clone https://github.com/joshua-sajeev/taskflow.git
cd taskflow
```

### 2. Start PostgreSQL with Docker
```bash
docker-compose up -d
```
This starts the PostgreSQL database inside a container.

### 3. Run the Application
```bash
go run main.go
```
The API server will start at `http://localhost:8080`

## API Endpoints
### 1. Create a Job
**Request:**
```http
POST /jobs/
```
**Body:**
```json
{
    "task": "Task 1"
}
```

### 2. Worker Pool Execution
- The worker pool picks up jobs from the queue and processes them concurrently.
- Logs indicate when a job is processed and updated.

## Running Multiple Jobs Concurrently
To submit multiple jobs simultaneously, use the following:
```bash
hey -n 10 -c 5 -m POST -H "Content-Type: application/json" -d '{"task":"Task 1"}' http://localhost:8080/jobs/
```
- `-n 10` sends 10 requests.
- `-c 5` processes 5 jobs concurrently.

Alternatively, using `curl` in a loop:
```bash
for i in {1..10}; do
    curl -X POST http://localhost:8080/jobs/ \
    -H "Content-Type: application/json" \
    -d '{"task":"Task 1"}' &
done
wait
```
This ensures jobs are submitted in parallel.

## Future Enhancements
- Implement job retry mechanism.
- Add job priority levels.
- Build a frontend dashboard for job monitoring.

---

**Taskflow** is designed for efficient background task execution, making it ideal for job scheduling needs. 🚀

