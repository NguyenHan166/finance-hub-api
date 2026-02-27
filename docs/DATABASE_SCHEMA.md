# Database Schema - Finance Hub MongoDB

## Overview

MongoDB database cho Finance Hub API. Database name: `fmp_app`

## Connection

```
mongodb+srv://nvhan166:han1662003@cluster0.evbdltl.mongodb.net/fmp_app?appName=Cluster0
```

---

## Collections

### 1. users

Lưu thông tin user cơ bản. Auth được handle bởi Supabase, collection này store metadata.

**Schema:**

```javascript
{
  _id: String,                    // UUID from Supabase
  email: String,                  // Unique
  full_name: String,
  avatar_url: String,
  currency: String,               // Default: "VND"
  timezone: String,               // Default: "Asia/Ho_Chi_Minh"
  preferences: {
    default_account_id: String,
    theme: String,                // "light" | "dark" | "system"
    language: String,             // "vi" | "en"
    date_format: String,          // "dd/MM/yyyy"
    notification_enabled: Boolean
  },
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ created_at: -1 });
```

---

### 2. accounts

Tài khoản tài chính (ví, ngân hàng, thẻ tín dụng).

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,                // Reference to users._id
  name: String,                   // "Ví tiền mặt", "Techcombank"
  type: String,                   // "cash" | "bank" | "credit"
  currency: String,               // "VND"
  balance: Number,                // Current balance in cents
  icon: String,                   // Emoji or icon name
  color: String,                  // Hex color code

  // Bank-specific fields
  bank_bin: String,               // Bank identification number (VietQR)
  bank_code: String,              // "VCB", "TCB", etc.
  bank_logo: String,              // URL to bank logo
  account_number: String,         // Masked or full account number

  // Credit card fields
  card_number: String,            // Masked card number
  credit_limit: Number,           // Credit limit
  statement_date: Number,         // Day of month (1-31)
  due_date: Number,               // Day of month (1-31)

  // Status
  is_active: Boolean,             // Default: true
  is_excluded_from_total: Boolean,// Exclude from net worth calculation

  // Metadata
  display_order: Number,          // For sorting in UI
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.accounts.createIndex({ user_id: 1, created_at: -1 });
db.accounts.createIndex({ user_id: 1, type: 1 });
db.accounts.createIndex({ user_id: 1, is_active: 1 });
db.accounts.createIndex({ user_id: 1, display_order: 1 });
```

**Validation Rules:**

- `balance` phải >= 0 cho cash và bank accounts
- `balance` có thể < 0 cho credit accounts (debt)
- `type` phải là một trong: "cash", "bank", "credit"
- Nếu `type` = "bank" thì nên có `bank_code`
- Nếu `type` = "credit" thì nên có `credit_limit`

**Sample Document:**

```json
{
  "_id": "550e8400-e29b-41d4-a716-446655440001",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Vietcombank - Lương",
  "type": "bank",
  "currency": "VND",
  "balance": 25000000,
  "icon": "🏦",
  "color": "#3B82F6",
  "bank_bin": "970436",
  "bank_code": "VCB",
  "bank_logo": "https://api.vietqr.io/img/VCB.png",
  "account_number": "1234567890",
  "is_active": true,
  "is_excluded_from_total": false,
  "display_order": 1,
  "created_at": ISODate("2026-01-15T10:30:00Z"),
  "updated_at": ISODate("2026-02-27T14:30:00Z")
}
```

---

### 3. categories

Danh mục giao dịch (thu nhập, chi tiêu).

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,                // Reference to users._id
  name: String,                   // "Ăn uống", "Lương"
  type: String,                   // "income" | "expense" | "both"
  parent_id: String,              // Reference to categories._id (for sub-categories)
  icon: String,                   // Emoji
  color: String,                  // Hex color
  is_default: Boolean,            // System default categories
  display_order: Number,
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.categories.createIndex({ user_id: 1, type: 1 });
db.categories.createIndex({ user_id: 1, parent_id: 1 });
db.categories.createIndex({ user_id: 1, is_default: 1 });
db.categories.createIndex({ user_id: 1, display_order: 1 });
```

**Validation:**

- `type` phải là "income", "expense", hoặc "both"
- Không thể xóa categories có `is_default: true`
- Không thể xóa categories đang được sử dụng trong transactions

**Default Categories:**
System sẽ tạo sẵn các categories sau cho user mới:

**Income:**

- Lương 💰
- Thưởng 🎁
- Đầu tư 📈
- Thu nhập khác 💵

**Expense:**

- Ăn uống 🍜
- Di chuyển 🚗
- Mua sắm 🛍️
- Nhà cửa 🏠
- Giải trí 🎮
- Sức khỏe ⚕️
- Giáo dục 📚
- Chi phí khác 💸

**Sample Document:**

```json
{
  "_id": "650e8400-e29b-41d4-a716-446655440002",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Ăn uống",
  "type": "expense",
  "parent_id": null,
  "icon": "🍜",
  "color": "#F59E0B",
  "is_default": true,
  "display_order": 1,
  "created_at": ISODate("2026-01-01T00:00:00Z"),
  "updated_at": ISODate("2026-01-01T00:00:00Z")
}
```

---

### 4. transactions

Giao dịch tài chính (thu, chi, chuyển khoản).

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,                // Reference to users._id
  type: String,                   // "income" | "expense" | "transfer"
  amount: Number,                 // Amount in cents (always positive)
  date_time_iso: Date,            // Transaction date/time

  // Account references
  account_id: String,             // Reference to accounts._id (from account)
  to_account_id: String,          // Reference to accounts._id (for transfers only)

  // Categorization
  category_id: String,            // Reference to categories._id
  merchant: String,               // Merchant/vendor name
  note: String,                   // User note/memo
  tags: [String],                 // Array of tags

  // Attachments
  attachment_url: String,         // URL to receipt/invoice image
  attachment_type: String,        // "image" | "pdf"

  // Metadata
  location: {
    latitude: Number,
    longitude: Number,
    address: String
  },
  is_recurring: Boolean,          // Part of recurring transaction
  recurring_id: String,           // Reference to recurring_transactions._id
  source: String,                 // "manual" | "sms" | "api" | "recurring"

  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.transactions.createIndex({ user_id: 1, date_time_iso: -1 });
db.transactions.createIndex({ user_id: 1, type: 1, date_time_iso: -1 });
db.transactions.createIndex({ user_id: 1, account_id: 1, date_time_iso: -1 });
db.transactions.createIndex({ user_id: 1, category_id: 1 });
db.transactions.createIndex({ user_id: 1, merchant: 1 });
db.transactions.createIndex({ user_id: 1, tags: 1 });
db.transactions.createIndex({ user_id: 1, created_at: -1 });

// Text search index for merchant and note
db.transactions.createIndex({
    merchant: "text",
    note: "text",
    tags: "text",
});

// Compound index for month queries
db.transactions.createIndex({
    user_id: 1,
    date_time_iso: 1,
});
```

**Validation:**

- `amount` phải > 0
- `type` phải là "income", "expense", hoặc "transfer"
- Nếu `type` = "transfer" thì `to_account_id` là required
- Nếu `type` = "transfer" thì `account_id` != `to_account_id`
- `date_time_iso` không được trong tương lai quá xa (max 1 ngày)

**Sample Documents:**

**Expense Transaction:**

```json
{
  "_id": "750e8400-e29b-41d4-a716-446655440003",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "expense",
  "amount": 450000,
  "date_time_iso": ISODate("2026-02-27T12:00:00Z"),
  "account_id": "550e8400-e29b-41d4-a716-446655440001",
  "to_account_id": null,
  "category_id": "650e8400-e29b-41d4-a716-446655440002",
  "merchant": "Phở 24",
  "note": "Ăn trưa team",
  "tags": ["food", "work"],
  "attachment_url": null,
  "source": "manual",
  "is_recurring": false,
  "created_at": ISODate("2026-02-27T12:05:00Z"),
  "updated_at": ISODate("2026-02-27T12:05:00Z")
}
```

**Transfer Transaction:**

```json
{
  "_id": "750e8400-e29b-41d4-a716-446655440004",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "transfer",
  "amount": 1000000,
  "date_time_iso": ISODate("2026-02-27T10:00:00Z"),
  "account_id": "550e8400-e29b-41d4-a716-446655440001",
  "to_account_id": "550e8400-e29b-41d4-a716-446655440005",
  "category_id": null,
  "merchant": null,
  "note": "Chuyển sang tiết kiệm",
  "tags": ["savings"],
  "source": "manual",
  "created_at": ISODate("2026-02-27T10:05:00Z"),
  "updated_at": ISODate("2026-02-27T10:05:00Z")
}
```

---

### 5. budgets

Ngân sách theo tháng (toàn bộ hoặc theo category).

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,                // Reference to users._id
  month: String,                  // "YYYY-MM" format
  scope: String,                  // "total" | "category"
  category_id: String,            // Reference to categories._id (if scope="category")
  limit: Number,                  // Budget limit
  spent: Number,                  // Current spending (calculated)
  alert_enabled: Boolean,
  alert_threshold: Number,        // Percentage (0-100)
  last_alert_at: Date,            // Last time alert was sent
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.budgets.createIndex(
    {
        user_id: 1,
        month: 1,
        scope: 1,
        category_id: 1,
    },
    { unique: true },
);
db.budgets.createIndex({ user_id: 1, month: 1 });
db.budgets.createIndex({ user_id: 1, alert_enabled: 1 });
```

**Validation:**

- `month` format phải là "YYYY-MM"
- `scope` phải là "total" hoặc "category"
- Nếu `scope` = "category" thì `category_id` là required
- `limit` phải > 0
- `alert_threshold` phải 0-100 nếu `alert_enabled` = true
- Mỗi user chỉ có 1 budget per month/scope/category (unique constraint)

**Business Logic:**

- `spent` được tính tự động từ transactions
- Alert được trigger khi % vượt `alert_threshold`
- Không gửi alert spam (check `last_alert_at`)

**Sample Document:**

```json
{
  "_id": "850e8400-e29b-41d4-a716-446655440005",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "month": "2026-02",
  "scope": "category",
  "category_id": "650e8400-e29b-41d4-a716-446655440002",
  "limit": 5000000,
  "spent": 4200000,
  "alert_enabled": true,
  "alert_threshold": 90,
  "last_alert_at": null,
  "created_at": ISODate("2026-02-01T00:00:00Z"),
  "updated_at": ISODate("2026-02-27T14:30:00Z")
}
```

---

### 6. recurring_transactions

Template cho giao dịch định kỳ.

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,
  name: String,                   // "Monthly Rent", "Salary"
  type: String,                   // "income" | "expense"
  amount: Number,
  frequency: String,              // "daily" | "weekly" | "monthly" | "yearly"
  interval: Number,               // Every N frequency units (e.g., every 2 weeks)
  start_date: Date,
  end_date: Date,                 // Optional
  next_date: Date,                // Next occurrence

  // Template data
  account_id: String,
  category_id: String,
  merchant: String,
  note: String,
  tags: [String],

  // Status
  is_active: Boolean,
  auto_create: Boolean,           // Auto-create transaction on next_date

  // Stats
  created_count: Number,          // Number of transactions created
  last_created_at: Date,

  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.recurring_transactions.createIndex({ user_id: 1, is_active: 1 });
db.recurring_transactions.createIndex({ next_date: 1, is_active: 1 });
db.recurring_transactions.createIndex({ user_id: 1, frequency: 1 });
```

**Sample Document:**

```json
{
  "_id": "950e8400-e29b-41d4-a716-446655440006",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Lương hàng tháng",
  "type": "income",
  "amount": 20000000,
  "frequency": "monthly",
  "interval": 1,
  "start_date": ISODate("2026-01-01T00:00:00Z"),
  "end_date": null,
  "next_date": ISODate("2026-03-01T00:00:00Z"),
  "account_id": "550e8400-e29b-41d4-a716-446655440001",
  "category_id": "650e8400-e29b-41d4-a716-446655440010",
  "merchant": "Company XYZ",
  "note": "Monthly salary",
  "tags": ["salary"],
  "is_active": true,
  "auto_create": true,
  "created_count": 2,
  "last_created_at": ISODate("2026-02-01T00:00:00Z"),
  "created_at": ISODate("2026-01-01T00:00:00Z"),
  "updated_at": ISODate("2026-02-01T00:05:00Z")
}
```

---

### 7. alerts

Alerts và insights cho user.

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,
  severity: String,               // "info" | "warn" | "danger"
  type: String,                   // "budget" | "forecast" | "insight" | "anomaly"
  title: String,
  description: String,

  // Actions
  cta_label: String,              // Call-to-action button text
  cta_route: String,              // Route to navigate

  // Metadata
  is_read: Boolean,
  dismissed_at: Date,
  related_entity_type: String,    // "budget" | "transaction" | "account"
  related_entity_id: String,

  created_at: Date,
  updated_at: Date
}
```

**Indexes:**

```javascript
db.alerts.createIndex({ user_id: 1, created_at: -1 });
db.alerts.createIndex({ user_id: 1, is_read: 1 });
db.alerts.createIndex({ user_id: 1, severity: 1 });
db.alerts.createIndex({ created_at: 1 }, { expireAfterSeconds: 2592000 }); // 30 days TTL
```

**Sample Document:**

```json
{
  "_id": "a50e8400-e29b-41d4-a716-446655440007",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "severity": "warn",
  "type": "budget",
  "title": "Ngân sách Mua sắm sắp vượt",
  "description": "Bạn đã chi 91.7% ngân sách Mua sắm (2.75M/3M)",
  "cta_label": "Xem chi tiết",
  "cta_route": "/budgets",
  "is_read": false,
  "dismissed_at": null,
  "related_entity_type": "budget",
  "related_entity_id": "850e8400-e29b-41d4-a716-446655440005",
  "created_at": ISODate("2026-02-26T15:00:00Z"),
  "updated_at": ISODate("2026-02-26T15:00:00Z")
}
```

---

### 8. forecasts

AI-generated spending forecasts.

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,
  month: String,                  // "YYYY-MM"
  predicted_total_expense: Number,
  low: Number,                    // Lower bound
  high: Number,                   // Upper bound
  confidence: Number,             // 0-100
  explanation_bullets: [String],  // Array of explanation points
  model_version: String,          // AI model version
  generated_at: Date,
  created_at: Date
}
```

**Indexes:**

```javascript
db.forecasts.createIndex({ user_id: 1, month: 1 }, { unique: true });
db.forecasts.createIndex({ user_id: 1, generated_at: -1 });
```

**Sample Document:**

```json
{
  "_id": "b50e8400-e29b-41d4-a716-446655440008",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "month": "2026-03",
  "predicted_total_expense": 16500000,
  "low": 14000000,
  "high": 19000000,
  "confidence": 78,
  "explanation_bullets": [
    "Dựa trên trung bình 3 tháng gần đây",
    "Có tăng nhẹ do tháng 3 thường có chi tiêu du lịch",
    "Lưu ý: Tết Thanh Minh có thể tăng chi tiêu gia đình"
  ],
  "model_version": "v1.0",
  "generated_at": ISODate("2026-02-27T10:00:00Z"),
  "created_at": ISODate("2026-02-27T10:00:00Z")
}
```

---

### 9. chat_messages

AI chat conversation history.

**Schema:**

```javascript
{
  _id: String,                    // UUID
  user_id: String,
  session_id: String,             // Group messages by session
  role: String,                   // "user" | "assistant"
  text: String,                   // Message text

  // AI response metadata
  answer_cards: [{
    title: String,
    metrics: [{
      label: String,
      value: String
    }],
    explanation_bullets: [String],
    cta_label: String,
    cta_route: String
  }],

  // Context
  context: {
    month: String,
    account_id: String,
    query_type: String            // "spending", "budget", "forecast", etc.
  },

  // Metadata
  tokens_used: Number,
  response_time_ms: Number,
  model: String,

  created_at: Date
}
```

**Indexes:**

```javascript
db.chat_messages.createIndex({ user_id: 1, created_at: -1 });
db.chat_messages.createIndex({ user_id: 1, session_id: 1, created_at: 1 });
db.chat_messages.createIndex(
    { created_at: 1 },
    { expireAfterSeconds: 7776000 },
); // 90 days TTL
```

---

## Relationships

```
users (1) ----< (N) accounts
users (1) ----< (N) categories
users (1) ----< (N) transactions
users (1) ----< (N) budgets
users (1) ----< (N) recurring_transactions
users (1) ----< (N) alerts
users (1) ----< (N) forecasts
users (1) ----< (N) chat_messages

accounts (1) ----< (N) transactions [account_id]
accounts (1) ----< (N) transactions [to_account_id]
categories (1) ----< (N) transactions
categories (1) ----< (N) budgets
categories (1) ----< (N) categories [parent_id] (self-referencing)
recurring_transactions (1) ----< (N) transactions
```

---

## MongoDB Queries Examples

### Get user's total balance

```javascript
db.accounts.aggregate([
    { $match: { user_id: "user-uuid", is_active: true } },
    {
        $group: {
            _id: "$user_id",
            total_balance: { $sum: "$balance" },
        },
    },
]);
```

### Get monthly transactions with pagination

```javascript
db.transactions
    .find({
        user_id: "user-uuid",
        date_time_iso: {
            $gte: ISODate("2026-02-01T00:00:00Z"),
            $lt: ISODate("2026-03-01T00:00:00Z"),
        },
    })
    .sort({ date_time_iso: -1 })
    .skip(0)
    .limit(20);
```

### Get category spending report

```javascript
db.transactions.aggregate([
    {
        $match: {
            user_id: "user-uuid",
            type: "expense",
            date_time_iso: {
                $gte: ISODate("2026-02-01T00:00:00Z"),
                $lt: ISODate("2026-03-01T00:00:00Z"),
            },
        },
    },
    {
        $group: {
            _id: "$category_id",
            total: { $sum: "$amount" },
            count: { $sum: 1 },
        },
    },
    { $sort: { total: -1 } },
]);
```

### Update account balance on transaction create

```javascript
// For expense
db.accounts.updateOne(
    { _id: "account-uuid" },
    {
        $inc: { balance: -amount },
        $set: { updated_at: new Date() },
    },
);

// For income
db.accounts.updateOne(
    { _id: "account-uuid" },
    {
        $inc: { balance: amount },
        $set: { updated_at: new Date() },
    },
);

// For transfer (2 operations)
db.accounts.updateOne(
    { _id: "from-account-uuid" },
    { $inc: { balance: -amount } },
);
db.accounts.updateOne(
    { _id: "to-account-uuid" },
    { $inc: { balance: amount } },
);
```

### Calculate budget spent

```javascript
db.transactions.aggregate([
    {
        $match: {
            user_id: "user-uuid",
            type: "expense",
            category_id: "category-uuid",
            date_time_iso: {
                $gte: ISODate("2026-02-01T00:00:00Z"),
                $lt: ISODate("2026-03-01T00:00:00Z"),
            },
        },
    },
    {
        $group: {
            _id: null,
            total_spent: { $sum: "$amount" },
        },
    },
]);
```

### Text search transactions

```javascript
db.transactions.find({
    $text: { $search: "highlands coffee" },
    user_id: "user-uuid",
});
```

---

## Data Migrations

### Initial Setup Script

```javascript
// Create indexes
db.users.createIndex({ email: 1 }, { unique: true });
db.accounts.createIndex({ user_id: 1, created_at: -1 });
db.transactions.createIndex({ user_id: 1, date_time_iso: -1 });
db.categories.createIndex({ user_id: 1, type: 1 });
db.budgets.createIndex(
    { user_id: 1, month: 1, scope: 1, category_id: 1 },
    { unique: true },
);

// Create default categories for new user
function createDefaultCategories(userId) {
    const categories = [
        { name: "Lương", type: "income", icon: "💰", color: "#10B981" },
        { name: "Thưởng", type: "income", icon: "🎁", color: "#14B8A6" },
        { name: "Ăn uống", type: "expense", icon: "🍜", color: "#F59E0B" },
        { name: "Di chuyển", type: "expense", icon: "🚗", color: "#3B82F6" },
        { name: "Mua sắm", type: "expense", icon: "🛍️", color: "#EC4899" },
        { name: "Nhà cửa", type: "expense", icon: "🏠", color: "#8B5CF6" },
    ];

    categories.forEach((cat, index) => {
        db.categories.insertOne({
            _id: generateUUID(),
            user_id: userId,
            ...cat,
            parent_id: null,
            is_default: true,
            display_order: index + 1,
            created_at: new Date(),
            updated_at: new Date(),
        });
    });
}
```

---

## Backup & Restore

### Backup

```bash
mongodump --uri="mongodb+srv://..." --db=fmp_app --out=/backup/
```

### Restore

```bash
mongorestore --uri="mongodb+srv://..." --db=fmp_app /backup/fmp_app/
```

---

## Performance Considerations

1. **Transactions Collection**: Largest collection, sử dụng compound indexes cho queries thường xuyên
2. **Budget Spent Calculation**: Cache spent value, chỉ recalculate khi có transaction mới
3. **Account Balance**: Denormalized - cập nhật trực tiếp khi có transaction
4. **TTL Indexes**: Auto-delete old alerts (30 days) và chat messages (90 days)
5. **Pagination**: Luôn sử dụng pagination cho transactions và chat messages

---

**Last Updated**: February 27, 2026
