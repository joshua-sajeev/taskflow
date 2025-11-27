#!/bin/bash
# demo.sh - Robust demonstration script for TaskFlow API

BASE_URL="http://localhost:8080/api"
EMAIL="demo@example.com"
PASSWORD="SecurePass123!"

echo "Starting TaskFlow Demo..."
echo ""

# Check API health
echo "Checking if API is reachable..."
curl -s "${BASE_URL}/health" >/dev/null
if [ $? -ne 0 ]; then
  echo "Error: TaskFlow API is not reachable at ${BASE_URL}"
  exit 1
fi
echo "API is reachable."
echo ""

# Register user
echo "1) Registering new user..."
REGISTER=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")

echo "$REGISTER"
echo ""

# Login
echo "2) Logging in..."
LOGIN=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")

TOKEN=$(echo "$LOGIN" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Error: Failed to retrieve token"
  echo "Login response: $LOGIN"
  exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."
echo ""

# Create tasks
echo "3) Creating tasks..."
CREATE_TASK=$(curl -s -X POST "${BASE_URL}/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"task":"Buy groceries"}')

echo "$CREATE_TASK"
echo ""

# List tasks
echo "4) Listing all tasks..."
TASKS=$(curl -s -X GET "${BASE_URL}/tasks" \
  -H "Authorization: Bearer $TOKEN")

echo "$TASKS"
echo ""

echo "Demo complete. Visit http://localhost:8080/swagger/index.html for interactive API documentation."
