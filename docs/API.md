# Finance Hub API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Health Check

#### GET /health

Check if the API is running.

**Response:**

```json
{
    "success": true,
    "message": "Service is healthy",
    "data": {
        "status": "ok"
    }
}
```

---

### Accounts

#### POST /api/v1/accounts

Create a new account.

**Request Body:**

```json
{
    "name": "Main Checking",
    "type": "checking",
    "balance": 1000.0,
    "currency": "USD",
    "bank_name": "Chase Bank"
}
```

**Response:** (201 Created)

```json
{
    "success": true,
    "message": "Account created successfully",
    "data": {
        "id": "uuid",
        "user_id": "uuid",
        "name": "Main Checking",
        "type": "checking",
        "balance": 1000.0,
        "currency": "USD",
        "bank_name": "Chase Bank",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    }
}
```

#### GET /api/v1/accounts

Get all accounts.

**Query Parameters:**

- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 10, max: 100)

**Response:** (200 OK)

```json
{
  "success": true,
  "message": "Accounts retrieved successfully",
  "data": {
    "data": [
      {
        "id": "uuid",
        "name": "Main Checking",
        "type": "checking",
        "balance": 1000.00,
        ...
      }
    ],
    "page": 1,
    "limit": 10,
    "total_items": 5,
    "total_pages": 1
  }
}
```

#### GET /api/v1/accounts/:id

Get account by ID.

**Response:** (200 OK)

```json
{
  "success": true,
  "message": "Account retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "Main Checking",
    ...
  }
}
```

#### PUT /api/v1/accounts/:id

Update an account.

**Request Body:**

```json
{
    "name": "Updated Account Name",
    "balance": 1500.0,
    "is_active": true
}
```

**Response:** (200 OK)

#### DELETE /api/v1/accounts/:id

Delete an account.

**Response:** (200 OK)

```json
{
    "success": true,
    "message": "Account deleted successfully"
}
```

---

### Transactions

#### POST /api/v1/transactions

Create a new transaction.

**Request Body:**

```json
{
    "account_id": "uuid",
    "category_id": "uuid",
    "type": "expense",
    "amount": 50.0,
    "description": "Grocery shopping",
    "transaction_date": "2024-01-15T10:30:00Z",
    "notes": "Weekly groceries"
}
```

**Response:** (201 Created)

```json
{
    "success": true,
    "message": "Transaction created successfully",
    "data": {
        "id": "uuid",
        "account_id": "uuid",
        "category_id": "uuid",
        "type": "expense",
        "amount": 50.0,
        "description": "Grocery shopping",
        "transaction_date": "2024-01-15T10:30:00Z",
        "notes": "Weekly groceries",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

#### GET /api/v1/transactions

Get all transactions.

**Query Parameters:**

- `page` (int, optional): Page number
- `limit` (int, optional): Items per page

**Response:** (200 OK) - Similar paginated response as accounts

#### GET /api/v1/transactions/:id

Get transaction by ID.

#### DELETE /api/v1/transactions/:id

Delete a transaction.

---

### Categories

#### POST /api/v1/categories

Create a new category.

**Request Body:**

```json
{
    "name": "Groceries",
    "type": "expense",
    "icon": "🛒",
    "color": "#ef4444"
}
```

**Response:** (201 Created)

#### GET /api/v1/categories

Get all categories.

**Query Parameters:**

- `type` (string, optional): Filter by type (income/expense)

**Response:** (200 OK)

```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "Groceries",
      "type": "expense",
      "icon": "🛒",
      "color": "#ef4444",
      ...
    }
  ]
}
```

#### GET /api/v1/categories/:id

Get category by ID.

#### DELETE /api/v1/categories/:id

Delete a category (cannot delete default categories).

---

## Error Responses

All error responses follow this format:

```json
{
    "success": false,
    "message": "Error message",
    "error": "Detailed error information"
}
```

### Common Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Validation error or bad input
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Data Types

### Account Types

- `checking` - Checking account
- `savings` - Savings account
- `credit` - Credit card account
- `investment` - Investment account

### Transaction Types

- `income` - Income transaction
- `expense` - Expense transaction
- `transfer` - Transfer between accounts

### Category Types

- `income` - Income category
- `expense` - Expense category

## Notes

1. All amounts are in decimal format with 2 decimal places
2. All dates are in ISO 8601 format
3. All IDs are UUIDs
4. Pagination defaults: page=1, limit=10
5. Maximum items per page: 100
