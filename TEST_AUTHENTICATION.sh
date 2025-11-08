#!/bin/bash
# Manual Testing Guide for JWT Authentication System

echo "=== JWT Authentication System Test Guide ==="
echo ""
echo "This guide will help you manually test the authentication system."
echo ""

echo "STEP 1: Set up the database"
echo "----------------------------"
echo "Run: lapis migrate"
echo "This will create the users and posts tables."
echo ""

echo "STEP 2: Start the server"
echo "------------------------"
echo "Run: lapis server"
echo "This will start the development server (usually on http://localhost:8080)"
echo ""

echo "STEP 3: Test the login page"
echo "---------------------------"
echo "Visit: http://localhost:8080/protected/login"
echo "You should see:"
echo "  - Login form with username and password fields"
echo "  - Sign up form with username, password, and confirm password fields"
echo ""

echo "STEP 4: Test signup"
echo "-------------------"
echo "1. Fill in the signup form:"
echo "   - Username: testuser"
echo "   - Password: password123"
echo "   - Confirm Password: password123"
echo "2. Click 'Sign Up'"
echo "3. You should be redirected to /protected/dashboard"
echo "4. You should see: 'Welcome, testuser! This is a protected page. Your XP: 0'"
echo ""

echo "STEP 5: Test logout"
echo "------------------"
echo "Visit: http://localhost:8080/protected/logout"
echo "You should be redirected back to /protected/login"
echo ""

echo "STEP 6: Test login"
echo "------------------"
echo "1. Fill in the login form:"
echo "   - Username: testuser"
echo "   - Password: password123"
echo "2. Click 'Login'"
echo "3. You should be redirected to /protected/dashboard"
echo ""

echo "STEP 7: Test authentication middleware"
echo "--------------------------------------"
echo "1. Clear your cookies or use incognito mode"
echo "2. Try to visit: http://localhost:8080/protected/dashboard"
echo "3. You should be redirected to /protected/login (not authenticated)"
echo ""

echo "STEP 8: Test invalid credentials"
echo "--------------------------------"
echo "1. Try to login with wrong password"
echo "2. You should see error message: 'Invalid username or password'"
echo ""

echo "STEP 9: Test duplicate signup"
echo "-----------------------------"
echo "1. Try to sign up with username 'testuser' again"
echo "2. You should see error message: 'User already exists'"
echo ""

echo "=== Testing Complete ==="
echo ""
echo "If all steps worked as expected, the authentication system is working correctly!"
echo ""
echo "To test programmatically, you can use curl:"
echo ""
echo "# Signup"
echo 'curl -X POST http://localhost:8080/protected/signup \\'
echo '  -H "Content-Type: application/x-www-form-urlencoded" \\'
echo '  -d "username=testuser2&password=pass123&confirm_password=pass123" \\'
echo '  -c cookies.txt -L -v'
echo ""
echo "# Login"
echo 'curl -X POST http://localhost:8080/protected/login \\'
echo '  -H "Content-Type: application/x-www-form-urlencoded" \\'
echo '  -d "username=testuser2&password=pass123" \\'
echo '  -c cookies.txt -L -v'
echo ""
echo "# Access protected route"
echo 'curl http://localhost:8080/protected/dashboard \\'
echo '  -b cookies.txt -L -v'
