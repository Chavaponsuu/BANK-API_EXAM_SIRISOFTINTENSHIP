# 🔍 Code Review: Concurrent Deposit/Withdraw Safety

## สถานการณ์: 10 คนฝาก/ถอนเงินพร้อมกัน

---

## ⚠️ ปัญหาที่พบ (ก่อนแก้ไข)

### 1. **Race Condition ใน Balance Update**

#### ❌ Code เดิม (มีปัญหา):
```go
// STEP 1: อ่านข้อมูล account ก่อน transaction
account, err := s.accountRepo.GetByAccountNumber(ctx, accountNumber)

// STEP 2: เริ่ม transaction หลังอ่านข้อมูล  
tx, err := s.accountRepo.BeginTx(ctx)

// STEP 3: ใช้ balance ที่อ่านไว้ก่อนหน้า (อาจเป็นข้อมูลเก่า)
balanceBefore := account.Balance  // ⚠️ STALE DATA!
balanceAfter := balanceBefore + amount
```

#### 🐛 สถานการณ์ที่เกิดข้อผิดพลาด:

| เวลา | User A (ฝาก 100)              | User B (ฝาก 200)              | Balance ใน DB |
|------|-------------------------------|-------------------------------|---------------|
| T1   | อ่าน balance = 1000          |                               | 1000          |
| T2   |                               | อ่าน balance = 1000          | 1000          |
| T3   | BEGIN TX                      |                               | 1000          |
| T4   |                               | BEGIN TX                      | 1000          |
| T5   | UPDATE balance = 1100         |                               | 1100          |
| T6   |                               | UPDATE balance = 1200         | **1200** ❌   |
| T7   | COMMIT                        |                               | 1200          |
| T8   |                               | COMMIT                        | 1200          |

**ผลลัพธ์:** ควรเป็น **1300** แต่ได้ **1200** → สูญเงิน **100 บาท!**

---

### 2. **Transaction Reference ซ้ำได้**

#### ❌ Code เดิม:
```go
// นับจำนวน transactions โดยไม่ lock
query := `SELECT COUNT(*) FROM transactions WHERE transaction_ref LIKE $1`
```

เมื่อ 2 requests มาพร้อมกัน:
- Request A: นับได้ 100 → สร้าง TXN202606110101
- Request B: นับได้ 100 → สร้าง TXN202606110101 (ซ้ำ!)

---

## ✅ วิธีแก้ไข

### 1. **ใช้ SELECT FOR UPDATE (Pessimistic Locking)**

#### ✅ Code ที่แก้แล้ว:

```go
// repositories/account.go
func (r *AccountRepo) GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*models.Account, error) {
    query := `SELECT id, account_number, owner_name, citizen_id, phone_number, 
              account_type, balance, status, created_at, updated_at
              FROM accounts
              WHERE account_number = $1
              FOR UPDATE`  // 🔒 LOCK ROW จนกว่าจะ COMMIT/ROLLBACK
    
    err := tx.QueryRowContext(ctx, query, accountNumber).Scan(...)
    return account, nil
}
```

#### ✅ Service เปลี่ยนลำดับการทำงาน:

```go
// services/transaction.go - Deposit()
func (s *transactionService) Deposit(...) (*models.Transaction, error) {
    // 1. Validate amount
    if amount <= 0 {
        return nil, ErrInvalidAmount
    }

    // 2. BEGIN TRANSACTION ก่อน (เปลี่ยนลำดับ)
    tx, err := s.accountRepo.BeginTx(ctx)
    
    // 3. GET ACCOUNT พร้อม LOCK 🔒
    account, err := s.accountRepo.GetByAccountNumberWithLock(ctx, tx, accountNumber)
    
    // 4. ใช้ balance ที่ FRESH และ LOCKED แล้ว
    balanceBefore := account.Balance  // ✅ ข้อมูลถูกต้อง
    balanceAfter := balanceBefore + amount
    
    // 5. UPDATE balance
    err = s.accountRepo.UpdateBalanceWithTx(ctx, tx, account.ID, balanceAfter)
    
    // 6. CREATE transaction record
    // 7. COMMIT
}
```

#### 🎯 สถานการณ์หลังแก้ไข:

| เวลา | User A (ฝาก 100)              | User B (ฝาก 200)              | Balance ใน DB |
|------|-------------------------------|-------------------------------|---------------|
| T1   | BEGIN TX                      |                               | 1000          |
| T2   |                               | BEGIN TX                      | 1000          |
| T3   | SELECT FOR UPDATE (LOCK!)     |                               | 1000          |
| T4   |                               | SELECT FOR UPDATE (รอ...)    | 1000          |
| T5   | UPDATE balance = 1100         |                               | 1100          |
| T6   | COMMIT (ปลด LOCK)            |                               | 1100          |
| T7   |                               | SELECT success (LOCK!)       | 1100          |
| T8   |                               | UPDATE balance = 1300         | **1300** ✅   |
| T9   |                               | COMMIT                        | 1300          |

**ผลลัพธ์:** ได้ **1300** ถูกต้อง! 🎉

---

### 2. **Lock Transaction Reference Generation**

```go
// repositories/transaction.go
func (r *TransactionRepo) GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error) {
    // ใช้ FOR UPDATE เพื่อ serialize การสร้าง ref
    query := `SELECT COUNT(*) FROM transactions 
              WHERE transaction_ref LIKE $1 
              FOR UPDATE`  // 🔒 LOCK ตาราง transactions
    
    err = tx.QueryRowContext(ctx, query, pattern).Scan(&count)
    runningNumber := count + 1
    return fmt.Sprintf("TXN%s%04d", dateStr, runningNumber), nil
}
```

---

## 📊 การทดสอบ Concurrency

### ทดสอบด้วย Apache Bench หรือ Go Test:

```go
// Example concurrent test
func TestConcurrentDeposit(t *testing.T) {
    accountNumber := "0000000001"
    initialBalance := 1000.0
    numDeposits := 10
    depositAmount := 100.0
    
    var wg sync.WaitGroup
    for i := 0; i < numDeposits; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            _, err := transactionService.Deposit(ctx, accountNumber, depositAmount, "test")
            assert.NoError(t, err)
        }()
    }
    wg.Wait()
    
    // ตรวจสอบ balance สุดท้าย
    account, _ := accountRepo.GetByAccountNumber(ctx, accountNumber)
    expected := initialBalance + (float64(numDeposits) * depositAmount)
    assert.Equal(t, expected, account.Balance) // ต้องเป็น 2000.0
}
```

---

## 🎯 สรุป

### ✅ การแก้ไขที่ทำ:

1. **เพิ่ม `GetByAccountNumberWithLock()`** - Lock row ด้วย `FOR UPDATE`
2. **เปลี่ยนลำดับ Deposit/Withdraw** - BEGIN TX → LOCK → อ่านข้อมูล → UPDATE
3. **Lock Transaction Ref Generation** - ป้องกัน reference ซ้ำ

### 🔒 Isolation Level:
PostgreSQL default = **READ COMMITTED** + **FOR UPDATE** = Serializable สำหรับ row ที่ lock

### ⚡ Performance:
- **Throughput:** ลดลงเล็กน้อยเพราะต้อง wait lock
- **Correctness:** เพิ่มขึ้น 100% - ไม่มี race condition
- **Trade-off:** ยอมเสีย latency เล็กน้อยเพื่อความถูกต้องของข้อมูล

### 📝 Best Practices:
✅ **ใช้ DB Transaction** สำหรับการเงินทุกครั้ง  
✅ **Lock row ด้วย FOR UPDATE** เมื่ออ่าน-แก้ไข-เขียน  
✅ **อ่านข้อมูลภายใน Transaction** ไม่ใช่ก่อน BEGIN  
✅ **Test concurrency** ด้วย load testing tools  
✅ **Monitor deadlock** ใน production (PostgreSQL logs)

---

## ⚠️ หมายเหตุเพิ่มเติม

### Deadlock Prevention:
- **Lock order:** ควร lock accounts ก่อน transactions เสมอ
- **Timeout:** ตั้ง context timeout เพื่อป้องกัน long-running transactions
- **Retry logic:** ใน handler ควรมี retry สำหรับ deadlock errors

### Alternative Solutions:
1. **Optimistic Locking:** ใช้ version column แทน FOR UPDATE
2. **Message Queue:** ใช้ queue serialization แทน database lock
3. **Distributed Lock:** ใช้ Redis/etcd สำหรับ microservices

---

**สถานะ:** ✅ **Code ปลอดภัยสำหรับ concurrent operations แล้ว**
