#!/bin/bash

# Finance Hub API - Authentication Testing Script
# Usage: ./test_auth.sh

BASE_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Finance Hub - Authentication Tests"
echo "=========================================="
echo ""

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
curl -s -X GET http://localhost:8080/health | jq '.'
echo ""
echo ""

# Test 2: Register New User
echo -e "${YELLOW}Test 2: Register New User${NC}"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "Test1234",
    "confirm_password": "Test1234",
    "full_name": "Test User"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# Extract access token
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.refresh_token')

if [ "$ACCESS_TOKEN" != "null" ]; then
  echo -e "${GREEN}✓ Registration successful${NC}"
else
  echo -e "${RED}✗ Registration failed${NC}"
fi
echo ""
echo ""

# Test 3: Login
echo -e "${YELLOW}Test 3: Login with Email/Password${NC}"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "Test1234"
  }')

echo "$LOGIN_RESPONSE" | jq '.'

# Update access token from login
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')

if [ "$ACCESS_TOKEN" != "null" ]; then
  echo -e "${GREEN}✓ Login successful${NC}"
else
  echo -e "${RED}✗ Login failed${NC}"
fi
echo ""
echo ""

# Test 4: Get Profile (Protected Route)
echo -e "${YELLOW}Test 4: Get User Profile (Protected)${NC}"
PROFILE_RESPONSE=$(curl -s -X GET $BASE_URL/auth/profile \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$PROFILE_RESPONSE" | jq '.'

if echo "$PROFILE_RESPONSE" | jq -e '.success == true' > /dev/null; then
  echo -e "${GREEN}✓ Profile retrieved successfully${NC}"
else
  echo -e "${RED}✗ Failed to get profile${NC}"
fi
echo ""
echo ""

# Test 5: Access Protected Route without Token
echo -e "${YELLOW}Test 5: Access Protected Route Without Token (Should Fail)${NC}"
NO_AUTH_RESPONSE=$(curl -s -X GET $BASE_URL/auth/profile)
echo "$NO_AUTH_RESPONSE" | jq '.'

if echo "$NO_AUTH_RESPONSE" | jq -e '.success == false' > /dev/null; then
  echo -e "${GREEN}✓ Correctly rejected unauthorized request${NC}"
else
  echo -e "${RED}✗ Should have rejected request${NC}"
fi
echo ""
echo ""

# Test 6: Refresh Token
echo -e "${YELLOW}Test 6: Refresh Access Token${NC}"
REFRESH_RESPONSE=$(curl -s -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "$REFRESH_RESPONSE" | jq '.'

NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.data.access_token')

if [ "$NEW_ACCESS_TOKEN" != "null" ]; then
  echo -e "${GREEN}✓ Token refresh successful${NC}"
  ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
else
  echo -e "${RED}✗ Token refresh failed${NC}"
fi
echo ""
echo ""

# Test 7: Change Password
echo -e "${YELLOW}Test 7: Change Password${NC}"
CHANGE_PWD_RESPONSE=$(curl -s -X POST $BASE_URL/auth/change-password \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "Test1234",
    "new_password": "NewPass123"
  }')

echo "$CHANGE_PWD_RESPONSE" | jq '.'

if echo "$CHANGE_PWD_RESPONSE" | jq -e '.success == true' > /dev/null; then
  echo -e "${GREEN}✓ Password changed successfully${NC}"
else
  echo -e "${RED}✗ Password change failed${NC}"
fi
echo ""
echo ""

# Test 8: Login with New Password
echo -e "${YELLOW}Test 8: Login with New Password${NC}"
NEW_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "NewPass123"
  }')

echo "$NEW_LOGIN_RESPONSE" | jq '.'

if echo "$NEW_LOGIN_RESPONSE" | jq -e '.success == true' > /dev/null; then
  echo -e "${GREEN}✓ Login with new password successful${NC}"
else
  echo -e "${RED}✗ Login with new password failed${NC}"
fi
echo ""
echo ""

# Test 9: Logout
echo -e "${YELLOW}Test 9: Logout${NC}"
LOGOUT_RESPONSE=$(curl -s -X POST $BASE_URL/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$LOGOUT_RESPONSE" | jq '.'

if echo "$LOGOUT_RESPONSE" | jq -e '.success == true' > /dev/null; then
  echo -e "${GREEN}✓ Logout successful${NC}"
else
  echo -e "${RED}✗ Logout failed${NC}"
fi
echo ""
echo ""

# Test 10: Invalid Credentials
echo -e "${YELLOW}Test 10: Login with Invalid Credentials (Should Fail)${NC}"
INVALID_LOGIN=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "WrongPassword"
  }')

echo "$INVALID_LOGIN" | jq '.'

if echo "$INVALID_LOGIN" | jq -e '.success == false' > /dev/null; then
  echo -e "${GREEN}✓ Correctly rejected invalid credentials${NC}"
else
  echo -e "${RED}✗ Should have rejected invalid credentials${NC}"
fi
echo ""
echo ""

# Test 11: Weak Password (Should Fail)
echo -e "${YELLOW}Test 11: Register with Weak Password (Should Fail)${NC}"
WEAK_PWD=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "weak@example.com",
    "password": "weak",
    "confirm_password": "weak",
    "full_name": "Weak User"
  }')

echo "$WEAK_PWD" | jq '.'

if echo "$WEAK_PWD" | jq -e '.success == false' > /dev/null; then
  echo -e "${GREEN}✓ Correctly rejected weak password${NC}"
else
  echo -e "${RED}✗ Should have rejected weak password${NC}"
fi
echo ""
echo ""

# Test 12: Password Mismatch (Should Fail)
echo -e "${YELLOW}Test 12: Register with Password Mismatch (Should Fail)${NC}"
MISMATCH_PWD=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mismatch@example.com",
    "password": "Test1234",
    "confirm_password": "Different123",
    "full_name": "Mismatch User"
  }')

echo "$MISMATCH_PWD" | jq '.'

if echo "$MISMATCH_PWD" | jq -e '.success == false' > /dev/null; then
  echo -e "${GREEN}✓ Correctly rejected password mismatch${NC}"
else
  echo -e "${RED}✗ Should have rejected password mismatch${NC}"
fi
echo ""
echo ""

echo "=========================================="
echo -e "${GREEN}All Authentication Tests Complete!${NC}"
echo "=========================================="
echo ""
echo "Saved tokens for manual testing:"
echo "ACCESS_TOKEN=$ACCESS_TOKEN"
echo "REFRESH_TOKEN=$REFRESH_TOKEN"
