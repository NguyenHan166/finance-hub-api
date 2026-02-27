# API Documentation - Finance Hub API

## Overview

RESTful API cho ứng dụng quản lý tài chính cá nhân. Base URL: `http://localhost:8080/api/v1`

## Authentication Overview

Tất cả endpoints (trừ Authentication endpoints và Health check) yêu cầu JWT token từ Supabase.

**Headers:**

```
Authorization: Bearer <supabase_jwt_token>
```

User ID được extract từ JWT token và tự động inject vào context.

---

## 0. Authentication API

### 0.1 Register

**POST** `/auth/register`

Đăng ký user mới. Authentication được handle bởi Supabase.

**Request Body:**

```json
{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "full_name": "Nguyễn Văn A"
}
```

**Validation:**

- `email`: required, valid email format
- `password`: required, ≥8 characters, phải có chữ hoa, chữ thường, số
- `full_name`: required, 1-100 characters

**Response 201:**

```json
{
    "status": "success",
    "message": "User registered successfully",
    "data": {
        "user": {
            "id": "uuid-string",
            "email": "user@example.com",
            "full_name": "Nguyễn Văn A",
            "created_at": "2026-02-27T10:00:00Z"
        },
        "session": {
            "access_token": "eyJhbGc...",
            "refresh_token": "refresh_token_string",
            "expires_in": 3600,
            "token_type": "bearer"
        }
    }
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Email already exists",
    "code": "EMAIL_EXISTS"
}
```

### 0.2 Login

**POST** `/auth/login`

Đăng nhập user. Authentication được handle bởi Supabase.

**Request Body:**

```json
{
    "email": "user@example.com",
    "password": "SecurePass123!"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "user": {
            "id": "uuid-string",
            "email": "user@example.com",
            "full_name": "Nguyễn Văn A",
            "avatar_url": "https://example.com/avatar.jpg"
        },
        "session": {
            "access_token": "eyJhbGc...",
            "refresh_token": "refresh_token_string",
            "expires_in": 3600,
            "token_type": "bearer"
        }
    }
}
```

**Response 401:**

```json
{
    "status": "error",
    "message": "Invalid email or password",
    "code": "INVALID_CREDENTIALS"
}
```

### 0.3 Logout

**POST** `/auth/logout`

Đăng xuất user (invalidate token).

**Headers:** Requires Authorization

**Response 200:**

```json
{
    "status": "success",
    "message": "Logged out successfully"
}
```

### 0.4 Refresh Token

**POST** `/auth/refresh`

Refresh access token khi hết hạn.

**Request Body:**

```json
{
    "refresh_token": "refresh_token_string"
}
```

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "access_token": "eyJhbGc...",
        "refresh_token": "new_refresh_token",
        "expires_in": 3600
    }
}
```

### 0.5 Get Current User

**GET** `/auth/me`

Lấy thông tin user hiện tại từ JWT token.

**Headers:** Requires Authorization

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "id": "uuid-string",
        "email": "user@example.com",
        "full_name": "Nguyễn Văn A",
        "avatar_url": "https://example.com/avatar.jpg",
        "phone": "+84901234567",
        "created_at": "2026-01-15T10:30:00Z",
        "updated_at": "2026-02-27T10:30:00Z"
    }
}
```

### 0.6 Forgot Password

**POST** `/auth/forgot-password`

Gửi email reset password.

**Request Body:**

```json
{
    "email": "user@example.com"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Password reset email sent"
}
```

### 0.7 Reset Password

**POST** `/auth/reset-password`

Reset password với token từ email.

**Request Body:**

```json
{
    "token": "reset_token_from_email",
    "new_password": "NewSecurePass123!"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Password reset successfully"
}
```

### 0.8 Login with Google (OAuth)

**GET** `/auth/google`

Initiate Google OAuth login flow. Redirects user to Google consent screen.

**Query Parameters:**

- `redirect_uri` (string, optional): URL to redirect after successful login, default = frontend URL

**Example:**

```
GET /api/v1/auth/google?redirect_uri=http://localhost:5173/auth/callback
```

**Response:**

Redirects to Google OAuth consent screen:

```
https://accounts.google.com/o/oauth2/v2/auth?
  client_id=YOUR_GOOGLE_CLIENT_ID
  &redirect_uri=http://localhost:8080/api/v1/auth/google/callback
  &response_type=code
  &scope=openid%20profile%20email
  &state=random_state_string
```

**User Flow:**

1. Frontend redirects user to `/api/v1/auth/google`
2. Backend redirects to Google consent screen
3. User logs in with Google and grants permissions
4. Google redirects back to `/api/v1/auth/google/callback`
5. Backend processes callback and redirects to frontend with token

### 0.9 Google OAuth Callback

**GET** `/auth/google/callback`

Handle OAuth callback từ Google. Endpoint này được Google gọi sau khi user authorize.

**Query Parameters:**

- `code` (string, required): Authorization code từ Google
- `state` (string, required): State parameter để verify request
- `error` (string, optional): Error code nếu user từ chối

**Success Flow:**

Sau khi verify code với Google, backend tạo/update user và redirect về frontend:

```
Redirect to: http://localhost:5173/auth/callback?token=eyJhbGc...&refresh_token=refresh_token_string
```

**Frontend sẽ nhận:**

```javascript
// Parse URL params
const urlParams = new URLSearchParams(window.location.search);
const accessToken = urlParams.get("token");
const refreshToken = urlParams.get("refresh_token");

// Save tokens and redirect to dashboard
localStorage.setItem("access_token", accessToken);
localStorage.setItem("refresh_token", refreshToken);
window.location.href = "/dashboard";
```

**Error Flow:**

Nếu có lỗi, redirect về frontend với error:

```
Redirect to: http://localhost:5173/login?error=access_denied&error_description=User%20denied%20access
```

**Response Data Structure (trong params):**

```
token=eyJhbGc...jwt_token
refresh_token=refresh_token_string
expires_in=3600
user_id=uuid-string
email=user@gmail.com
full_name=Nguyễn Văn A
avatar_url=https://lh3.googleusercontent.com/...
```

### 0.10 Login with Google (Direct Token)

**POST** `/auth/google/token`

Đăng nhập với Google ID token (alternative flow, dùng khi frontend tự handle Google Sign-In).

**Use Case:**

Khi frontend sử dụng Google Sign-In JavaScript library để lấy ID token trực tiếp, sau đó gửi token này lên backend để verify và tạo session.

**Request Body:**

```json
{
    "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjI3ZTc..."
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Login with Google successful",
    "data": {
        "user": {
            "id": "uuid-string",
            "email": "user@gmail.com",
            "full_name": "Nguyễn Văn A",
            "avatar_url": "https://lh3.googleusercontent.com/a/...",
            "auth_provider": "google",
            "created_at": "2026-02-27T10:00:00Z"
        },
        "session": {
            "access_token": "eyJhbGc...",
            "refresh_token": "refresh_token_string",
            "expires_in": 3600,
            "token_type": "bearer"
        },
        "is_new_user": false
    }
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Invalid Google ID token",
    "code": "INVALID_TOKEN"
}
```

**Frontend Implementation Example:**

```javascript
// Using Google Sign-In JavaScript library
function handleGoogleSignIn() {
    google.accounts.id.initialize({
        client_id: "YOUR_GOOGLE_CLIENT_ID",
        callback: handleCredentialResponse,
    });

    google.accounts.id.prompt();
}

async function handleCredentialResponse(response) {
    const idToken = response.credential;

    // Send to backend
    const res = await fetch("/api/v1/auth/google/token", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id_token: idToken }),
    });

    const data = await res.json();

    if (data.status === "success") {
        localStorage.setItem("access_token", data.data.session.access_token);
        localStorage.setItem("refresh_token", data.data.session.refresh_token);
        window.location.href = "/dashboard";
    }
}
```

### OAuth Configuration Notes

**Supabase OAuth Setup:**

1. **Enable Google Provider** trong Supabase Dashboard:
    - Go to Authentication → Providers
    - Enable Google
    - Add Google Client ID và Client Secret

2. **Google Cloud Console Setup:**
    - Create OAuth 2.0 Client ID
    - Authorized redirect URIs:
        - `https://your-project.supabase.co/auth/v1/callback`
        - `http://localhost:8080/api/v1/auth/google/callback` (development)
    - Authorized JavaScript origins:
        - `http://localhost:5173` (development)
        - `https://your-domain.com` (production)

3. **Environment Variables:**

```env
# .env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/google/callback
FRONTEND_URL=http://localhost:5173
```

**Security Considerations:**

- ✅ Always verify `state` parameter để prevent CSRF attacks
- ✅ Validate Google ID token signature với Google's public keys
- ✅ Check token expiration time
- ✅ Verify `aud` (audience) claim matches your Client ID
- ✅ Use HTTPS in production
- ✅ Store tokens securely (httpOnly cookies recommended for web)

**User Linking:**

Nếu user đã có account với email, Google OAuth sẽ:

- Link Google account với existing account
- Update avatar_url từ Google
- Set `auth_provider` = "google"
- User có thể login bằng cả email/password và Google

---

## 1. Accounts API

### 1.1 List All Accounts

**GET** `/accounts`

Lấy danh sách tất cả tài khoản của user.

**Query Parameters:**

- `page` (integer, optional): Page number, default = 1
- `limit` (integer, optional): Items per page, default = 10
- `sort_by` (string, optional): Field to sort by (name, balance, created_at), default = created_at
- `order` (string, optional): Sort order (asc, desc), default = desc

**Response 200:**

```json
{
    "status": "success",
    "message": "Accounts retrieved successfully",
    "data": {
        "data": [
            {
                "id": "uuid-string",
                "user_id": "uuid-string",
                "name": "Ví tiền mặt",
                "type": "cash",
                "currency": "VND",
                "balance": 5000000,
                "icon": "💵",
                "color": "#10B981",
                "bank_bin": null,
                "bank_code": null,
                "bank_logo": null,
                "account_number": null,
                "card_number": null,
                "created_at": "2026-01-15T10:30:00Z",
                "updated_at": "2026-01-15T10:30:00Z"
            }
        ],
        "page": 1,
        "limit": 10,
        "total_items": 5,
        "total_pages": 1
    }
}
```

### 1.2 Get Account Summary

**GET** `/accounts/summary`

Lấy tổng quan tất cả tài khoản (total balance, net worth change).

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "total_balance": 42500000,
        "accounts": [
            /* array of accounts */
        ],
        "net_worth_change": 2350000,
        "net_worth_change_percent": 5.8
    }
}
```

### 1.3 Get Account by ID

**GET** `/accounts/:id`

Lấy chi tiết một tài khoản.

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "id": "uuid",
        "user_id": "uuid",
        "name": "Techcombank",
        "type": "bank",
        "currency": "VND",
        "balance": 25000000,
        "bank_code": "TCB",
        "bank_logo": "https://api.vietqr.io/img/TCB.png",
        "account_number": "19036587456",
        "created_at": "2026-01-15T10:30:00Z",
        "updated_at": "2026-02-20T15:45:00Z"
    }
}
```

**Response 404:**

```json
{
    "status": "error",
    "message": "Account not found"
}
```

### 1.4 Create Account

**POST** `/accounts`

Tạo tài khoản mới.

**Request Body:**

```json
{
    "name": "Vietcombank",
    "type": "bank",
    "balance": 10000000,
    "icon": "🏦",
    "color": "#3B82F6",
    "bank_bin": "970436",
    "bank_code": "VCB",
    "bank_logo": "https://api.vietqr.io/img/VCB.png",
    "account_number": "1234567890"
}
```

**Validation Rules:**

- `name` (required, string, 1-100 chars)
- `type` (required, enum: "cash", "bank", "credit")
- `balance` (required, number >= 0)
- `icon` (optional, string)
- `color` (optional, string, hex color)
- `bank_code` (optional, string, for bank accounts)
- `account_number` (optional, string, for bank accounts)
- `card_number` (optional, string, for credit accounts)

**Response 201:**

```json
{
    "status": "success",
    "message": "Account created successfully",
    "data": {
        /* created account object */
    }
}
```

### 1.5 Update Account

**PUT** `/accounts/:id`

Cập nhật thông tin tài khoản (không cập nhật balance trực tiếp).

**Request Body:** (all fields optional)

```json
{
    "name": "VCB - Lương",
    "icon": "💰",
    "color": "#EF4444"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Account updated successfully",
    "data": {
        /* updated account */
    }
}
```

### 1.6 Delete Account

**DELETE** `/accounts/:id`

Xóa tài khoản. Chỉ xóa được nếu không có transaction nào liên quan.

**Response 200:**

```json
{
    "status": "success",
    "message": "Account deleted successfully"
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Cannot delete account with existing transactions"
}
```

---

## 2. Transactions API

### 2.1 List Transactions

**GET** `/transactions`

Lấy danh sách giao dịch với filters.

**Query Parameters:**

- `page` (integer): Page number
- `limit` (integer): Items per page
- `month` (string): Filter by month YYYY-MM
- `account_id` (string): Filter by account
- `category_id` (string): Filter by category
- `type` (string): Filter by type (income, expense, transfer)
- `search` (string): Search in merchant, note, tags
- `start_date` (string): Start date ISO 8601
- `end_date` (string): End date ISO 8601

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "data": [
            {
                "id": "uuid",
                "user_id": "uuid",
                "type": "expense",
                "amount": 150000,
                "date_time_iso": "2026-01-20T12:30:00Z",
                "account_id": "uuid",
                "to_account_id": null,
                "category_id": "uuid",
                "merchant": "Highlands Coffee",
                "note": "Họp team",
                "tags": ["work", "food"],
                "attachment_url": null,
                "created_at": "2026-01-20T12:35:00Z",
                "updated_at": "2026-01-20T12:35:00Z"
            }
        ],
        "page": 1,
        "limit": 20,
        "total_items": 156,
        "total_pages": 8
    }
}
```

### 2.2 Get Transaction by ID

**GET** `/transactions/:id`

Lấy chi tiết giao dịch.

**Response 200:**

```json
{
    "status": "success",
    "data": {
        /* transaction object */
    }
}
```

### 2.3 Create Transaction

**POST** `/transactions`

Tạo giao dịch mới. Tự động cập nhật balance của account(s).

**Request Body:**

```json
{
    "type": "expense",
    "amount": 500000,
    "date_time_iso": "2026-02-27T14:30:00Z",
    "account_id": "uuid",
    "category_id": "uuid",
    "merchant": "Shopee",
    "note": "Mua quần áo",
    "tags": ["shopping", "clothes"],
    "attachment_url": "https://storage.example.com/receipts/abc.jpg"
}
```

**For Transfer Transaction:**

```json
{
    "type": "transfer",
    "amount": 1000000,
    "date_time_iso": "2026-02-27T10:00:00Z",
    "account_id": "uuid-from",
    "to_account_id": "uuid-to",
    "note": "Chuyển tiền tiết kiệm"
}
```

**Validation Rules:**

- `type` (required, enum: "income", "expense", "transfer")
- `amount` (required, number > 0)
- `date_time_iso` (required, ISO 8601 datetime)
- `account_id` (required, valid account UUID)
- `to_account_id` (required if type=transfer, valid account UUID)
- `category_id` (optional, valid category UUID)
- `merchant` (optional, string, max 200 chars)
- `note` (optional, string, max 500 chars)
- `tags` (optional, array of strings)
- `attachment_url` (optional, string, valid URL)

**Business Logic:**

- For **expense**: decrease account balance
- For **income**: increase account balance
- For **transfer**: decrease from_account, increase to_account
- Validate sufficient balance for expense/transfer
- Category is optional for transfer transactions

**Response 201:**

```json
{
    "status": "success",
    "message": "Transaction created successfully",
    "data": {
        /* created transaction */
    }
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Insufficient balance"
}
```

### 2.4 Update Transaction

**PUT** `/transactions/:id`

Cập nhật giao dịch. Sẽ revert balance changes của transaction cũ và apply lại với data mới.

**Request Body:** (all fields optional except type)

```json
{
    "amount": 600000,
    "merchant": "Shopee (updated)",
    "note": "Mua quần áo + phụ kiện"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Transaction updated successfully",
    "data": {
        /* updated transaction */
    }
}
```

### 2.5 Delete Transaction

**DELETE** `/transactions/:id`

Xóa giao dịch và revert balance changes.

**Response 200:**

```json
{
    "status": "success",
    "message": "Transaction deleted successfully",
    "data": {
        /* deleted transaction object */
    }
}
```

### 2.6 Bulk Delete Transactions

**POST** `/transactions/bulk-delete`

Xóa nhiều transactions cùng lúc.

**Request Body:**

```json
{
    "ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "3 transactions deleted successfully",
    "data": {
        "deleted_count": 3
    }
}
```

---

## 3. Categories API

### 3.1 List Categories

**GET** `/categories`

Lấy tất cả categories của user.

**Query Parameters:**

- `type` (string, optional): Filter by type (income, expense, both)

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "id": "uuid",
            "user_id": "uuid",
            "name": "Ăn uống",
            "type": "expense",
            "parent_id": null,
            "icon": "🍜",
            "color": "#F59E0B",
            "is_default": true,
            "created_at": "2026-01-01T00:00:00Z",
            "updated_at": "2026-01-01T00:00:00Z"
        }
    ]
}
```

### 3.2 Get Category by ID

**GET** `/categories/:id`

**Response 200:**

```json
{
    "status": "success",
    "data": {
        /* category object */
    }
}
```

### 3.3 Create Category

**POST** `/categories`

Tạo category mới.

**Request Body:**

```json
{
    "name": "Đầu tư",
    "type": "expense",
    "parent_id": null,
    "icon": "📈",
    "color": "#8B5CF6"
}
```

**Validation:**

- `name` (required, string, 1-100 chars)
- `type` (required, enum: "income", "expense", "both")
- `parent_id` (optional, valid category UUID)
- `icon` (optional, string, emoji)
- `color` (optional, string, hex color)

**Response 201:**

```json
{
    "status": "success",
    "message": "Category created successfully",
    "data": {
        /* created category */
    }
}
```

### 3.4 Update Category

**PUT** `/categories/:id`

Cập nhật category. Không thể update default categories.

**Request Body:**

```json
{
    "name": "Đầu tư Crypto",
    "color": "#6366F1"
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Category updated successfully",
    "data": {
        /* updated category */
    }
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Cannot update default category"
}
```

### 3.5 Delete Category

**DELETE** `/categories/:id`

Xóa category. Không thể xóa default categories hoặc categories đang được sử dụng.

**Response 200:**

```json
{
    "status": "success",
    "message": "Category deleted successfully"
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Cannot delete category in use"
}
```

---

## 4. Budgets API

### 4.1 List Budgets

**GET** `/budgets`

Lấy danh sách budgets theo tháng.

**Query Parameters:**

- `month` (required, string): YYYY-MM format

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "id": "uuid",
            "user_id": "uuid",
            "month": "2026-02",
            "scope": "total",
            "category_id": null,
            "limit": 20000000,
            "spent": 15714000,
            "alert_enabled": true,
            "alert_threshold": 80,
            "created_at": "2026-02-01T00:00:00Z",
            "updated_at": "2026-02-27T10:00:00Z"
        },
        {
            "id": "uuid",
            "user_id": "uuid",
            "month": "2026-02",
            "scope": "category",
            "category_id": "uuid-food",
            "limit": 5000000,
            "spent": 4200000,
            "alert_enabled": true,
            "alert_threshold": 90,
            "created_at": "2026-02-01T00:00:00Z",
            "updated_at": "2026-02-27T10:00:00Z"
        }
    ]
}
```

### 4.2 Get Budget Detail

**GET** `/budgets/:id`

**Response 200:**

```json
{
    "status": "success",
    "data": {
        /* budget object */
    }
}
```

### 4.3 Create or Update Budget

**POST** `/budgets`

Tạo hoặc cập nhật budget. Nếu đã tồn tại budget cho month/scope/category thì update, không thì create.

**Request Body:**

```json
{
    "month": "2026-03",
    "scope": "category",
    "category_id": "uuid",
    "limit": 3000000,
    "alert_enabled": true,
    "alert_threshold": 85
}
```

**Validation:**

- `month` (required, string, YYYY-MM)
- `scope` (required, enum: "total", "category")
- `category_id` (required if scope=category)
- `limit` (required, number > 0)
- `alert_enabled` (required, boolean)
- `alert_threshold` (optional, number 0-100, required if alert_enabled=true)

**Response 201:**

```json
{
    "status": "success",
    "message": "Budget created successfully",
    "data": {
        /* budget object */
    }
}
```

### 4.4 Delete Budget

**DELETE** `/budgets/:id`

**Response 200:**

```json
{
    "status": "success",
    "message": "Budget deleted successfully"
}
```

---

## 5. Reports API

### 5.1 Get Overview Report

**GET** `/reports/overview`

Lấy báo cáo tổng quan trong khoảng thời gian.

**Query Parameters:**

- `start_date` (required, string): ISO 8601 date
- `end_date` (required, string): ISO 8601 date

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "total_income": 25000000,
        "total_expense": 15714000,
        "net_saving": 9286000,
        "saving_rate": 37.1,
        "transaction_count": 87,
        "avg_daily_expense": 566963,
        "compared_to_prev_month": {
            "income": 5.2,
            "expense": -8.5,
            "saving": 15.3
        }
    }
}
```

### 5.2 Get Category Report

**GET** `/reports/by-category`

Báo cáo chi tiêu theo category.

**Query Parameters:**

- `start_date` (required)
- `end_date` (required)

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "category_id": "uuid",
            "category_name": "Ăn uống",
            "amount": 4400000,
            "percentage": 28,
            "transaction_count": 32,
            "trend": "up"
        }
    ]
}
```

### 5.3 Get Merchant Report

**GET** `/reports/by-merchant`

Báo cáo chi tiêu theo merchant.

**Query Parameters:**

- `start_date` (required)
- `end_date` (required)
- `limit` (optional, default=20): Top N merchants

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "merchant": "Highlands Coffee",
            "amount": 980000,
            "transaction_count": 15,
            "percentage": 6.2
        }
    ]
}
```

### 5.4 Get Spending Trend

**GET** `/reports/spending-trend`

Xu hướng chi tiêu theo tuần/tháng.

**Query Parameters:**

- `start_date` (required)
- `end_date` (required)
- `interval` (optional, enum: "day", "week", "month", default="week")

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "period": "2026-02-24",
            "label": "Tuần 22-28/2",
            "amount": 3500000,
            "transaction_count": 12
        }
    ]
}
```

---

## 6. Alerts & Insights API

### 6.1 List Alerts

**GET** `/alerts`

Lấy danh sách alerts và insights.

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "id": "uuid",
            "user_id": "uuid",
            "severity": "warn",
            "title": "Ngân sách Mua sắm sắp vượt",
            "description": "Bạn đã chi 91.7% ngân sách Mua sắm (2.75M/3M)",
            "cta_label": "Xem chi tiết",
            "cta_route": "/budgets",
            "created_at": "2026-02-26T15:00:00Z",
            "is_read": false
        }
    ]
}
```

### 6.2 Dismiss Alert

**DELETE** `/alerts/:id`

Xóa/dismiss một alert.

**Response 200:**

```json
{
    "status": "success",
    "message": "Alert dismissed"
}
```

### 6.3 Get Forecast

**GET** `/forecasts/:month`

Lấy dự báo chi tiêu cho tháng.

**Path Parameters:**

- `month` (string): YYYY-MM

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "month": "2026-03",
        "predicted_total_expense": 16500000,
        "low": 14000000,
        "high": 19000000,
        "explanation_bullets": [
            "Dựa trên trung bình 3 tháng gần đây",
            "Có tăng nhẹ do tháng 3 thường có chi tiêu du lịch",
            "Lưu ý: Tết Thanh Minh có thể tăng chi tiêu gia đình"
        ],
        "generated_at": "2026-02-27T10:00:00Z"
    }
}
```

---

## 7. AI Chat API

### 7.1 Send Chat Message

**POST** `/ai/chat`

Gửi message cho AI assistant và nhận response.

**Request Body:**

```json
{
    "text": "Tôi chi tiêu như thế nào trong tháng này?",
    "context": {
        "month": "2026-02",
        "account_id": null
    }
}
```

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "message_id": "uuid",
        "reply_text": "Dựa trên dữ liệu chi tiêu của bạn trong tháng 2/2026...",
        "answer_cards": [
            {
                "title": "Phân tích chi tiêu tháng 2/2026",
                "metrics": [
                    { "label": "Tổng chi", "value": "15,714,000 ₫" },
                    { "label": "TB/ngày", "value": "566,963 ₫" }
                ],
                "explanation_bullets": [
                    "Chi tiêu Ăn uống chiếm 28%",
                    "Mua sắm tăng 15% so với tháng trước"
                ],
                "cta_label": "Xem báo cáo",
                "cta_route": "/reports"
            }
        ],
        "timestamp": "2026-02-27T14:30:00Z"
    }
}
```

### 7.2 Get Chat History

**GET** `/ai/chat/history`

Lấy lịch sử chat với AI.

**Query Parameters:**

- `limit` (optional, default=50): Number of messages

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "id": "uuid",
            "role": "user",
            "text": "Tôi chi tiêu như thế nào?",
            "timestamp": "2026-02-27T14:30:00Z"
        },
        {
            "id": "uuid",
            "role": "assistant",
            "text": "Dựa trên dữ liệu...",
            "timestamp": "2026-02-27T14:30:05Z",
            "answer_cards": [
                /* ... */
            ]
        }
    ]
}
```

---

## 8. Bank Integration API

### 8.1 Get Bank List

**GET** `/banks`

Lấy danh sách ngân hàng Việt Nam (từ VietQR API).

**Response 200:**

```json
{
    "status": "success",
    "data": [
        {
            "bin": "970436",
            "code": "VCB",
            "name": "Vietcombank",
            "logo": "https://api.vietqr.io/img/VCB.png"
        }
    ]
}
```

### 8.2 Parse Bank Transaction

**POST** `/banks/parse-transaction`

Parse SMS ngân hàng để tạo transaction tự động.

**Request Body:**

```json
{
    "sms_content": "VCB: -150,000 VND tai HIGHLANDS COFFEE luc 12:30 20/01/2026. SD: 25,000,000 VND"
}
```

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "type": "expense",
        "amount": 150000,
        "merchant": "HIGHLANDS COFFEE",
        "date_time_iso": "2026-01-20T12:30:00Z",
        "balance_after": 25000000,
        "bank_code": "VCB"
    }
}
```

---

## 9. Health Check

### 9.1 Health Status

**GET** `/health`

Check API health (no auth required).

**Response 200:**

```json
{
    "status": "ok",
    "timestamp": "2026-02-27T14:30:00Z",
    "database": "connected",
    "version": "1.0.0"
}
```

---

## 10. User & Settings API

### 10.1 Get User Profile

**GET** `/users/profile`

Lấy thông tin profile của user hiện tại.

**Headers:** Requires Authorization

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "id": "uuid-string",
        "email": "user@example.com",
        "full_name": "Nguyễn Văn A",
        "avatar_url": "https://example.com/avatar.jpg",
        "phone": "+84901234567",
        "preferences": {
            "currency": "VND",
            "language": "vi",
            "date_format": "DD/MM/YYYY",
            "theme": "light",
            "notifications": {
                "email_budget_alerts": true,
                "email_unusual_spending": true,
                "email_monthly_summary": true,
                "push_budget_exceeded": true,
                "push_low_balance": true
            }
        },
        "created_at": "2026-01-15T10:30:00Z",
        "updated_at": "2026-02-27T10:30:00Z"
    }
}
```

### 10.2 Update User Profile

**PUT** `/users/profile`

Cập nhật thông tin profile.

**Headers:** Requires Authorization

**Request Body:**

```json
{
    "full_name": "Nguyễn Văn B",
    "phone": "+84901234567",
    "avatar_url": "https://example.com/new-avatar.jpg"
}
```

**Validation:**

- `full_name`: optional, 1-100 characters
- `phone`: optional, valid phone format
- `avatar_url`: optional, valid URL

**Response 200:**

```json
{
    "status": "success",
    "message": "Profile updated successfully",
    "data": {
        "id": "uuid-string",
        "email": "user@example.com",
        "full_name": "Nguyễn Văn B",
        "avatar_url": "https://example.com/new-avatar.jpg",
        "phone": "+84901234567",
        "updated_at": "2026-02-27T11:00:00Z"
    }
}
```

### 10.3 Change Password

**PUT** `/users/password`

Đổi password của user.

**Headers:** Requires Authorization

**Request Body:**

```json
{
    "current_password": "OldPassword123!",
    "new_password": "NewPassword123!",
    "confirm_password": "NewPassword123!"
}
```

**Validation:**

- `current_password`: required
- `new_password`: required, ≥8 characters, phải có chữ hoa, chữ thường, số
- `confirm_password`: required, phải match với new_password

**Response 200:**

```json
{
    "status": "success",
    "message": "Password changed successfully"
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "Current password is incorrect",
    "code": "INVALID_PASSWORD"
}
```

### 10.4 Get User Preferences

**GET** `/users/preferences`

Lấy user preferences (currency, language, theme, notifications).

**Headers:** Requires Authorization

**Response 200:**

```json
{
    "status": "success",
    "data": {
        "currency": "VND",
        "language": "vi",
        "date_format": "DD/MM/YYYY",
        "number_format": "1,000,000.00",
        "theme": "light",
        "start_of_week": "monday",
        "fiscal_month_start": 1,
        "notifications": {
            "email_budget_alerts": true,
            "email_unusual_spending": true,
            "email_monthly_summary": true,
            "email_weekly_digest": false,
            "push_budget_exceeded": true,
            "push_low_balance": true,
            "push_transaction_reminders": false,
            "in_app_all_alerts": true,
            "in_app_insights": true
        }
    }
}
```

### 10.5 Update User Preferences

**PUT** `/users/preferences`

Cập nhật user preferences.

**Headers:** Requires Authorization

**Request Body:**

```json
{
    "currency": "USD",
    "language": "en",
    "theme": "dark",
    "notifications": {
        "email_budget_alerts": false,
        "push_budget_exceeded": true
    }
}
```

**Response 200:**

```json
{
    "status": "success",
    "message": "Preferences updated successfully",
    "data": {
        "currency": "USD",
        "language": "en",
        "date_format": "DD/MM/YYYY",
        "theme": "dark",
        "notifications": {
            "email_budget_alerts": false,
            "email_unusual_spending": true,
            "email_monthly_summary": true,
            "push_budget_exceeded": true,
            "push_low_balance": true
        }
    }
}
```

### 10.6 Upload Avatar

**POST** `/users/avatar`

Upload avatar image.

**Headers:**

- Requires Authorization
- `Content-Type: multipart/form-data`

**Request Body (Form Data):**

- `avatar`: Image file (JPG, PNG, max 5MB)

**Response 200:**

```json
{
    "status": "success",
    "message": "Avatar uploaded successfully",
    "data": {
        "avatar_url": "https://storage.example.com/avatars/uuid-string.jpg"
    }
}
```

**Response 400:**

```json
{
    "status": "error",
    "message": "File size exceeds 5MB limit",
    "code": "FILE_TOO_LARGE"
}
```

### 10.7 Export User Data

**POST** `/users/export-data`

Export toàn bộ data của user (GDPR compliance).

**Headers:** Requires Authorization

**Query Parameters:**

- `format` (string, optional): Export format (json, csv), default = json

**Response 200:**

```json
{
    "status": "success",
    "message": "Data export initiated",
    "data": {
        "download_url": "https://storage.example.com/exports/user-data-uuid.json",
        "expires_at": "2026-02-28T10:00:00Z"
    }
}
```

**Note:** Export có thể mất vài phút. Download URL có thời hạn 24 giờ.

### 10.8 Delete Account

**DELETE** `/users/account`

Xóa tài khoản user và toàn bộ dữ liệu (KHÔNG THỂ HOÀN TÁC).

**Headers:** Requires Authorization

**Request Body:**

```json
{
    "password": "UserPassword123!",
    "confirmation": "DELETE"
}
```

**Validation:**

- `password`: required, phải đúng password hiện tại
- `confirmation`: required, phải là string "DELETE"

**Response 200:**

```json
{
    "status": "success",
    "message": "Account deleted successfully"
}
```

**Response 403:**

```json
{
    "status": "error",
    "message": "Password is incorrect",
    "code": "INVALID_PASSWORD"
}
```

**What gets deleted:**

- User profile
- All accounts
- All transactions
- All budgets
- All categories (custom ones)
- All alerts
- All chat messages
- All forecasts

---

## Error Responses

### Standard Error Format

```json
{
    "status": "error",
    "message": "Human readable error message",
    "error": "Detailed error info (dev mode only)",
    "code": "ERROR_CODE"
}
```

### Common HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Missing/invalid JWT
- `403 Forbidden` - User không có quyền
- `404 Not Found` - Resource không tồn tại
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Server error

### Error Codes

- `INVALID_INPUT` - Validation error
- `UNAUTHORIZED` - Auth error
- `NOT_FOUND` - Resource not found
- `INSUFFICIENT_BALANCE` - Không đủ tiền
- `CATEGORY_IN_USE` - Category đang được sử dụng
- `ACCOUNT_HAS_TRANSACTIONS` - Account có transactions
- `DATABASE_ERROR` - Database error
- `EXTERNAL_API_ERROR` - External service error

---

## Rate Limiting

- **Default**: 100 requests per minute per user
- **Burst**: 20 requests per second
- Headers:
    - `X-RateLimit-Limit`: Total requests allowed
    - `X-RateLimit-Remaining`: Remaining requests
    - `X-RateLimit-Reset`: Reset timestamp

---

## Pagination

All list endpoints support pagination:

**Request:**

```
GET /transactions?page=2&limit=20
```

**Response:**

```json
{
    "data": [
        /* items */
    ],
    "page": 2,
    "limit": 20,
    "total_items": 156,
    "total_pages": 8
}
```

---

## Sorting

Supported via `sort_by` and `order` query params:

```
GET /transactions?sort_by=amount&order=desc
```

Common sort fields:

- `created_at` (default)
- `updated_at`
- `amount`
- `date_time_iso`
- `name`

---

## Date/Time Format

- **ISO 8601**: `2026-02-27T14:30:00Z`
- **Date only**: `2026-02-27`
- **Month**: `2026-02`
- Timezone: UTC

---

## Currency

- **Default**: VND (Vietnamese Dong)
- **Format**: Integer (không có decimals)
- **Example**: `1500000` = 1,500,000 VND

---

## Examples

### Create Expense Transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 450000,
    "date_time_iso": "2026-02-27T12:00:00Z",
    "account_id": "uuid-123",
    "category_id": "uuid-food",
    "merchant": "Phở 24",
    "note": "Ăn trưa team",
    "tags": ["food", "work"]
  }'
```

### Get Monthly Transactions

```bash
curl -X GET "http://localhost:8080/api/v1/transactions?month=2026-02&limit=50" \
  -H "Authorization: Bearer <token>"
```

### Create Budget

```bash
curl -X POST http://localhost:8080/api/v1/budgets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "month": "2026-03",
    "scope": "total",
    "limit": 20000000,
    "alert_enabled": true,
    "alert_threshold": 80
  }'
```

---

## WebSocket API (Future)

### Real-time Updates

**WS** `/ws/updates`

Subscribe để nhận real-time updates về:

- New transactions
- Budget alerts
- Balance changes
- AI insights

**Message Format:**

```json
{
    "type": "transaction.created",
    "data": {
        /* transaction object */
    },
    "timestamp": "2026-02-27T14:30:00Z"
}
```

---

**Last Updated**: February 27, 2026  
**API Version**: 1.0.0
