# Feature Specifications - Finance Hub

## Table of Contents

1. [Feature Overview](#feature-overview)
2. [Dashboard](#1-dashboard)
3. [Accounts Management](#2-accounts-management)
4. [Transactions](#3-transactions)
5. [Categories](#4-categories)
6. [Budgets](#5-budgets)
7. [Reports & Analytics](#6-reports--analytics)
8. [AI Chat Assistant](#7-ai-chat-assistant)
9. [Alerts & Insights](#8-alerts--insights)
10. [Settings](#9-settings)

---

## Feature Overview

Finance Hub là ứng dụng quản lý tài chính cá nhân toàn diện với các tính năng chính:

```
┌─────────────────────────────────────────────────────────┐
│                     Core Features                        │
├─────────────────────────────────────────────────────────┤
│ 📊 Dashboard         │ Tổng quan tài chính real-time    │
│ 💰 Accounts          │ Quản lý ví, ngân hàng, thẻ       │
│ 💸 Transactions      │ Ghi chép thu chi hàng ngày        │
│ 📁 Categories        │ Phân loại giao dịch               │
│ 🎯 Budgets           │ Đặt ngân sách và theo dõi         │
│ 📈 Reports           │ Báo cáo và phân tích chi tiết     │
│ 🤖 AI Assistant      │ Trợ lý tài chính thông minh       │
│ 🔔 Alerts            │ Cảnh báo và insights tự động      │
└─────────────────────────────────────────────────────────┘
```

---

## 1. Dashboard

### Overview

Dashboard là trang chính, cung cấp snapshot tổng quan về tình hình tài chính hiện tại của user.

### User Stories

- **US-D1**: Là user, tôi muốn xem tổng quan tài chính tháng hiện tại (thu, chi, tiết kiệm) để nắm được tình hình
- **US-D2**: Là user, tôi muốn xem xu hướng chi tiêu qua các tuần để phát hiện patterns
- **US-D3**: Là user, tôi muốn xem top categories chi tiêu nhiều nhất để biết nên cắt giảm ở đâu
- **US-D4**: Là user, tôi muốn xem dự báo chi tiêu tháng tới để lập kế hoạch
- **US-D5**: Là user, tôi muốn thấy các alerts quan trọng (budget vượt, unusual spending) để kịp thời xử lý

### UI Components

#### 1.1 Summary Cards

**Location**: Top row  
**Display**: 4 cards in a row (responsive)

**Cards:**

1. **Total Income** (Tổng thu)
    - Amount: VND format
    - Change vs last month: +X% (green) or -X% (red)
    - Icon: 💰 or trending up arrow

2. **Total Expense** (Tổng chi)
    - Amount: VND format
    - Change vs last month
    - Icon: 💸 or trending down arrow

3. **Net Saving** (Tiết kiệm)
    - Amount: Income - Expense
    - Percentage of income
    - Icon: 📊 or piggy bank

4. **Budget Status** (Tình trạng ngân sách)
    - Progress bar: Spent / Limit
    - Percentage used
    - Status color: green (<70%), orange (70-90%), red (>90%)

**Data Source:**

```
GET /api/v1/reports/overview?start_date=2026-02-01&end_date=2026-02-28
GET /api/v1/budgets?month=2026-02
```

#### 1.2 Spend Trend Chart

**Location**: Left column, top  
**Display**: Line chart với 4 tuần trong tháng

**Data Points:**

- X-axis: Tuần 1-4 (hoặc Tuần 23-27/2)
- Y-axis: Amount (VND, formatted)
- Tooltip: Total amount + số transactions

**Calculation Logic:**

```javascript
// Group transactions by week
const weeks = groupTransactionsByWeek(transactions, currentMonth);

weeks.forEach((week) => {
    const expenses = week.transactions.filter((t) => t.type === "expense");
    week.totalAmount = sum(expenses.map((t) => t.amount));
    week.count = expenses.length;
});
```

**Data Source:**

```
GET /api/v1/reports/spending-trend?start_date=...&end_date=...&interval=week
```

#### 1.3 Top Categories Section

**Location**: Left column, bottom-left  
**Display**: Horizontal bar chart (top 5 categories)

**Data:**

- Category name + icon
- Amount + percentage of total expense
- Bar width proportional to percentage

**Interaction:**

- Click category → navigate to `/reports?category_id=xxx`

**Data Source:**

```
GET /api/v1/reports/by-category?start_date=...&end_date=...
```

#### 1.4 Forecast Card

**Location**: Left column, bottom-right  
**Display**: Card với predicted amount + explanation

**Data:**

- Predicted total expense cho tháng tới
- Confidence range (low - high)
- 3-4 explanation bullets (why this prediction?)
- CTA button: "Thiết lập ngân sách" → `/budgets`

**Data Source:**

```
GET /api/v1/forecasts/2026-03
```

#### 1.5 Alerts List

**Location**: Right column  
**Display**: Vertical list of alert cards

**Alert Types:**

1. **Budget Alert** (warn/danger)
    - Title: "Ngân sách [Category] sắp vượt"
    - Description: "Bạn đã chi X% ngân sách"
    - CTA: "Xem chi tiết" → `/budgets`

2. **Unusual Spending** (info)
    - Title: "Chi tiêu bất thường phát hiện"
    - Description: "Chi [Category] tăng X% so với tháng trước"
    - CTA: "Xem chi tiết"

3. **Forecast Warning** (warn)
    - Title: "Dự báo vượt ngân sách"
    - Description: "Tháng tới dự kiến vượt XM ₫"
    - CTA: "Điều chỉnh"

**Interaction:**

- Dismiss button (X icon)
- Click alert → navigate to related page

**Data Source:**

```
GET /api/v1/alerts
```

### Business Logic

#### Monthly Summary Calculation

```typescript
function calculateMonthlySummary(transactions: Transaction[]) {
    const incomes = transactions.filter((t) => t.type === "income");
    const expenses = transactions.filter((t) => t.type === "expense");

    const totalIncome = sum(incomes.map((t) => t.amount));
    const totalExpense = sum(expenses.map((t) => t.amount));
    const netSaving = totalIncome - totalExpense;
    const savingRate = totalIncome > 0 ? (netSaving / totalIncome) * 100 : 0;

    return {
        totalIncome,
        totalExpense,
        netSaving,
        savingRate,
        transactionCount: transactions.length,
        avgDailyExpense: totalExpense / getDaysInMonth(),
    };
}
```

#### Comparison with Previous Month

```typescript
function calculateComparison(currentMonth: Summary, previousMonth: Summary) {
    return {
        income:
            ((currentMonth.totalIncome - previousMonth.totalIncome) /
                previousMonth.totalIncome) *
            100,
        expense:
            ((currentMonth.totalExpense - previousMonth.totalExpense) /
                previousMonth.totalExpense) *
            100,
        saving:
            ((currentMonth.netSaving - previousMonth.netSaving) /
                previousMonth.netSaving) *
            100,
    };
}
```

### Acceptance Criteria

- ✅ Dashboard loads within 2 seconds
- ✅ All cards show accurate real-time data
- ✅ Charts are interactive (tooltips on hover)
- ✅ Responsive layout (mobile, tablet, desktop)
- ✅ Auto-refresh data when month selector changes
- ✅ Alerts are dismissable and persistent
- ✅ Empty states shown when no data available

---

## 2. Accounts Management

### Overview

Quản lý tất cả tài khoản tài chính: ví tiền mặt, tài khoản ngân hàng, thẻ tín dụng.

### User Stories

- **US-A1**: Là user, tôi muốn tạo tài khoản mới (cash/bank/credit) để theo dõi số dư
- **US-A2**: Là user, tôi muốn xem danh sách tất cả tài khoản với số dư hiện tại
- **US-A3**: Là user, tôi muốn chỉnh sửa thông tin tài khoản (tên, icon, màu)
- **US-A4**: Là user, tôi muốn xóa tài khoản không dùng nữa (nếu không có transactions)
- **US-A5**: Là user, tôi muốn thấy total net worth (tổng tất cả accounts)
- **US-A6**: Là user, tôi muốn liên kết tài khoản ngân hàng với VietQR để auto-fill bank info

### Account Types

#### 1. Cash (Tiền mặt)

**Properties:**

- Name: "Ví tiền mặt"
- Type: cash
- Balance: Current amount
- Icon: 💵 (customizable)
- Color: #10B981 (customizable)

**Use Cases:**

- Tiền mặt trong ví
- Tiền tại nhà
- Petty cash

#### 2. Bank (Ngân hàng)

**Properties:**

- Name: "Vietcombank - Lương"
- Type: bank
- Balance: Current balance
- Bank Code: VCB
- Bank Logo: URL from VietQR
- Account Number: 1234567890 (optional, masked)
- Icon & Color: Auto from bank or custom

**Use Cases:**

- Tài khoản lương
- Tài khoản tiết kiệm
- Tài khoản thanh toán

#### 3. Credit (Thẻ tín dụng)

**Properties:**

- Name: "Techcombank Credit Card"
- Type: credit
- Balance: Current debt (negative number or positive for credit available)
- Credit Limit: 50,000,000
- Card Number: \***\* \*\*** \*\*\*\* 1234 (masked)
- Statement Date: 15 (ngày đóng sổ)
- Due Date: 5 (ngày đáo hạn)

**Use Cases:**

- Thẻ tín dụng các ngân hàng
- Credit line

### UI Components

#### 2.1 Account List Page

**Route**: `/accounts`

**Layout:**

- Header: "Tài khoản" + total net worth + "Thêm tài khoản" button
- Summary card: Total balance across all accounts, change vs last month
- Account cards grid (3 columns on desktop)

**Account Card:**

```
┌─────────────────────────────────────┐
│ 🏦 Vietcombank - Lương              │
│                                     │
│ 25,000,000 ₫                        │
│ Tài khoản ngân hàng                 │
│                                     │
│ [Xem chi tiết] [⋮]                  │
└─────────────────────────────────────┘
```

**Actions (⋮ menu):**

- Chỉnh sửa
- Xem transactions
- Xóa

#### 2.2 Create/Edit Account Modal

**Trigger**: "Thêm tài khoản" button or Edit action

**Form Fields:**

**Step 1: Choose Type**

- Radio buttons: Cash | Bank | Credit
- Visual cards with icons

**Step 2: Account Details**

For **Cash**:

- Name\* (text input)
- Initial Balance\* (number input, VND)
- Icon (emoji picker)
- Color (color picker)

For **Bank**:

- Name\* (auto-filled if using VietQR)
- Bank Selection (dropdown with logo, auto-fill from VietQR API)
    - Options: VCB, TCB, ACB, VPBank, etc.
- Account Number (optional, text)
- Initial Balance\* (number)
- Icon (auto from bank or custom)
- Color (auto from bank or custom)

For **Credit**:

- Name\* (e.g., "Techcombank Visa")
- Card Number (optional, masked)
- Credit Limit\* (number)
- Current Balance/Debt (number, default 0)
- Statement Date (day of month, 1-31)
- Due Date (day of month, 1-31)
- Icon & Color

**Validation:**

- Name: required, 1-100 chars
- Balance: required, number >= 0 (for cash/bank), any number (for credit)
- Bank code: required for bank type
- Credit limit: required for credit type

**Actions:**

- Hủy (cancel)
- Lưu (save)

#### 2.3 Account Detail Page

**Route**: `/accounts/:id`

**Sections:**

1. **Header**
    - Account name + icon
    - Current balance (large, prominent)
    - Last updated time
    - Edit button

2. **Quick Stats**
    - Total income this month
    - Total expense this month
    - Transaction count

3. **Recent Transactions**
    - List of last 10 transactions for this account
    - Link: "Xem tất cả" → `/transactions?account_id=xxx`

4. **Balance History Chart** (future)
    - Line chart showing balance over time

### Business Logic

#### Account Creation

```typescript
async function createAccount(input: AccountInput): Promise<Account> {
    // Validation
    if (!input.name || input.name.length < 1) {
        throw new Error("Name is required");
    }

    if (!["cash", "bank", "credit"].includes(input.type)) {
        throw new Error("Invalid account type");
    }

    if (input.balance === undefined || input.balance < 0) {
        throw new Error("Invalid balance");
    }

    // For bank accounts, validate bank code
    if (input.type === "bank" && !input.bankCode) {
        throw new Error("Bank code is required for bank accounts");
    }

    // Create account
    const account: Account = {
        id: generateUUID(),
        userId: getCurrentUserId(),
        name: input.name,
        type: input.type,
        currency: "VND",
        balance: input.balance,
        icon: input.icon || getDefaultIcon(input.type),
        color: input.color || getDefaultColor(input.type),
        bankCode: input.bankCode,
        bankLogo: input.bankCode ? getBankLogo(input.bankCode) : null,
        accountNumber: input.accountNumber,
        cardNumber: input.cardNumber,
        creditLimit: input.creditLimit,
        isActive: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    };

    // Save to database
    await AccountRepository.create(account);

    return account;
}
```

#### Balance Update Logic

Balance được update tự động khi có transaction:

```typescript
async function createTransaction(input: TransactionInput) {
    const transaction = await TransactionRepository.create(input);

    // Update account balance
    if (input.type === "income") {
        await AccountRepository.updateBalance(input.accountId, +input.amount);
    } else if (input.type === "expense") {
        await AccountRepository.updateBalance(input.accountId, -input.amount);
    } else if (input.type === "transfer") {
        await AccountRepository.updateBalance(input.accountId, -input.amount);
        await AccountRepository.updateBalance(
            input.toAccountId!,
            +input.amount,
        );
    }

    return transaction;
}
```

#### Account Deletion

```typescript
async function deleteAccount(accountId: string): Promise<void> {
    // Check if account has transactions
    const transactionCount =
        await TransactionRepository.countByAccount(accountId);

    if (transactionCount > 0) {
        throw new Error("Cannot delete account with existing transactions");
    }

    await AccountRepository.delete(accountId);
}
```

### VietQR Integration

**Get Bank List:**

```
GET https://api.vietqr.io/v2/banks
```

**Response:**

```json
{
    "code": "00",
    "desc": "Success",
    "data": [
        {
            "id": 17,
            "name": "Ngân hàng TMCP Công Thương Việt Nam",
            "code": "VCB",
            "bin": "970436",
            "shortName": "Vietcombank",
            "logo": "https://api.vietqr.io/img/VCB.png",
            "transferSupported": 1
        }
    ]
}
```

### Acceptance Criteria

- ✅ User có thể tạo 3 loại tài khoản (cash, bank, credit)
- ✅ VietQR integration hoạt động, auto-fill bank info
- ✅ Total net worth tính chính xác (tổng tất cả accounts)
- ✅ Balance update real-time khi có transaction
- ✅ Không thể xóa account có transactions
- ✅ Account cards hiển thị bank logo nếu có
- ✅ Form validation hoạt động đúng
- ✅ Modal có thể đóng bằng ESC key

---

## 3. Transactions

### Overview

Ghi chép và quản lý tất cả giao dịch tài chính: thu nhập, chi tiêu, chuyển khoản.

### User Stories

- **US-T1**: Là user, tôi muốn thêm giao dịch chi tiêu nhanh (số tiền, category, merchant)
- **US-T2**: Là user, tôi muốn thêm giao dịch thu nhập (số tiền, nguồn)
- **US-T3**: Là user, tôi muốn ghi chép chuyển khoản giữa các tài khoản
- **US-T4**: Là user, tôi muốn lọc transactions theo tháng, account, category, type
- **US-T5**: Là user, tôi muốn search transactions theo merchant name hoặc note
- **US-T6**: Là user, tôi muốn chỉnh sửa hoặc xóa transaction
- **US-T7**: Là user, tôi muốn xóa nhiều transactions cùng lúc (bulk delete)
- **US-T8**: Là user, tôi muốn thêm file đính kèm (receipt/invoice) cho transaction
- **US-T9**: Là user, tôi muốn tag transactions để dễ tìm kiếm

### Transaction Types

#### 1. Expense (Chi tiêu)

**Required Fields:**

- Amount (>0)
- Account (ví/ngân hàng nào chi tiền)
- Category (phân loại, e.g., Ăn uống)
- Date & Time

**Optional Fields:**

- Merchant (tên cửa hàng, e.g., "Highlands Coffee")
- Note (ghi chú thêm)
- Tags (e.g., ["work", "team-lunch"])
- Attachment (ảnh hóa đơn)

**Effect:**

- Account balance giảm
- Budget spent tăng (nếu có budget)

#### 2. Income (Thu nhập)

**Required Fields:**

- Amount (>0)
- Account (ví/ngân hàng nhận tiền)
- Category (phân loại, e.g., Lương)
- Date & Time

**Optional Fields:**

- Merchant/Source (nguồn thu nhập)
- Note
- Tags

**Effect:**

- Account balance tăng

#### 3. Transfer (Chuyển khoản)

**Required Fields:**

- Amount (>0)
- From Account
- To Account
- Date & Time

**Optional Fields:**

- Note (mục đích chuyển)
- Tags

**Effect:**

- From account balance giảm
- To account balance tăng
- Không ảnh hưởng budget (không phải thu/chi)

### UI Components

#### 3.1 Transaction List Page

**Route**: `/transactions`

**Layout:**

- Header: "Giao dịch" + "Thêm giao dịch" button
- Filter bar (sticky)
- Transaction list (grouped by date)
- Pagination

**Filter Bar:**

```
[Month Selector] [Account Selector] [Category Selector] [Type: All▾] [Search: 🔍]
```

**Filters:**

- Month: Dropdown, default = current month
- Account: Dropdown, "Tất cả tài khoản" + accounts list
- Category: Dropdown, "Tất cả danh mục" + categories list
- Type: Dropdown, "Tất cả" | "Chi tiêu" | "Thu nhập" | "Chuyển khoản"
- Search: Text input, search in merchant/note/tags

**Transaction List:**
Grouped by date, sorted descending:

```
Hôm nay - 27/02/2026
  ┌─────────────────────────────────────────────┐
  │ 🍜 Phở 24             -150,000 ₫   12:30    │
  │ Ăn uống              Ví tiền mặt            │
  └─────────────────────────────────────────────┘
  ┌─────────────────────────────────────────────┐
  │ ☕ Highlands Coffee   -85,000 ₫    09:15    │
  │ Ăn uống              Vietcombank            │
  └─────────────────────────────────────────────┘

Hôm qua - 26/02/2026
  ┌─────────────────────────────────────────────┐
  │ 🛍️ Shopee             -500,000 ₫   20:00    │
  │ Mua sắm              Techcombank Credit     │
  └─────────────────────────────────────────────┘
```

**Transaction Card:**

- Icon (category icon hoặc emoji)
- Merchant name (bold)
- Amount (color: red for expense, green for income)
- Time
- Category name (small text)
- Account name (small text)
- Click → open detail modal
- Long press / hover → show actions (Edit, Delete)

**Bulk Actions:**

- Checkbox on each transaction
- "Select all" checkbox in header
- Bulk actions bar appears when ≥1 selected:
    ```
    [2 selected] [Xóa] [Cancel]
    ```

#### 3.2 Add/Edit Transaction Modal

**Trigger**: "Thêm giao dịch" button hoặc click transaction card

**Form Layout:**

**Step 1: Type Selection** (for Add only)

- 3 tabs: Chi tiêu | Thu nhập | Chuyển khoản

**Step 2: Transaction Details**

**For Expense:**

```
┌─────────────────────────────────────┐
│ Chi tiêu                             │
├─────────────────────────────────────┤
│ Số tiền *                            │
│ [1,500,000] ₫                        │
│                                      │
│ Tài khoản *                          │
│ [Chọn tài khoản ▾]                   │
│                                      │
│ Danh mục *                           │
│ [Chọn danh mục ▾]                    │
│                                      │
│ Merchant (không bắt buộc)            │
│ [Nhập tên cửa hàng]                  │
│                                      │
│ Ngày & giờ *                         │
│ [27/02/2026] [14:30]                 │
│                                      │
│ Ghi chú (không bắt buộc)             │
│ [Nhập ghi chú...]                    │
│                                      │
│ Tags                                 │
│ [work] [x] [Add tag +]               │
│                                      │
│ Đính kèm hóa đơn                     │
│ [📎 Upload file]                     │
│                                      │
│         [Hủy]      [Lưu]             │
└─────────────────────────────────────┘
```

**For Income:**
Similar, nhưng:

- Title: "Thu nhập"
- Merchant → "Nguồn thu nhập"
- Category options: only income categories

**For Transfer:**

```
Số tiền *
[1,000,000] ₫

Từ tài khoản *
[Vietcombank ▾]

Đến tài khoản *
[Tiết kiệm ▾]

Ngày & giờ *
[27/02/2026] [10:00]

Ghi chú
[Chuyển tiết kiệm tháng]

[Hủy]      [Lưu]
```

**Validation:**

- Amount: required, >0, number
- Account: required
- To Account: required for transfer, cannot be same as from account
- Category: required for expense/income
- Date: required, không quá 1 năm trong quá khứ hoặc tương lai

**Auto-complete:**

- Merchant: suggest từ recent merchants
- Amount: suggest common amounts (50k, 100k, 500k)
- Date/Time: default = now

#### 3.3 Transaction Detail Modal

**Trigger**: Click transaction in list

**Display:**

```
┌───────────────────────────────────────┐
│ Chi tiết giao dịch         [✏️] [🗑️]    │
├───────────────────────────────────────┤
│ 🍜 Phở 24                              │
│                                        │
│ Số tiền:         -150,000 ₫            │
│ Loại:            Chi tiêu              │
│ Danh mục:        Ăn uống               │
│ Tài khoản:       Ví tiền mặt           │
│ Ngày giờ:        27/02/2026 12:30      │
│ Ghi chú:         Ăn trưa team          │
│ Tags:            #work #lunch          │
│ Đính kèm:        [📄 receipt.jpg]      │
│                                        │
│              [Đóng]                    │
└───────────────────────────────────────┘
```

Actions:

- Edit icon (✏️) → open edit modal
- Delete icon (🗑️) → confirm delete dialog

#### 3.4 Attachment Upload

**Supported Formats:**

- Images: JPG, PNG, WebP (max 5MB)
- PDF: Receipt, invoice (max 10MB)

**Upload Flow:**

1. Click "Upload file" button
2. File picker opens
3. Select file
4. Upload progress shown
5. Preview thumbnail displayed
6. Can remove before saving

**Storage:**

- Upload to cloud storage (Supabase Storage hoặc AWS S3)
- Store URL in transaction.attachment_url

### Business Logic

#### Create Transaction

```typescript
async function createTransaction(
    input: TransactionInput,
): Promise<Transaction> {
    // Validation
    if (input.amount <= 0) {
        throw new Error("Amount must be greater than 0");
    }

    if (!input.accountId) {
        throw new Error("Account is required");
    }

    if (input.type === "transfer" && !input.toAccountId) {
        throw new Error("To account is required for transfers");
    }

    if (input.type === "transfer" && input.accountId === input.toAccountId) {
        throw new Error("Cannot transfer to the same account");
    }

    if (["income", "expense"].includes(input.type) && !input.categoryId) {
        throw new Error("Category is required");
    }

    // Verify account belongs to user
    const account = await AccountRepository.getById(input.accountId, userId);
    if (!account) {
        throw new Error("Account not found");
    }

    // Check sufficient balance for expense/transfer
    if (["expense", "transfer"].includes(input.type)) {
        if (account.balance < input.amount) {
            throw new Error("Insufficient balance");
        }
    }

    // Create transaction
    const transaction = await TransactionRepository.create({
        ...input,
        userId,
        id: generateUUID(),
        createdAt: new Date(),
        updatedAt: new Date(),
    });

    // Update account balance
    await updateAccountBalance(input, transaction);

    // Recalculate budget if expense
    if (input.type === "expense") {
        await recalculateBudget(userId, input.categoryId, input.dateTimeISO);
    }

    // Check and create alerts
    await checkBudgetAlerts(userId);

    return transaction;
}
```

#### Update Account Balance

```typescript
async function updateAccountBalance(
    input: TransactionInput,
    transaction: Transaction,
) {
    if (input.type === "income") {
        await AccountRepository.updateBalance(input.accountId, +input.amount);
    } else if (input.type === "expense") {
        await AccountRepository.updateBalance(input.accountId, -input.amount);
    } else if (input.type === "transfer") {
        await AccountRepository.updateBalance(input.accountId, -input.amount);
        await AccountRepository.updateBalance(
            input.toAccountId!,
            +input.amount,
        );
    }
}
```

#### Delete Transaction

```typescript
async function deleteTransaction(transactionId: string, userId: string) {
    const transaction = await TransactionRepository.getById(
        transactionId,
        userId,
    );

    if (!transaction) {
        throw new Error("Transaction not found");
    }

    // Revert account balance changes
    if (transaction.type === "income") {
        await AccountRepository.updateBalance(
            transaction.accountId,
            -transaction.amount,
        );
    } else if (transaction.type === "expense") {
        await AccountRepository.updateBalance(
            transaction.accountId,
            +transaction.amount,
        );
    } else if (transaction.type === "transfer") {
        await AccountRepository.updateBalance(
            transaction.accountId,
            +transaction.amount,
        );
        await AccountRepository.updateBalance(
            transaction.toAccountId!,
            -transaction.amount,
        );
    }

    // Delete transaction
    await TransactionRepository.delete(transactionId, userId);

    // Recalculate budget if needed
    if (transaction.type === "expense" && transaction.categoryId) {
        await recalculateBudget(
            userId,
            transaction.categoryId,
            transaction.dateTimeISO,
        );
    }
}
```

#### Bulk Delete

```typescript
async function bulkDeleteTransactions(
    transactionIds: string[],
    userId: string,
) {
    // Validate all transactions belong to user
    const transactions = await TransactionRepository.getByIds(
        transactionIds,
        userId,
    );

    if (transactions.length !== transactionIds.length) {
        throw new Error("Some transactions not found");
    }

    // Delete each transaction (revert balance)
    for (const transaction of transactions) {
        await deleteTransaction(transaction.id, userId);
    }

    return { deletedCount: transactions.length };
}
```

### Acceptance Criteria

- ✅ User có thể tạo 3 loại transactions (income, expense, transfer)
- ✅ Form validation hoạt động đúng
- ✅ Account balance update real-time và chính xác
- ✅ Không thể expense/transfer khi insufficient balance
- ✅ Filters hoạt động (month, account, category, type, search)
- ✅ Search có thể tìm trong merchant, note, tags
- ✅ Bulk delete hoạt động, có confirm dialog
- ✅ Attachment upload thành công (image, PDF)
- ✅ Transaction list grouped by date, sorted desc
- ✅ Pagination hoạt động mượt mà
- ✅ Tags có thể add/remove dễ dàng

---

## 4. Categories

### Overview

Quản lý danh mục để phân loại transactions.

### User Stories

- **US-C1**: Là user, tôi muốn xem danh sách categories với icons và colors
- **US-C2**: Là user, tôi muốn tạo custom category mới
- **US-C3**: Là user, tôi muốn chỉnh sửa category (tên, icon, color)
- **US-C4**: Là user, tôi muốn xóa category không dùng (nếu không có transactions)
- **US-C5**: Là user, tôi muốn phân biệt expense categories và income categories
- **US-C6**: Là user, tôi muốn tạo sub-categories (parent-child relationship)

### Default Categories

**Income Categories:**

- 💰 Lương
- 🎁 Thưởng
- 📈 Đầu tư
- 💵 Thu nhập khác

**Expense Categories:**

- 🍜 Ăn uống
- 🚗 Di chuyển
- 🛍️ Mua sắm
- 🏠 Nhà cửa
- 🎮 Giải trí
- ⚕️ Sức khỏe
- 📚 Giáo dục
- 💸 Chi phí khác

### UI Components

#### 4.1 Category List Page

**Route**: `/settings/categories` (hoặc `/categories`)

**Layout:**

- Tabs: "Chi tiêu" | "Thu nhập"
- "Thêm danh mục" button
- Category grid (4 columns)

**Category Card:**

```
┌─────────────────────┐
│  🍜                  │
│                     │
│  Ăn uống            │
│  12 giao dịch       │
│                     │
│  [⋮]                │
└─────────────────────┘
```

Actions (⋮ menu):

- Chỉnh sửa
- Xem giao dịch
- Xóa (disabled nếu đang được dùng hoặc is_default)

#### 4.2 Create/Edit Category Modal

**Form Fields:**

- Name\* (text input)
- Type\* (expense | income | both)
- Icon\* (emoji picker)
- Color\* (color picker với preset colors)
- Parent Category (dropdown, optional, for sub-categories)

**Validation:**

- Name: required, 1-100 chars
- Type: required
- Icon: required
- Color: required, valid hex

**Emoji Picker:**
Categorized emojis:

- Food & Drink: 🍜🍕🍔🍰☕
- Transport: 🚗🚕🚌🚎🚲
- Shopping: 🛍️👕👗👠💄
- House: 🏠🛋️🛏️🚿💡
- Entertainment: 🎮🎬🎭🎪🎨
- Health: ⚕️💊💉🏥
- Education: 📚📖✏️🎓
- Money: 💰💵💳💸

**Color Picker:**
Preset colors (Tailwind palette):

- Red: #EF4444
- Orange: #F59E0B
- Yellow: #EAB308
- Green: #10B981
- Blue: #3B82F6
- Purple: #8B5CF6
- Pink: #EC4899
- Gray: #6B7280

#### 4.3 Sub-Categories

**Use Case:**
Chia Ăn uống thành:

- Ăn uống (parent)
    - Ăn sáng (child)
    - Ăn trưa (child)
    - Ăn tối (child)
    - Cà phê (child)

**Display:**
Category list shows hierarchy:

```
🍜 Ăn uống
  └ ☕ Cà phê
  └ 🍱 Ăn trưa
  └ 🍳 Ăn sáng
```

**Transaction Creation:**
When selecting category, show tree structure:

```
Danh mục ▾
  🍜 Ăn uống
    ☕ Cà phê
    🍱 Ăn trưa
    🍳 Ăn sáng
  🚗 Di chuyển
    ⛽ Xăng xe
    🅿️ Đỗ xe
```

### Business Logic

#### Create Category

```typescript
async function createCategory(input: CategoryInput): Promise<Category> {
    // Validation
    if (!input.name || input.name.length < 1) {
        throw new Error("Name is required");
    }

    if (!["income", "expense", "both"].includes(input.type)) {
        throw new Error("Invalid category type");
    }

    // Check parent exists if provided
    if (input.parentId) {
        const parent = await CategoryRepository.getById(input.parentId, userId);
        if (!parent) {
            throw new Error("Parent category not found");
        }

        // Ensure parent has same type
        if (parent.type !== input.type && parent.type !== "both") {
            throw new Error("Parent category type mismatch");
        }
    }

    // Create category
    const category = await CategoryRepository.create({
        ...input,
        userId,
        id: generateUUID(),
        isDefault: false,
        createdAt: new Date(),
        updatedAt: new Date(),
    });

    return category;
}
```

#### Delete Category

```typescript
async function deleteCategory(categoryId: string, userId: string) {
    const category = await CategoryRepository.getById(categoryId, userId);

    if (!category) {
        throw new Error("Category not found");
    }

    // Cannot delete default categories
    if (category.isDefault) {
        throw new Error("Cannot delete default category");
    }

    // Check if category is in use
    const transactionCount = await TransactionRepository.countByCategory(
        categoryId,
        userId,
    );

    if (transactionCount > 0) {
        throw new Error("Cannot delete category with existing transactions");
    }

    // Check if category has children
    const children = await CategoryRepository.getChildren(categoryId, userId);

    if (children.length > 0) {
        throw new Error("Cannot delete category with sub-categories");
    }

    await CategoryRepository.delete(categoryId, userId);
}
```

### Acceptance Criteria

- ✅ Default categories tạo tự động khi user đăng ký
- ✅ User có thể tạo custom categories
- ✅ Emoji picker và color picker hoạt động
- ✅ Sub-categories (parent-child) hoạt động
- ✅ Không thể xóa default categories
- ✅ Không thể xóa categories đang được dùng
- ✅ Category filter trong transactions hoạt động
- ✅ Category icons hiển thị trong transaction list

---

## 5. Budgets

### Overview

Đặt ngân sách cho từng category và theo dõi chi tiêu so với kế hoạch.

### User Stories

- **US-B1**: Là user, tôi muốn đặt ngân sách hàng tháng cho từng category
- **US-B2**: Là user, tôi muốn xem tiến trình chi tiêu so với budget (progress bar)
- **US-B3**: Là user, tôi muốn nhận cảnh báo khi sắp vượt budget (80%, 100%)
- **US-B4**: Là user, tôi muốn xem budget overview cho cả tháng
- **US-B5**: Là user, tôi muốn copy budgets sang tháng sau để tiện
- **US-B6**: Là user, tôi muốn so sánh budget vs actual qua nhiều tháng

### Budget Model

```typescript
interface Budget {
    id: string;
    userId: string;
    categoryId: string;
    amount: number; // Limit/ngân sách đặt ra
    spent: number; // Đã chi (calculated)
    month: string; // "2026-02"
    status: "normal" | "warning" | "exceeded";
    percentage: number; // (spent / amount) * 100
    createdAt: Date;
    updatedAt: Date;
}
```

### UI Components

#### 5.1 Budget List Page

**Route**: `/budgets`

**Layout:**

- Header: "Ngân sách" + "Thêm ngân sách" button + Month selector
- Overview card: Total budget, total spent, remaining
- Progress bar: Overall progress
- Budget cards list

**Overview Card:**

```
┌──────────────────────────────────────────┐
│ Tổng quan ngân sách - Tháng 2/2026       │
├──────────────────────────────────────────┤
│ Tổng ngân sách:     15,000,000 ₫         │
│ Đã chi:             12,500,000 ₫         │
│ Còn lại:             2,500,000 ₫         │
│                                          │
│ ████████████░░░░░░ 83%                   │
└──────────────────────────────────────────┘
```

**Budget Card:**

```
┌──────────────────────────────────────────┐
│ 🍜 Ăn uống                               │
│                                          │
│ 3,500,000 ₫ / 5,000,000 ₫                │
│ ███████░░░░░░░░░░ 70%          [⋮]       │
│                                          │
│ Còn lại: 1,500,000 ₫                     │
│ Trung bình: ~125,000 ₫/ngày              │
└──────────────────────────────────────────┘
```

**Status Colors:**

- Normal (<70%): Green progress bar
- Warning (70-100%): Orange progress bar
- Exceeded (>100%): Red progress bar

**Actions (⋮):**

- Chỉnh sửa ngân sách
- Xem giao dịch
- Xóa

**Empty State:**

```
No budgets set for this month

Đặt ngân sách để theo dõi chi tiêu!

[Thêm ngân sách đầu tiên]
```

#### 5.2 Add/Edit Budget Modal

**Form Fields:**

- Category\* (dropdown, chỉ expense categories)
- Amount\* (number input, VND)
- Month\* (month picker, default = current month)
- Rollover (checkbox): "Chuyển số dư sang tháng sau"

**Validation:**

- Category: required
- Amount: required, >0
- Một category chỉ có 1 budget per month

**Smart Suggestions:**
Show suggested budget based on:

```
Gợi ý ngân sách:
  Trung bình 3 tháng gần nhất: 4,200,000 ₫
  Tháng trước: 3,800,000 ₫

[Sử dụng gợi ý]
```

#### 5.3 Budget Detail Modal

**Trigger**: Click budget card

**Display:**

```
┌────────────────────────────────────────┐
│ Ngân sách Ăn uống - Tháng 2/2026       │
├────────────────────────────────────────┤
│ Ngân sách:      5,000,000 ₫            │
│ Đã chi:         3,500,000 ₫ (70%)      │
│ Còn lại:        1,500,000 ₫            │
│                                        │
│ ███████░░░░░░░░░░ 70%                  │
│                                        │
│ Chi tiết:                              │
│ • 25 giao dịch                         │
│ • Trung bình: 140,000 ₫/giao dịch      │
│ • Dự kiến hết vào: 05/03 (còn 6 ngày)  │
│                                        │
│ Giao dịch gần nhất:                    │
│ • 27/02 - Phở 24 - 150,000 ₫          │
│ • 26/02 - Highlands - 85,000 ₫        │
│ • 25/02 - Lotteria - 120,000 ₫        │
│                                        │
│ [Xem tất cả giao dịch]                 │
│                                        │
│           [Đóng]  [Chỉnh sửa]          │
└────────────────────────────────────────┘
```

#### 5.4 Budget Comparison View

**Route**: `/budgets/compare`

**Display**: Multi-month comparison table

```
Danh mục       | 12/2025    | 01/2026    | 02/2026
---------------------------------------------------------
Ăn uống        | 4.2M/5M    | 3.8M/5M    | 3.5M/5M
               | 84%        | 76%        | 70%
---------------------------------------------------------
Di chuyển      | 1.5M/2M    | 1.8M/2M    | 2.1M/2M ⚠️
               | 75%        | 90%        | 105%
---------------------------------------------------------
Mua sắm        | 2.0M/3M    | 2.5M/3M    | 1.8M/3M
               | 67%        | 83%        | 60%
```

**Insights:**

- Highlight trends: increasing/decreasing
- Flag exceeded budgets
- Show YoY comparison

### Business Logic

#### Create Budget

```typescript
async function createBudget(input: BudgetInput): Promise<Budget> {
    // Validation
    if (input.amount <= 0) {
        throw new Error("Budget amount must be greater than 0");
    }

    // Check category is expense type
    const category = await CategoryRepository.getById(input.categoryId, userId);
    if (
        !category ||
        (category.type !== "expense" && category.type !== "both")
    ) {
        throw new Error("Budget can only be set for expense categories");
    }

    // Check for existing budget
    const existing = await BudgetRepository.getByMonthCategory(
        userId,
        input.month,
        input.categoryId,
    );

    if (existing) {
        throw new Error("Budget already exists for this category and month");
    }

    // Calculate current spent
    const startDate = `${input.month}-01`;
    const endDate = lastDayOfMonth(input.month);

    const transactions = await TransactionRepository.list({
        userId,
        type: "expense",
        categoryId: input.categoryId,
        startDate,
        endDate,
    });

    const spent = transactions.reduce((sum, t) => sum + t.amount, 0);

    // Calculate status
    const percentage = (spent / input.amount) * 100;
    let status: BudgetStatus;

    if (percentage >= 100) {
        status = "exceeded";
    } else if (percentage >= 70) {
        status = "warning";
    } else {
        status = "normal";
    }

    // Create budget
    const budget = await BudgetRepository.create({
        ...input,
        userId,
        id: generateUUID(),
        spent,
        percentage,
        status,
        createdAt: new Date(),
        updatedAt: new Date(),
    });

    return budget;
}
```

#### Recalculate Budget on Transaction

```typescript
async function recalculateBudget(
    userId: string,
    categoryId: string,
    transactionDate: string,
) {
    const month = transactionDate.substring(0, 7); // "2026-02"

    const budget = await BudgetRepository.getByMonthCategory(
        userId,
        month,
        categoryId,
    );

    if (!budget) {
        return; // No budget set
    }

    // Recalculate spent
    const startDate = `${month}-01`;
    const endDate = lastDayOfMonth(month);

    const transactions = await TransactionRepository.list({
        userId,
        type: "expense",
        categoryId,
        startDate,
        endDate,
    });

    const spent = transactions.reduce((sum, t) => sum + t.amount, 0);
    const percentage = (spent / budget.amount) * 100;

    // Update status
    let status: BudgetStatus;
    if (percentage >= 100) {
        status = "exceeded";
    } else if (percentage >= 80) {
        status = "warning";
    } else {
        status = "normal";
    }

    // Update budget
    await BudgetRepository.update(budget.id, {
        spent,
        percentage,
        status,
        updatedAt: new Date(),
    });

    // Create alert if needed
    if (status === "warning" && budget.status !== "warning") {
        await AlertService.create({
            userId,
            type: "budget_warning",
            severity: "warning",
            title: `Ngân sách ${budget.category.name} sắp vượt`,
            message: `Bạn đã chi ${percentage.toFixed(1)}% ngân sách`,
            relatedId: budget.id,
            relatedType: "budget",
        });
    }

    if (status === "exceeded" && budget.status !== "exceeded") {
        await AlertService.create({
            userId,
            type: "budget_exceeded",
            severity: "danger",
            title: `Vượt ngân sách ${budget.category.name}`,
            message: `Bạn đã chi vượt ${(percentage - 100).toFixed(1)}%`,
            relatedId: budget.id,
            relatedType: "budget",
        });
    }
}
```

#### Copy Budgets to Next Month

```typescript
async function copyBudgetsToNextMonth(userId: string, fromMonth: string) {
    const budgets = await BudgetRepository.list({ userId, month: fromMonth });

    const toMonth = addMonths(fromMonth, 1);

    const created = [];

    for (const budget of budgets) {
        // Check if budget already exists
        const existing = await BudgetRepository.getByMonthCategory(
            userId,
            toMonth,
            budget.categoryId,
        );

        if (!existing) {
            const newBudget = await createBudget({
                userId,
                categoryId: budget.categoryId,
                amount: budget.amount, // Keep same amount
                month: toMonth,
            });

            created.push(newBudget);
        }
    }

    return created;
}
```

### Acceptance Criteria

- ✅ User có thể tạo budget cho expense categories
- ✅ Budget progress tính chính xác real-time
- ✅ Status colors (green/orange/red) hiển thị đúng
- ✅ Alerts tạo tự động khi đạt 80% và 100%
- ✅ Không thể tạo duplicate budget (same category + month)
- ✅ Budget suggestions dựa trên historical data
- ✅ Copy budgets to next month hoạt động
- ✅ Multi-month comparison view hiển thị trends
- ✅ Budget recalculates khi add/edit/delete transaction

---

## 6. Reports & Analytics

### Overview

Báo cáo và phân tích chi tiết về tình hình tài chính.

### User Stories

- **US-R1**: Là user, tôi muốn xem báo cáo theo ngày/tuần/tháng/năm
- **US-R2**: Là user, tôi muốn xem phân tích chi tiêu theo categories (pie chart)
- **US-R3**: Là user, tôi muốn xem trend thu/chi qua các tháng (line chart)
- **US-R4**: Là user, tôi muốn so sánh tháng này vs tháng trước
- **US-R5**: Là user, tôi muốn export báo cáo ra PDF/Excel
- **US-R6**: Là user, tôi muốn xem top merchants chi tiêu nhiều nhất
- **US-R7**: Là user, tôi muốn phân tích cash flow (money in vs money out)

### Report Types

#### 6.1 Overview Report

**Route**: `/reports/overview`

**Time Range Selector:**

- Last 7 days
- This month
- Last month
- This year
- Custom range

**Sections:**

**A. Summary Cards**

- Total Income
- Total Expense
- Net Saving
- Saving Rate (%)

**B. Income vs Expense Chart**
Line chart, 2 lines:

- Income (green line)
- Expense (red line)
- X-axis: Time (days/weeks/months)
- Y-axis: Amount (VND)

**C. Category Breakdown**
Pie chart hoặc Donut chart:

- Each slice = một category
- Size = percentage of total expense
- Colors = category colors
- Click slice → drill down to transactions

**D. Top Categories Table**

```
Danh mục          Số tiền        % Tổng chi    Transactions
---------------------------------------------------------------
Ăn uống          5,200,000 ₫     35%           48 giao dịch
Di chuyển        2,800,000 ₫     19%           25 giao dịch
Mua sắm          2,100,000 ₫     14%           12 giao dịch
```

**E. Month-over-Month Comparison**

```
Metric              Tháng này   Tháng trước   Thay đổi
-----------------------------------------------------------
Thu nhập           15,000,000   14,500,000    +3.4% ↑
Chi tiêu           12,000,000   11,800,000    +1.7% ↑
Tiết kiệm           3,000,000    2,700,000   +11.1% ↑
Tỷ lệ tiết kiệm         20.0%        18.6%    +1.4pp ↑
```

#### 6.2 Category Report

**Route**: `/reports/by-category`

**Filters:**

- Category selector (multi-select)
- Time range
- Account filter

**Display:**

**A. Category Comparison Bar Chart**
Horizontal bars, sorted by amount desc:

```
Ăn uống      ████████████████ 5,200,000 ₫
Di chuyển    ████████ 2,800,000 ₫
Mua sắm      ██████ 2,100,000 ₫
```

**B. Category Trends**
Line chart:

- X-axis: Months (last 6 months)
- Y-axis: Amount
- Multiple lines (1 per category)
- Legend with colors

**C. Transaction Details Table**
For selected category:

```
Ngày        Merchant       Số tiền       Tài khoản
-------------------------------------------------------
27/02/26    Phở 24        150,000 ₫     Ví tiền mặt
26/02/26    Highlands      85,000 ₫     Vietcombank
```

#### 6.3 Account Report

**Route**: `/reports/by-account`

**Display:**

**A. Account Balance Chart**
Stacked bar chart:

- X-axis: Accounts
- Y-axis: Balance
- Colors by account type

**B. Account Activity**
For each account:

- Total inflows (income + transfers in)
- Total outflows (expense + transfers out)
- Net change
- Current balance

**C. Account Trends**
Line chart showing balance history over time

#### 6.4 Cash Flow Report

**Route**: `/reports/cash-flow`

**Display:**

**A. Cash Flow Chart**
Waterfall chart or stacked bar:

```
Starting Balance
  + Income
  - Fixed Expenses (rent, utilities)
  - Variable Expenses (food, shopping)
  = Ending Balance
```

**B. Cash Flow Statement**

```
Dòng tiền tháng 2/2026

Số dư đầu kỳ:           10,000,000 ₫

Dòng tiền thu:
  Thu nhập              15,000,000 ₫
  Chuyển khoản vào       1,000,000 ₫
  ───────────────────────────────────
  Tổng thu:             16,000,000 ₫

Dòng tiền chi:
  Chi phí cố định        5,000,000 ₫
  Chi phí biến đổi       7,000,000 ₫
  Chuyển khoản ra          500,000 ₫
  ───────────────────────────────────
  Tổng chi:             12,500,000 ₫

Thay đổi ròng:          +3,500,000 ₫
Số dư cuối kỳ:          13,500,000 ₫
```

#### 6.5 Merchant Report

**Route**: `/reports/by-merchant`

**Display:**

**A. Top Merchants Table**

```
Cửa hàng           Tổng chi      Số lần     Trung bình
---------------------------------------------------------------
Highlands Coffee   1,200,000 ₫    15          80,000 ₫
Shopee             2,500,000 ₫     8         312,500 ₫
Circle K             450,000 ₫    12          37,500 ₫
```

**B. Merchant Frequency Chart**
Bar chart: Số lần mua hàng tại mỗi merchant

**C. Merchant Spending Trend**
Line chart: Chi tiêu tại top merchants qua các tháng

### Export Features

#### PDF Export

**Contents:**

- Report header (title, time range, generated date)
- Summary section with key metrics
- Charts (rendered as images)
- Tables (formatted)
- Footer (app branding)

**Implementation:**

```typescript
async function exportReportPDF(reportData: ReportData): Promise<Blob> {
    const pdf = new jsPDF();

    // Header
    pdf.setFontSize(20);
    pdf.text("Báo cáo tài chính", 20, 20);
    pdf.setFontSize(12);
    pdf.text(`Kỳ báo cáo: ${reportData.timeRange}`, 20, 30);

    // Summary
    pdf.text(`Tổng thu: ${formatCurrency(reportData.totalIncome)}`, 20, 50);
    pdf.text(`Tổng chi: ${formatCurrency(reportData.totalExpense)}`, 20, 60);
    pdf.text(`Tiết kiệm: ${formatCurrency(reportData.netSaving)}`, 20, 70);

    // Chart (convert to image)
    const chartCanvas = document.querySelector("#expense-chart canvas");
    const chartImage = chartCanvas.toDataURL("image/png");
    pdf.addImage(chartImage, "PNG", 20, 90, 170, 100);

    // Table
    pdf.autoTable({
        head: [["Danh mục", "Số tiền", "% Tổng"]],
        body: reportData.categories.map((c) => [
            c.name,
            formatCurrency(c.amount),
            `${c.percentage}%`,
        ]),
        startY: 200,
    });

    return pdf.output("blob");
}
```

#### Excel Export

**Implementation:**

```typescript
async function exportReportExcel(reportData: ReportData): Promise<Blob> {
    const workbook = XLSX.utils.book_new();

    // Summary sheet
    const summaryData = [
        ["Metric", "Value"],
        ["Total Income", reportData.totalIncome],
        ["Total Expense", reportData.totalExpense],
        ["Net Saving", reportData.netSaving],
        ["Saving Rate", `${reportData.savingRate}%`],
    ];
    const summarySheet = XLSX.utils.aoa_to_sheet(summaryData);
    XLSX.utils.book_append_sheet(workbook, summarySheet, "Summary");

    // Transactions sheet
    const transactionsData = [
        ["Date", "Type", "Category", "Merchant", "Amount", "Account"],
        ...reportData.transactions.map((t) => [
            t.date,
            t.type,
            t.category,
            t.merchant,
            t.amount,
            t.account,
        ]),
    ];
    const transactionsSheet = XLSX.utils.aoa_to_sheet(transactionsData);
    XLSX.utils.book_append_sheet(workbook, transactionsSheet, "Transactions");

    // Categories sheet
    const categoriesData = [
        ["Category", "Amount", "Percentage", "Count"],
        ...reportData.categories.map((c) => [
            c.name,
            c.amount,
            `${c.percentage}%`,
            c.count,
        ]),
    ];
    const categoriesSheet = XLSX.utils.aoa_to_sheet(categoriesData);
    XLSX.utils.book_append_sheet(workbook, categoriesSheet, "Categories");

    // Generate Excel file
    const excelBuffer = XLSX.write(workbook, {
        type: "array",
        bookType: "xlsx",
    });

    return new Blob([excelBuffer], {
        type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    });
}
```

### Business Logic

#### Generate Overview Report

```typescript
async function generateOverviewReport(
    userId: string,
    startDate: string,
    endDate: string,
): Promise<OverviewReport> {
    // Get all transactions in range
    const transactions = await TransactionRepository.list({
        userId,
        startDate,
        endDate,
    });

    // Calculate totals
    const incomes = transactions.filter((t) => t.type === "income");
    const expenses = transactions.filter((t) => t.type === "expense");

    const totalIncome = sum(incomes.map((t) => t.amount));
    const totalExpense = sum(expenses.map((t) => t.amount));
    const netSaving = totalIncome - totalExpense;
    const savingRate = totalIncome > 0 ? (netSaving / totalIncome) * 100 : 0;

    // Group by category
    const byCategory = groupBy(expenses, "categoryId");
    const categories = Object.entries(byCategory)
        .map(([categoryId, txs]) => {
            const amount = sum(txs.map((t) => t.amount));
            return {
                categoryId,
                name: txs[0].category.name,
                icon: txs[0].category.icon,
                amount,
                percentage: (amount / totalExpense) * 100,
                count: txs.length,
            };
        })
        .sort((a, b) => b.amount - a.amount);

    // Spending trend (group by day/week/month)
    const interval = calculateInterval(startDate, endDate); // day/week/month
    const trend = groupTransactionsByInterval(expenses, interval);

    // Previous period comparison
    const previousPeriod = await generateOverviewReport(
        userId,
        subtractPeriod(startDate, endDate),
        startDate,
    );
    const comparison = calculateComparison(
        { totalIncome, totalExpense, netSaving, savingRate },
        previousPeriod,
    );

    return {
        timeRange: { startDate, endDate },
        summary: {
            totalIncome,
            totalExpense,
            netSaving,
            savingRate,
            transactionCount: transactions.length,
        },
        categories,
        trend,
        comparison,
        generatedAt: new Date(),
    };
}
```

### Acceptance Criteria

- ✅ Reports load within 3 seconds for 1 year of data
- ✅ All charts interactive (tooltips, click to drill down)
- ✅ Time range selector hoạt động (predefined + custom)
- ✅ Category pie chart colors match category colors
- ✅ Month-over-month comparison tính chính xác
- ✅ PDF export includes all charts and tables
- ✅ Excel export có multiple sheets (summary, transactions, categories)
- ✅ Reports responsive trên mobile
- ✅ Empty states shown khi không có data

---

## 7. AI Chat Assistant

### Overview

Trợ lý AI chatbot giúp user phân tích tài chính và trả lời câu hỏi.

### User Stories

- **US-AI1**: Là user, tôi muốn hỏi chatbot về tình hình tài chính ("Tháng này tôi chi bao nhiêu?")
- **US-AI2**: Là user, tôi muốn nhận insights tự động ("Chi tiêu tăng đột biến ở đâu?")
- **US-AI3**: Là user, tôi muốn hỏi về budgets ("Budget nào sắp vượt?")
- **US-AI4**: Là user, tôi muốn được gợi ý cách tiết kiệm
- **US-AI5**: Là user, tôi muốn chat history được lưu lại

### AI Capabilities

#### Financial Queries

User có thể hỏi:

- "Tháng này tôi chi bao nhiêu?"
- "Top 3 categories tôi chi nhiều nhất?"
- "So sánh chi tiêu tháng này vs tháng trước"
- "Tôi tiết kiệm được bao nhiêu tháng này?"
- "Account nào có số dư cao nhất?"

#### Budget Analysis

- "Budget nào đang vượt?"
- "Tôi cần cắt giảm ở đâu để không vượt ngân sách?"
- "Nếu tôi chi tiếp như này thì cuối tháng có vượt budget không?"

#### Insights & Recommendations

- "Phát hiện chi tiêu bất thường nào không?"
- "Tôi nên đặt budget bao nhiêu cho category X?"
- "Gợi ý cách tiết kiệm cho tôi"
- "Phân tích patterns chi tiêu của tôi"

#### Forecasting

- "Predict chi tiêu tháng sau"
- "Khi nào thì tôi tiết kiệm được 10 triệu?"

### UI Components

#### 7.1 Chat Interface

**Route**: `/ai-chat`

**Layout:**

```
┌────────────────────────────────────────────┐
│ 🤖 Trợ lý tài chính AI           [⋮]       │
├────────────────────────────────────────────┤
│                                            │
│  ┌──────────────────────────────┐          │
│  │ Chào bạn! Tôi có thể giúp     │          │
│  │ bạn phân tích tài chính.      │          │
│  │ Hãy hỏi tôi bất cứ điều gì!   │          │
│  └──────────────────────────────┘          │
│                                    🤖      │
│                                            │
│              ┌────────────────────────┐    │
│          👤  │ Tháng này tôi chi      │    │
│              │ bao nhiêu?             │    │
│              └────────────────────────┘    │
│                                            │
│  ┌──────────────────────────────┐          │
│  │ Tháng 2/2026 bạn đã chi:      │          │
│  │                               │          │
│  │ 💸 12,500,000 ₫               │          │
│  │                               │          │
│  │ Top 3 categories:             │          │
│  │ 1. Ăn uống - 5.2M (42%)       │          │
│  │ 2. Di chuyển - 2.8M (22%)     │          │
│  │ 3. Mua sắm - 2.1M (17%)       │          │
│  │                               │          │
│  │ [Xem chi tiết]                │          │
│  └──────────────────────────────┘          │
│                                    🤖      │
│                                            │
├────────────────────────────────────────────┤
│  [Type your message...]           [Send]   │
└────────────────────────────────────────────┘
```

**Features:**

- Message bubbles (user right, AI left)
- Typing indicator when AI is thinking
- Quick action buttons in AI responses
- Scroll to bottom button
- Chat history loads on mount

**Quick Suggestions** (shown when chat is empty):

```
Gợi ý câu hỏi:
- Tổng quan tháng này
- Budget nào sắp vượt?
- So sánh với tháng trước
- Chi tiêu bất thường
- Gợi ý tiết kiệm
```

#### 7.2 AI Response Types

**A. Text Response**
Simple text answer

**B. Data Card**
Structured data display:

```
┌─────────────────────────────────┐
│ 📊 Chi tiêu tháng 2/2026         │
├─────────────────────────────────┤
│ Tổng chi: 12,500,000 ₫           │
│ Số giao dịch: 87                │
│ Trung bình/ngày: 446,428 ₫       │
│                                 │
│ [Xem chi tiết] [Xem báo cáo]    │
└─────────────────────────────────┘
```

**C. Chart Response**
Inline charts (bar, pie, line)

**D. Action Buttons**

- "Xem chi tiết" → navigate to related page
- "Thiết lập ngân sách" → open budget modal
- "Xem giao dịch" → filter transactions

#### 7.3 Chat History

**Storage**: MongoDB collection `chat_messages`

**Schema:**

```typescript
interface ChatMessage {
    id: string;
    userId: string;
    role: "user" | "assistant";
    content: string;
    metadata?: {
        query_type?: string; // "financial_query" | "budget_analysis" | etc
        data?: any; // Structured data for rendering
    };
    createdAt: Date;
}
```

**Load History:**

- Load last 50 messages on mount
- Infinite scroll to load more (pagination)
- Group by conversation/session (optional)

### Business Logic

#### Process User Query

```typescript
async function processAIQuery(
    userId: string,
    query: string,
): Promise<AIResponse> {
    // Save user message
    await ChatMessageRepository.create({
        userId,
        role: "user",
        content: query,
        createdAt: new Date(),
    });

    // Classify query intent
    const intent = await classifyIntent(query);

    let response: AIResponse;

    switch (intent.type) {
        case "financial_summary":
            response = await generateFinancialSummary(userId, intent.params);
            break;

        case "budget_analysis":
            response = await analyzeBudgets(userId, intent.params);
            break;

        case "spending_comparison":
            response = await compareSpending(userId, intent.params);
            break;

        case "insights":
            response = await generateInsights(userId, intent.params);
            break;

        case "forecast":
            response = await generateForecast(userId, intent.params);
            break;

        default:
            response = {
                text: "Xin lỗi, tôi chưa hiểu câu hỏi của bạn. Bạn có thể hỏi về chi tiêu, ngân sách, hoặc tài khoản.",
                type: "text",
            };
    }

    // Save assistant message
    await ChatMessageRepository.create({
        userId,
        role: "assistant",
        content: response.text,
        metadata: response.metadata,
        createdAt: new Date(),
    });

    return response;
}
```

#### Intent Classification

```typescript
function classifyIntent(query: string): Intent {
    const q = query.toLowerCase();

    // Financial summary patterns
    if (/(tháng này|this month).*(chi|spend|expense)/i.test(q)) {
        return {
            type: "financial_summary",
            params: { period: "current_month" },
        };
    }

    // Budget patterns
    if (/(budget|ngân sách).*(vượt|exceed|over)/i.test(q)) {
        return {
            type: "budget_analysis",
            params: { checkExceeded: true },
        };
    }

    // Comparison patterns
    if (/(so sánh|compare).*(tháng trước|last month)/i.test(q)) {
        return {
            type: "spending_comparison",
            params: {
                current: "current_month",
                previous: "last_month",
            },
        };
    }

    // Insights patterns
    if (/(phát hiện|detect|unusual|bất thường)/i.test(q)) {
        return {
            type: "insights",
            params: { analysisType: "unusual_spending" },
        };
    }

    // Forecast patterns
    if (/(predict|dự đoán|tháng sau|next month)/i.test(q)) {
        return {
            type: "forecast",
            params: { period: "next_month" },
        };
    }

    return { type: "unknown", params: {} };
}
```

#### Generate Financial Summary

```typescript
async function generateFinancialSummary(
    userId: string,
    params: any,
): Promise<AIResponse> {
    const month =
        params.period === "current_month" ? getCurrentMonth() : params.month;

    const startDate = `${month}-01`;
    const endDate = lastDayOfMonth(month);

    // Get data
    const report = await ReportService.generateOverview(
        userId,
        startDate,
        endDate,
    );

    // Format response
    const text = `
Tháng ${month} bạn đã chi:

💸 ${formatCurrency(report.summary.totalExpense)}

Top 3 categories:
${report.categories
    .slice(0, 3)
    .map(
        (c, i) =>
            `${i + 1}. ${c.name} - ${formatCurrency(c.amount)} (${c.percentage.toFixed(1)}%)`,
    )
    .join("\n")}

${
    report.comparison.expense > 0
        ? `⚠️ Tăng ${report.comparison.expense.toFixed(1)}% so với tháng trước`
        : `✅ Giảm ${Math.abs(report.comparison.expense).toFixed(1)}% so với tháng trước`
}
  `.trim();

    return {
        text,
        type: "data_card",
        metadata: {
            query_type: "financial_summary",
            data: report,
            actions: [
                { label: "Xem chi tiết", link: "/reports/overview" },
                {
                    label: "Xem giao dịch",
                    link: `/transactions?month=${month}`,
                },
            ],
        },
    };
}
```

#### Analyze Budgets

```typescript
async function analyzeBudgets(
    userId: string,
    params: any,
): Promise<AIResponse> {
    const month = getCurrentMonth();
    const budgets = await BudgetService.list({ userId, month });

    // Find exceeded/warning budgets
    const exceeded = budgets.filter((b) => b.status === "exceeded");
    const warnings = budgets.filter((b) => b.status === "warning");

    if (exceeded.length === 0 && warnings.length === 0) {
        return {
            text: "✅ Tất cả ngân sách đều ổn! Bạn đang chi tiêu hợp lý.",
            type: "text",
        };
    }

    let text = "";

    if (exceeded.length > 0) {
        text += "⚠️ Các ngân sách đã vượt:\n\n";
        exceeded.forEach((b) => {
            text += `• ${b.category.name}: ${formatCurrency(b.spent)} / ${formatCurrency(b.amount)} (${b.percentage.toFixed(1)}%)\n`;
        });
        text += "\n";
    }

    if (warnings.length > 0) {
        text += "⚡ Các ngân sách sắp vượt:\n\n";
        warnings.forEach((b) => {
            text += `• ${b.category.name}: ${formatCurrency(b.spent)} / ${formatCurrency(b.amount)} (${b.percentage.toFixed(1)}%)\n`;
        });
    }

    text += "\n💡 Gợi ý: Hãy xem xét cắt giảm chi tiêu ở các categories này.";

    return {
        text,
        type: "text",
        metadata: {
            query_type: "budget_analysis",
            data: { exceeded, warnings },
            actions: [{ label: "Xem ngân sách", link: "/budgets" }],
        },
    };
}
```

### OpenAI Integration (Optional)

Nếu muốn AI thông minh hơn, integrate OpenAI GPT:

```typescript
async function generateAIResponse(
    userId: string,
    query: string,
): Promise<AIResponse> {
    // Get user's financial context
    const context = await getUserFinancialContext(userId);

    // Call OpenAI API
    const completion = await openai.chat.completions.create({
        model: "gpt-4",
        messages: [
            {
                role: "system",
                content: `You are a financial advisor assistant. 
        User's financial data: ${JSON.stringify(context)}.
        Answer questions about their finances in Vietnamese.
        Be concise and helpful.`,
            },
            {
                role: "user",
                content: query,
            },
        ],
    });

    const text = completion.choices[0].message.content;

    return { text, type: "text" };
}
```

### Acceptance Criteria

- ✅ Chat interface responsive và smooth
- ✅ Intent classification hoạt động cho common queries
- ✅ AI responses accurate dựa trên real data
- ✅ Quick action buttons navigate đúng
- ✅ Chat history loads và pagination hoạt động
- ✅ Typing indicator xuất hiện khi AI đang process
- ✅ Messages group by conversation
- ✅ Có quick suggestions khi chat empty
- ✅ Inline charts render correctly

---

## 8. Alerts & Insights

### Overview

Hệ thống cảnh báo và insights tự động giúp user nhận biết vấn đề tài chính.

### User Stories

- **US-AL1**: Là user, tôi muốn nhận alert khi budget sắp vượt (80%)
- **US-AL2**: Là user, tôi muốn nhận alert khi phát hiện unusual spending
- **US-AL3**: Là user, tôi muốn nhận insight về patterns chi tiêu
- **US-AL4**: Là user, tôi muốn dismiss hoặc snooze alerts
- **US-AL5**: Là user, tôi muốn xem alert history

### Alert Types

#### 1. Budget Alerts

**Trigger Conditions:**

- Budget reaches 80% → Warning alert
- Budget reaches 100% → Danger alert
- Budget projected to exceed → Forecast alert

**Alert Data:**

```typescript
{
  type: "budget_warning",
  severity: "warning",
  title: "Ngân sách Ăn uống sắp vượt",
  message: "Bạn đã chi 85% ngân sách",
  relatedId: budgetId,
  relatedType: "budget",
  actionLabel: "Xem chi tiết",
  actionLink: "/budgets"
}
```

#### 2. Unusual Spending Alerts

**Trigger Conditions:**

- Spending in a category increases >50% vs last month
- Large transaction detected (>2x average)
- Multiple transactions in short time

**Alert Data:**

```typescript
{
  type: "unusual_spending",
  severity: "info",
  title: "Chi tiêu bất thường phát hiện",
  message: "Chi Mua sắm tăng 65% so với tháng trước",
  relatedId: categoryId,
  relatedType: "category",
  actionLabel: "Xem giao dịch",
  actionLink: "/transactions?category=..."
}
```

#### 3. Low Balance Alerts

**Trigger Conditions:**

- Account balance < threshold (e.g., 500k)
- Account balance projected to run out

**Alert Data:**

```typescript
{
  type: "low_balance",
  severity: "warning",
  title: "Số dư tài khoản thấp",
  message: "Ví tiền mặt còn 350,000 ₫",
  relatedId: accountId,
  relatedType: "account",
  actionLabel: "Xem tài khoản",
  actionLink: "/accounts/:id"
}
```

#### 4. Forecast Alerts

**Trigger Conditions:**

- Forecast predicts budget exceed
- Forecast predicts negative balance

**Alert Data:**

```typescript
{
  type: "forecast_warning",
  severity: "warn",
  title: "Dự báo vượt ngân sách",
  message: "Nếu tiếp tục chi như này, bạn sẽ vượt 2.5M ₫",
  relatedId: forecastId,
  relatedType: "forecast",
  actionLabel: "Xem dự báo",
  actionLink: "/forecasts"
}
```

### Insight Types

#### 1. Spending Patterns

**Examples:**

- "Bạn thường chi nhiều nhất vào cuối tuần"
- "Chi tiêu Ăn uống tăng dần vào cuối tháng"
- "Bạn thường mua sắm online vào buổi tối"

#### 2. Saving Opportunities

**Examples:**

- "Nếu giảm chi Ăn uống 20%, bạn tiết kiệm được 1M/tháng"
- "Chi Grab tăng 40%, có thể cân nhắc đi xe bus?"
- "3 subscriptions không dùng phát hiện (Netflix, Spotify, Gym)"

#### 3. Comparison Insights

**Examples:**

- "Tháng này chi thấp hơn 15% so với trung bình 3 tháng"
- "Bạn đã tiết kiệm được nhiều hơn 2 tháng trước"
- "Chi Giải trí giảm đáng kể trong tháng này"

### UI Components

#### 8.1 Alerts Panel (Dashboard)

**Location**: Right sidebar on Dashboard

**Display**: List of active alerts, sorted by severity and date

**Alert Card:**

```
┌────────────────────────────────────┐
│ ⚠️ Ngân sách Ăn uống sắp vượt      │
│                                    │
│ Bạn đã chi 85% ngân sách           │
│                                    │
│ [Xem chi tiết]        [Dismiss] ✕  │
└────────────────────────────────────┘
```

**Severity Colors:**

- info: Blue border
- warning: Orange border
- danger: Red border

#### 8.2 Alerts Page

**Route**: `/alerts`

**Layout:**

- Tabs: "Active" | "History"
- Filter: severity, type, date range
- Alert list (same as dashboard panel)

**Bulk Actions:**

- Dismiss all
- Mark all as read

#### 8.3 Insights Panel

**Location**: Can be on Dashboard or separate page `/insights`

**Display**: Card-based layout with insight types

**Insight Card:**

```
┌────────────────────────────────────────┐
│ 💡 Spending Pattern                    │
├────────────────────────────────────────┤
│ Bạn thường chi nhiều nhất vào cuối tuần│
│                                        │
│ [Chart: Spending by day of week]      │
│                                        │
│ 💬 Tip: Hãy lập kế hoạch chi tiêu     │
│    cuối tuần để tránh chi quá nhiều   │
└────────────────────────────────────────┘
```

### Business Logic

#### Create Budget Alert

```typescript
async function checkBudgetAlerts(userId: string) {
    const month = getCurrentMonth();
    const budgets = await BudgetService.list({ userId, month });

    for (const budget of budgets) {
        const percentage = budget.percentage;

        // Check for 80% warning
        if (percentage >= 80 && percentage < 100) {
            const existingAlert = await AlertRepository.findOne({
                userId,
                type: "budget_warning",
                relatedId: budget.id,
                status: "active",
            });

            if (!existingAlert) {
                await AlertRepository.create({
                    userId,
                    type: "budget_warning",
                    severity: "warning",
                    title: `Ngân sách ${budget.category.name} sắp vượt`,
                    message: `Bạn đã chi ${percentage.toFixed(1)}% ngân sách`,
                    relatedId: budget.id,
                    relatedType: "budget",
                    actionLabel: "Xem chi tiết",
                    actionLink: `/budgets`,
                    status: "active",
                    createdAt: new Date(),
                });
            }
        }

        // Check for 100% exceeded
        if (percentage >= 100) {
            const existingAlert = await AlertRepository.findOne({
                userId,
                type: "budget_exceeded",
                relatedId: budget.id,
                status: "active",
            });

            if (!existingAlert) {
                await AlertRepository.create({
                    userId,
                    type: "budget_exceeded",
                    severity: "danger",
                    title: `Vượt ngân sách ${budget.category.name}`,
                    message: `Bạn đã chi vượt ${(percentage - 100).toFixed(1)}%`,
                    relatedId: budget.id,
                    relatedType: "budget",
                    actionLabel: "Xem chi tiết",
                    actionLink: `/budgets`,
                    status: "active",
                    createdAt: new Date(),
                });
            }
        }
    }
}
```

#### Detect Unusual Spending

```typescript
async function detectUnusualSpending(userId: string) {
    const currentMonth = getCurrentMonth();
    const lastMonth = getPreviousMonth();

    // Get current month expenses by category
    const currentExpenses = await TransactionRepository.aggregateByCategory({
        userId,
        type: "expense",
        month: currentMonth,
    });

    // Get last month expenses by category
    const lastExpenses = await TransactionRepository.aggregateByCategory({
        userId,
        type: "expense",
        month: lastMonth,
    });

    for (const current of currentExpenses) {
        const last = lastExpenses.find(
            (e) => e.categoryId === current.categoryId,
        );

        if (!last) continue;

        const increasePercentage =
            ((current.amount - last.amount) / last.amount) * 100;

        // Alert if increase > 50%
        if (increasePercentage > 50) {
            await AlertRepository.create({
                userId,
                type: "unusual_spending",
                severity: "info",
                title: "Chi tiêu bất thường phát hiện",
                message: `Chi ${current.category.name} tăng ${increasePercentage.toFixed(1)}% so với tháng trước`,
                relatedId: current.categoryId,
                relatedType: "category",
                actionLabel: "Xem giao dịch",
                actionLink: `/transactions?category=${current.categoryId}&month=${currentMonth}`,
                status: "active",
                createdAt: new Date(),
            });
        }
    }
}
```

#### Generate Spending Pattern Insight

```typescript
async function generateSpendingPatterns(userId: string): Promise<Insight[]> {
    const insights: Insight[] = [];

    // Analyze spending by day of week
    const last30Days = getLast30Days();
    const transactions = await TransactionRepository.list({
        userId,
        type: "expense",
        startDate: last30Days[0],
        endDate: last30Days[29],
    });

    const byDayOfWeek = groupBy(transactions, (t) =>
        getDayOfWeek(t.dateTimeISO),
    );

    const weekendSpending = (byDayOfWeek[6] || []).concat(byDayOfWeek[0] || []);
    const weekdaySpending = [1, 2, 3, 4, 5].flatMap(
        (d) => byDayOfWeek[d] || [],
    );

    const avgWeekend = average(weekendSpending.map((t) => t.amount));
    const avgWeekday = average(weekdaySpending.map((t) => t.amount));

    if (avgWeekend > avgWeekday * 1.5) {
        insights.push({
            type: "spending_pattern",
            title: "Chi tiêu cuối tuần cao",
            message:
                "Bạn thường chi nhiều hơn 50% vào cuối tuần so với ngày thường",
            recommendation:
                "Hãy lập kế hoạch chi tiêu cuối tuần để tránh chi quá nhiều",
            data: { avgWeekend, avgWeekday },
        });
    }

    // More pattern analysis...

    return insights;
}
```

### Acceptance Criteria

- ✅ Budget alerts trigger at 80% and 100%
- ✅ Unusual spending detection hoạt động accurate
- ✅ Alerts dismissable và có history
- ✅ Alert notifications real-time (hoặc polling)
- ✅ Insights helpful và actionable
- ✅ Alert cards hiển thị severity colors
- ✅ Click action button navigate đúng
- ✅ Bulk dismiss hoạt động

---

## 9. Settings

### Overview

Cấu hình app và preferences của user.

### User Stories

- **US-S1**: Là user, tôi muốn xem và chỉnh sửa profile (name, email, avatar)
- **US-S2**: Là user, tôi muốn đổi password
- **US-S3**: Là user, tôi muốn cấu hình notification preferences
- **US-S4**: Là user, tôi muốn chọn currency mặc định (VND, USD)
- **US-S5**: Là user, tôi muốn export toàn bộ data (GDPR compliance)
- **US-S6**: Là user, tôi muốn delete account
- **US-S7**: Là user, tôi muốn quản lý categories custom

### Settings Sections

#### 9.1 Profile Settings

**Route**: `/settings/profile`

**Fields:**

- Avatar (upload image)
- Full Name
- Email (verified)
- Phone (optional)

**Actions:**

- Update profile
- Change avatar

#### 9.2 Security Settings

**Route**: `/settings/security`

**Sections:**

- Change Password
- Two-Factor Authentication (future)
- Active Sessions (future)

**Change Password Form:**

- Current password
- New password
- Confirm new password

**Validation:**

- New password ≥ 8 characters
- Must include: uppercase, lowercase, number
- Confirm password matches

#### 9.3 Notification Settings

**Route**: `/settings/notifications`

**Options:**

```
Email Notifications:
  ☑ Budget alerts
  ☑ Unusual spending detected
  ☑ Monthly summary report
  ☐ Weekly digest

Push Notifications: (future)
  ☑ Budget exceeded
  ☑ Low balance warnings
  ☐ Transaction reminders

In-App Notifications:
  ☑ All alerts
  ☑ Insights
```

#### 9.4 Preferences

**Route**: `/settings/preferences`

**Options:**

- Default Currency: [VND ▾]
- Language: [Tiếng Việt ▾]
- Date Format: [DD/MM/YYYY ▾]
- Number Format: [1,000,000.00 ▾]
- Theme: [Light | Dark | Auto]
- Start of Week: [Monday ▾]
- Fiscal Month Start: [1 ▾] (for custom fiscal calendar)

#### 9.5 Data & Privacy

**Route**: `/settings/data`

**Sections:**

**Export Data:**

```
Xuất toàn bộ dữ liệu của bạn

File bao gồm: transactions, accounts, budgets, categories

[Xuất dữ liệu (JSON)] [Xuất dữ liệu (CSV)]
```

**Delete Account:**

```
Xóa tài khoản

⚠️ Hành động này không thể hoàn tác!

Tất cả dữ liệu của bạn sẽ bị xóa vĩnh viễn.

[Tôi hiểu, xóa tài khoản của tôi]
```

#### 9.6 About

**Route**: `/settings/about`

**Display:**

- App version
- Privacy Policy link
- Terms of Service link
- Contact support

#### 9.7 Categories Management

**Route**: `/settings/categories`

Quản lý categories để phân loại transactions. Đây là nơi tập trung để user tạo, chỉnh sửa, và tổ chức categories của mình.

**Layout:**

- Tabs: "Chi tiêu" | "Thu nhập"
- "Thêm danh mục" button (top right)
- Category grid (4 columns desktop, 2 columns tablet, 1 column mobile)
- Search bar để tìm categories

**Category Card:**

```
┌─────────────────────────┐
│  🍜                      │
│                         │
│  Ăn uống                │
│  42 giao dịch           │
│  3,500,000 ₫            │
│                         │
│  [Chỉnh sửa]   [⋮]      │
└─────────────────────────┘
```

**Card Actions (⋮ menu):**

- Chỉnh sửa danh mục
- Xem giao dịch
- Xóa (disabled nếu đang được dùng hoặc is_default)
- Set as default (cho custom categories)

**Add/Edit Category Modal:**

```
┌─────────────────────────────────────┐
│ Thêm danh mục mới          [✕]      │
├─────────────────────────────────────┤
│                                     │
│ Tên danh mục *                      │
│ [Nhập tên danh mục]                 │
│                                     │
│ Loại *                              │
│ ○ Chi tiêu  ○ Thu nhập  ○ Cả hai   │
│                                     │
│ Icon *                              │
│ [🍜] [Icon picker...]               │
│                                     │
│ Màu sắc *                           │
│ ████ [Color picker...]              │
│                                     │
│ Danh mục cha (tùy chọn)             │
│ [Chọn danh mục cha ▾]               │
│                                     │
│ Mô tả (tùy chọn)                    │
│ [Nhập mô tả...]                     │
│                                     │
│         [Hủy]      [Lưu]            │
└─────────────────────────────────────┘
```

**Icon Picker:**
Categorized emoji picker với các nhóm:

- 🍽️ Food & Drink: 🍜🍕🍔🍰☕🍺🍱
- 🚗 Transport: 🚗🚕🚌🚎🚲🚇✈️
- 🛍️ Shopping: 🛍️👕👗👠💄📱💻
- 🏠 House: 🏠🛋️🛏️🚿💡🔧
- 🎮 Entertainment: 🎮🎬🎭🎪🎨🎸📺
- ⚕️ Health: ⚕️💊💉🏥🏋️‍♂️🧘
- 📚 Education: 📚📖✏️🎓💼
- 💰 Money: 💰💵💳💸📈📊
- 🎁 Other: 🎁🎉🎂📦🔔

**Color Picker:**
Preset colors với visual preview:

```
[🔴] [🟠] [🟡] [🟢] [🔵] [🟣] [🟤] [⚫]
Red  Orange Yellow Green Blue Purple Brown Gray
```

**Predefined Colors:**

- Red: #EF4444
- Orange: #F59E0B
- Yellow: #EAB308
- Green: #10B981
- Blue: #3B82F6
- Purple: #8B5CF6
- Pink: #EC4899
- Gray: #6B7280

**Sub-Categories (Hierarchy View):**

Khi "Show hierarchy" được bật:

```
┌─────────────────────────────────────┐
│ 🍜 Ăn uống              [Expand ▾]  │
│   └─ ☕ Cà phê                      │
│   └─ 🍱 Ăn trưa                     │
│   └─ 🍳 Ăn sáng                     │
│   └─ 🍺 Bar/Pub                     │
│                                     │
│ 🚗 Di chuyển            [Expand ▾]  │
│   └─ ⛽ Xăng xe                     │
│   └─ 🅿️ Đỗ xe                      │
│   └─ 🚕 Taxi/Grab                   │
└─────────────────────────────────────┘
```

**Bulk Actions:**

Khi user select multiple categories (checkbox):

```
[3 selected] [Xóa] [Export] [Gộp vào...]
```

**Import/Export Categories:**

```
┌─────────────────────────────────────┐
│ Import/Export Danh mục              │
├─────────────────────────────────────┤
│ Import từ file:                     │
│ [Chọn file JSON/CSV]  [Import]      │
│                                     │
│ Export sang file:                   │
│ [ ] Chi tiêu only                   │
│ [ ] Thu nhập only                   │
│ [ ] Include default categories      │
│                                     │
│ [Xuất JSON] [Xuất CSV]              │
└─────────────────────────────────────┘
```

**Import Format (JSON):**

```json
{
    "categories": [
        {
            "name": "Ăn uống",
            "type": "expense",
            "icon": "🍜",
            "color": "#EF4444",
            "parent_id": null,
            "description": "Chi phí ăn uống hàng ngày"
        },
        {
            "name": "Cà phê",
            "type": "expense",
            "icon": "☕",
            "color": "#F59E0B",
            "parent_name": "Ăn uống"
        }
    ]
}
```

**Import Format (CSV):**

```csv
name,type,icon,color,parent_name,description
Ăn uống,expense,🍜,#EF4444,,Chi phí ăn uống hàng ngày
Cà phê,expense,☕,#F59E0B,Ăn uống,
```

**Reset to Default:**

```
┌─────────────────────────────────────┐
│ ⚠️ Khôi phục danh mục mặc định      │
├─────────────────────────────────────┤
│ Hành động này sẽ:                   │
│ • Xóa tất cả custom categories      │
│ • Khôi phục default categories      │
│ • Không ảnh hưởng transactions      │
│                                     │
│ Lưu ý: Transactions sẽ không có     │
│ category sau khi reset (cần chọn    │
│ lại category cho từng transaction)  │
│                                     │
│ Bạn có chắc chắn?                   │
│                                     │
│         [Hủy]   [Xác nhận]          │
└─────────────────────────────────────┘
```

**Search & Filter:**

```
[🔍 Tìm danh mục...]  [Type: All ▾]  [Sort: Name ▾]
```

Filters:

- Type: All | Chi tiêu | Thu nhập
- Sort: Name (A-Z) | Most used | Recent

**Statistics Card:**

Hiển thị ở top của page:

```
┌─────────────────────────────────────────────────────┐
│ 📊 Thống kê danh mục                                │
├─────────────────────────────────────────────────────┤
│ Tổng số: 18 danh mục                                │
│ • 12 danh mục chi tiêu                              │
│ • 4 danh mục thu nhập                               │
│ • 2 danh mục cả hai                                 │
│                                                     │
│ Top 3 danh mục được sử dụng nhiều nhất:            │
│ 1. 🍜 Ăn uống - 342 giao dịch                       │
│ 2. 🚗 Di chuyển - 156 giao dịch                     │
│ 3. 🛍️ Mua sắm - 98 giao dịch                        │
└─────────────────────────────────────────────────────┘
```

**Business Logic:**

**Validate Category:**

```typescript
function validateCategory(input: CategoryInput): ValidationResult {
    const errors: string[] = [];

    // Name validation
    if (!input.name || input.name.trim().length === 0) {
        errors.push("Tên danh mục không được để trống");
    }

    if (input.name.length > 100) {
        errors.push("Tên danh mục không được quá 100 ký tự");
    }

    // Type validation
    if (!["income", "expense", "both"].includes(input.type)) {
        errors.push("Loại danh mục không hợp lệ");
    }

    // Icon validation
    if (!input.icon || input.icon.trim().length === 0) {
        errors.push("Vui lòng chọn icon");
    }

    // Color validation
    if (!input.color || !/^#[0-9A-F]{6}$/i.test(input.color)) {
        errors.push("Màu sắc không hợp lệ");
    }

    return {
        isValid: errors.length === 0,
        errors,
    };
}
```

**Import Categories:**

```typescript
async function importCategories(
    userId: string,
    file: File,
    options: ImportOptions,
): Promise<ImportResult> {
    // Parse file (JSON or CSV)
    const data = await parseFile(file);

    const imported: Category[] = [];
    const errors: ImportError[] = [];

    for (const categoryData of data.categories) {
        try {
            // Validate
            const validation = validateCategory(categoryData);
            if (!validation.isValid) {
                errors.push({
                    name: categoryData.name,
                    error: validation.errors.join(", "),
                });
                continue;
            }

            // Check duplicate
            const existing = await CategoryRepository.findByName(
                userId,
                categoryData.name,
            );

            if (existing) {
                if (options.skipDuplicates) {
                    continue;
                } else if (options.overwriteDuplicates) {
                    await CategoryRepository.update(existing.id, categoryData);
                    imported.push(existing);
                    continue;
                }
            }

            // Create category
            const category = await CategoryRepository.create({
                ...categoryData,
                userId,
                isDefault: false,
            });

            imported.push(category);
        } catch (error) {
            errors.push({
                name: categoryData.name,
                error: error.message,
            });
        }
    }

    return {
        totalProcessed: data.categories.length,
        imported: imported.length,
        errors: errors.length,
        details: { imported, errors },
    };
}
```

**Merge Categories:**

```typescript
async function mergeCategories(
    userId: string,
    sourceCategoryIds: string[],
    targetCategoryId: string,
): Promise<MergeResult> {
    const target = await CategoryRepository.getById(targetCategoryId, userId);

    if (!target) {
        throw new Error("Target category not found");
    }

    let totalTransactionsMoved = 0;

    for (const sourceId of sourceCategoryIds) {
        // Update all transactions from source to target
        const count = await TransactionRepository.updateCategory(
            userId,
            sourceId,
            targetCategoryId,
        );

        totalTransactionsMoved += count;

        // Update budgets
        await BudgetRepository.updateCategory(
            userId,
            sourceId,
            targetCategoryId,
        );

        // Delete source category
        await CategoryRepository.delete(sourceId, userId);
    }

    return {
        targetCategory: target,
        categoriesDeleted: sourceCategoryIds.length,
        transactionsMoved: totalTransactionsMoved,
    };
}
```

**Reset to Default Categories:**

```typescript
async function resetToDefaultCategories(userId: string): Promise<ResetResult> {
    // Get all custom categories
    const customCategories = await CategoryRepository.list({
        userId,
        isDefault: false,
    });

    // Check if any custom category is in use
    const categoriesInUse: string[] = [];

    for (const category of customCategories) {
        const count = await TransactionRepository.countByCategory(
            userId,
            category.id,
        );

        if (count > 0) {
            categoriesInUse.push(category.name);
        }
    }

    if (categoriesInUse.length > 0) {
        throw new Error(
            `Không thể reset: các danh mục sau đang được sử dụng: ${categoriesInUse.join(", ")}`,
        );
    }

    // Delete all custom categories
    for (const category of customCategories) {
        await CategoryRepository.delete(category.id, userId);
    }

    // Ensure default categories exist
    await seedDefaultCategories(userId);

    return {
        deletedCount: customCategories.length,
        message: "Đã khôi phục về danh mục mặc định",
    };
}
```

**Acceptance Criteria:**

- ✅ User có thể tạo unlimited custom categories
- ✅ Emoji picker hiển thị đầy đủ và search được
- ✅ Color picker có preset colors và custom color input
- ✅ Sub-categories (hierarchy) hoạt động
- ✅ Không thể xóa default categories
- ✅ Không thể xóa categories đang có transactions
- ✅ Import từ JSON/CSV hoạt động
- ✅ Export sang JSON/CSV hoạt động
- ✅ Merge categories chuyển tất cả transactions
- ✅ Reset to default có confirmation dialog
- ✅ Search và filter hoạt động real-time
- ✅ Statistics card hiển thị đúng số liệu

### Business Logic

#### Update Profile

```typescript
async function updateProfile(
    userId: string,
    input: ProfileInput,
): Promise<User> {
    // Validation
    if (input.email && !isValidEmail(input.email)) {
        throw new Error("Invalid email format");
    }

    // Check email uniqueness if changed
    if (input.email) {
        const existing = await UserRepository.findByEmail(input.email);
        if (existing && existing.id !== userId) {
            throw new Error("Email already in use");
        }
    }

    // Upload avatar if provided
    let avatarUrl = undefined;
    if (input.avatar) {
        avatarUrl = await uploadAvatar(input.avatar);
    }

    // Update user
    const user = await UserRepository.update(userId, {
        fullName: input.fullName,
        email: input.email,
        phone: input.phone,
        avatarUrl,
        updatedAt: new Date(),
    });

    return user;
}
```

#### Change Password

```typescript
async function changePassword(
    userId: string,
    currentPassword: string,
    newPassword: string,
): Promise<void> {
    // Get user
    const user = await UserRepository.getById(userId);

    // Verify current password
    const isValid = await bcrypt.compare(currentPassword, user.passwordHash);
    if (!isValid) {
        throw new Error("Current password is incorrect");
    }

    // Validate new password
    if (newPassword.length < 8) {
        throw new Error("Password must be at least 8 characters");
    }

    if (
        !/[A-Z]/.test(newPassword) ||
        !/[a-z]/.test(newPassword) ||
        !/[0-9]/.test(newPassword)
    ) {
        throw new Error(
            "Password must include uppercase, lowercase, and number",
        );
    }

    // Hash new password
    const passwordHash = await bcrypt.hash(newPassword, 10);

    // Update user
    await UserRepository.update(userId, {
        passwordHash,
        updatedAt: new Date(),
    });
}
```

#### Export Data

```typescript
async function exportUserData(
    userId: string,
    format: "json" | "csv",
): Promise<Blob> {
    // Get all user data
    const user = await UserRepository.getById(userId);
    const accounts = await AccountRepository.list({ userId });
    const transactions = await TransactionRepository.list({ userId });
    const categories = await CategoryRepository.list({ userId });
    const budgets = await BudgetRepository.list({ userId });
    const alerts = await AlertRepository.list({ userId });

    const data = {
        user: { id: user.id, email: user.email, fullName: user.fullName },
        accounts,
        transactions,
        categories,
        budgets,
        alerts,
        exportedAt: new Date(),
    };

    if (format === "json") {
        return new Blob([JSON.stringify(data, null, 2)], {
            type: "application/json",
        });
    } else {
        // Convert to CSV (multiple files zipped)
        const zip = new JSZip();

        zip.file("accounts.csv", convertToCSV(accounts));
        zip.file("transactions.csv", convertToCSV(transactions));
        zip.file("categories.csv", convertToCSV(categories));
        zip.file("budgets.csv", convertToCSV(budgets));

        return await zip.generateAsync({ type: "blob" });
    }
}
```

#### Delete Account

```typescript
async function deleteAccount(userId: string, password: string): Promise<void> {
    // Verify password
    const user = await UserRepository.getById(userId);
    const isValid = await bcrypt.compare(password, user.passwordHash);

    if (!isValid) {
        throw new Error("Password is incorrect");
    }

    // Delete all user data (cascade)
    await TransactionRepository.deleteAll(userId);
    await BudgetRepository.deleteAll(userId);
    await AlertRepository.deleteAll(userId);
    await CategoryRepository.deleteAll(userId); // only custom categories
    await AccountRepository.deleteAll(userId);
    await ChatMessageRepository.deleteAll(userId);
    await UserRepository.delete(userId);

    // Optionally: store deleted user IDs for audit
    await DeletedUserRepository.create({
        userId,
        deletedAt: new Date(),
    });
}
```

### Acceptance Criteria

- ✅ Profile update works (name, email, avatar)
- ✅ Change password validates correctly
- ✅ Notification preferences save and apply
- ✅ Currency and date format changes reflect app-wide
- ✅ Theme switching (light/dark) works
- ✅ Export data generates valid JSON/CSV
- ✅ Delete account requires password confirmation
- ✅ Delete account removes all user data
- ✅ Settings page responsive on mobile

---

## Implementation Priorities

### Phase 1: MVP (Core Features)

1. ✅ Authentication (Login/Register)
2. ✅ Accounts Management
3. ✅ Transactions (Add/Edit/Delete)
4. ✅ Categories (Default + Custom)
5. ✅ Dashboard (Basic overview)

### Phase 2: Advanced Features

1. ✅ Budgets
2. ✅ Reports (Overview, Category, Account)
3. ✅ Alerts (Budget warnings)

### Phase 3: Intelligence

1. ✅ AI Chat Assistant
2. ✅ Forecasting
3. ✅ Insights & Patterns

### Phase 4: Polish

1. ✅ Settings & Preferences
2. ✅ Export/Import
3. ✅ Mobile optimization
4. ✅ Performance optimization

---

**Document Version**: 1.0  
**Last Updated**: February 27, 2026  
**Author**: Finance Hub Development Team
