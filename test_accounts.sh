#!/bin/bash

# Account Module Testing Script
# This script tests all account endpoints

BASE_URL="http://localhost:8080/api/v1"
TOKEN="" # Will be set after login

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print test result
print_result() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC} - $2"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ FAIL${NC} - $2"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Function to make API call
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -z "$data" ]; then
        curl -s -X $method \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            "$BASE_URL$endpoint"
    else
        curl -s -X $method \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint"
    fi
}

echo "=========================================="
echo "Account Module Test Suite"
echo "=========================================="
echo ""

# Step 1: Login to get token
echo -e "${YELLOW}Step 1: Login${NC}"
LOGIN_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "Test1234!"
    }' \
    "$BASE_URL/auth/login")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ "$TOKEN" != "null" ] && [ ! -z "$TOKEN" ]; then
    print_result 0 "Login successful"
else
    print_result 1 "Login failed"
    echo "Please create a test user first with email: test@example.com"
    exit 1
fi
echo ""

# Step 2: Test Get Bank List
echo -e "${YELLOW}Step 2: Get Bank List from VietQR${NC}"
BANKS_RESPONSE=$(api_call GET "/accounts/banks")
BANK_COUNT=$(echo $BANKS_RESPONSE | jq -r '.data | length')

if [ "$BANK_COUNT" -gt 0 ]; then
    print_result 0 "Get banks list (Count: $BANK_COUNT)"
else
    print_result 1 "Get banks list failed"
fi
echo ""

# Step 3: Search Banks
echo -e "${YELLOW}Step 3: Search Banks${NC}"
SEARCH_RESPONSE=$(api_call GET "/accounts/banks?q=vietcombank")
SEARCH_COUNT=$(echo $SEARCH_RESPONSE | jq -r '.data | length')

if [ "$SEARCH_COUNT" -gt 0 ]; then
    print_result 0 "Search banks (Query: vietcombank, Results: $SEARCH_COUNT)"
else
    print_result 1 "Search banks failed"
fi
echo ""

# Step 4: Create Cash Account
echo -e "${YELLOW}Step 4: Create Cash Account${NC}"
CASH_ACCOUNT=$(api_call POST "/accounts" '{
    "name": "Ví tiền mặt",
    "type": "cash",
    "balance": 5000000,
    "currency": "VND",
    "icon": "💵",
    "color": "#10B981"
}')

CASH_ACCOUNT_ID=$(echo $CASH_ACCOUNT | jq -r '.data.id')

if [ "$CASH_ACCOUNT_ID" != "null" ] && [ ! -z "$CASH_ACCOUNT_ID" ]; then
    print_result 0 "Create cash account (ID: $CASH_ACCOUNT_ID)"
else
    print_result 1 "Create cash account failed"
fi
echo ""

# Step 5: Create Bank Account with VietQR auto-fill
echo -e "${YELLOW}Step 5: Create Bank Account (VietQR auto-fill)${NC}"
BANK_ACCOUNT=$(api_call POST "/accounts" '{
    "name": "Vietcombank - Lương",
    "type": "bank",
    "balance": 10000000,
    "currency": "VND",
    "bank_code": "VCB",
    "account_number": "1234567890"
}')

BANK_ACCOUNT_ID=$(echo $BANK_ACCOUNT | jq -r '.data.id')
BANK_NAME=$(echo $BANK_ACCOUNT | jq -r '.data.bank_name')
BANK_LOGO=$(echo $BANK_ACCOUNT | jq -r '.data.bank_logo')

if [ "$BANK_ACCOUNT_ID" != "null" ] && [ ! -z "$BANK_ACCOUNT_ID" ]; then
    print_result 0 "Create bank account (ID: $BANK_ACCOUNT_ID)"
    
    if [ "$BANK_NAME" != "null" ] && [ ! -z "$BANK_NAME" ]; then
        print_result 0 "VietQR auto-fill bank name: $BANK_NAME"
    else
        print_result 1 "VietQR auto-fill bank name failed"
    fi
    
    if [ "$BANK_LOGO" != "null" ] && [ ! -z "$BANK_LOGO" ]; then
        print_result 0 "VietQR auto-fill bank logo"
    else
        print_result 1 "VietQR auto-fill bank logo failed"
    fi
else
    print_result 1 "Create bank account failed"
fi
echo ""

# Step 6: Create Credit Card Account
echo -e "${YELLOW}Step 6: Create Credit Card Account${NC}"
CREDIT_ACCOUNT=$(api_call POST "/accounts" '{
    "name": "Techcombank Visa",
    "type": "credit",
    "balance": 0,
    "currency": "VND",
    "bank_code": "TCB",
    "card_number": "**** **** **** 1234",
    "credit_limit": 50000000,
    "statement_date": 15,
    "due_date": 5
}')

CREDIT_ACCOUNT_ID=$(echo $CREDIT_ACCOUNT | jq -r '.data.id')

if [ "$CREDIT_ACCOUNT_ID" != "null" ] && [ ! -z "$CREDIT_ACCOUNT_ID" ]; then
    print_result 0 "Create credit card account (ID: $CREDIT_ACCOUNT_ID)"
else
    print_result 1 "Create credit card account failed"
fi
echo ""

# Step 7: List All Accounts
echo -e "${YELLOW}Step 7: List All Accounts${NC}"
ACCOUNTS_LIST=$(api_call GET "/accounts?page=1&limit=10")
ACCOUNTS_COUNT=$(echo $ACCOUNTS_LIST | jq -r '.data.total_items')

if [ "$ACCOUNTS_COUNT" -ge 3 ]; then
    print_result 0 "List accounts (Total: $ACCOUNTS_COUNT)"
else
    print_result 1 "List accounts failed"
fi
echo ""

# Step 8: Get Account Summary
echo -e "${YELLOW}Step 8: Get Account Summary${NC}"
SUMMARY=$(api_call GET "/accounts/summary")
TOTAL_ACCOUNTS=$(echo $SUMMARY | jq -r '.data.total_accounts')
TOTAL_BALANCE=$(echo $SUMMARY | jq -r '.data.total_balance')

if [ "$TOTAL_ACCOUNTS" -ge 3 ]; then
    print_result 0 "Get account summary (Accounts: $TOTAL_ACCOUNTS, Balance: $TOTAL_BALANCE)"
else
    print_result 1 "Get account summary failed"
fi
echo ""

# Step 9: Get Account by ID
echo -e "${YELLOW}Step 9: Get Account by ID${NC}"
ACCOUNT_DETAIL=$(api_call GET "/accounts/$CASH_ACCOUNT_ID")
ACCOUNT_NAME=$(echo $ACCOUNT_DETAIL | jq -r '.data.name')

if [ "$ACCOUNT_NAME" == "Ví tiền mặt" ]; then
    print_result 0 "Get account by ID (Name: $ACCOUNT_NAME)"
else
    print_result 1 "Get account by ID failed"
fi
echo ""

# Step 10: Update Account
echo -e "${YELLOW}Step 10: Update Account${NC}"
UPDATE_RESPONSE=$(api_call PUT "/accounts/$CASH_ACCOUNT_ID" '{
    "name": "Ví tiền mặt - Updated",
    "color": "#EF4444"
}')

UPDATED_NAME=$(echo $UPDATE_RESPONSE | jq -r '.data.name')

if [ "$UPDATED_NAME" == "Ví tiền mặt - Updated" ]; then
    print_result 0 "Update account (New name: $UPDATED_NAME)"
else
    print_result 1 "Update account failed"
fi
echo ""

# Step 11: Test validation - Create account without required fields
echo -e "${YELLOW}Step 11: Test Validation (Missing required fields)${NC}"
INVALID_RESPONSE=$(api_call POST "/accounts" '{
    "name": "",
    "type": "invalid_type"
}')

ERROR_MESSAGE=$(echo $INVALID_RESPONSE | jq -r '.error')

if [ "$ERROR_MESSAGE" != "null" ]; then
    print_result 0 "Validation works - Caught invalid input"
else
    print_result 1 "Validation failed - Should reject invalid input"
fi
echo ""

# Step 12: Test validation - Create credit card without credit_limit
echo -e "${YELLOW}Step 12: Test Validation (Credit card without limit)${NC}"
INVALID_CREDIT=$(api_call POST "/accounts" '{
    "name": "Test Credit",
    "type": "credit",
    "currency": "VND"
}')

ERROR_MESSAGE=$(echo $INVALID_CREDIT | jq -r '.error')

if [[ "$ERROR_MESSAGE" == *"credit limit"* ]]; then
    print_result 0 "Validation works - Caught missing credit limit"
else
    print_result 1 "Validation failed - Should require credit limit"
fi
echo ""

# Step 13: Delete Account
echo -e "${YELLOW}Step 13: Delete Account${NC}"
DELETE_RESPONSE=$(api_call DELETE "/accounts/$CASH_ACCOUNT_ID")
DELETE_STATUS=$(echo $DELETE_RESPONSE | jq -r '.status')

if [ "$DELETE_STATUS" == "success" ]; then
    print_result 0 "Delete account"
else
    print_result 1 "Delete account failed"
fi
echo ""

# Print Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "Total Tests:  $TOTAL_TESTS"
echo -e "${GREEN}Passed:       $PASSED_TESTS${NC}"
echo -e "${RED}Failed:       $FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    exit 1
fi
