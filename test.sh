#!/bin/bash

set -e

moonc .

rm -f abcms.test.sqlite

pkill nginx || true

echo "Running migrations..."
lapis migrate test

echo "Starting server..."
lapis server test &
SERVER_PID=$!

# Wait for server to start
sleep 2

trap 'kill $SERVER_PID' EXIT

# Test signup
echo "Testing signup..."
curl -s -d "option=signup&username=testuser&password=password" -X POST http://localhost:8080/login | grep -q "Redirecting to /dashboard"

# Test signup with existing username
echo "Testing signup with existing username..."
curl -s -d "option=signup&username=testuser&password=password" -X POST http://localhost:8080/login | grep -q "Username already exists"

# Test login
echo "Testing login..."
curl -s -d "option=login&username=testuser&password=password" -X POST http://localhost:8080/login | grep -q "Redirecting to /dashboard"


echo "All tests passed!"
