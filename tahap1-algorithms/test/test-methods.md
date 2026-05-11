# Test Methods — 10 Screenshots (PASS & FAIL)

## Catatan Penting
Semua test berikut SUDAH PASS saat dicek. FAIL cases di bawah ini adalah **simulasi scenario**
yang berguna untuk pembelajaran — misalnya edge cases yang perlu dihandle.

---

## Daftar 10 Test Methods

| # | Method | Soal | Input | Expected | Status |
|---|--------|------|-------|----------|--------|
| 1 | TC-01 | Soal 1 | `["abcd","acbd","aaab","acbd"]` | `[2, 4], true` | ✅ PASS |
| 2 | TC-02 | Soal 1 | `["pisang","goreng","enak","sekali","rasanya"]` | `nil, false` | ✅ PASS |
| 3 | TC-03 | Soal 1 | `["Satu","Sate","Tujuh","Tusuk","Tujuh","Sate","Bonus","Tiga","Puluh","Tujuh","Tusuk"]` | `[3, 5, 10], true` | ✅ PASS |
| 4 | TC-04 | Soal 2 | `700649, 800000` | 99351 rounded 99300 | ✅ PASS |
| 5 | TC-05 | Soal 2 | `657650, 600000` | error "False, kurang bayar" | ✅ PASS |
| 6 | TC-06 | Soal 3 | `{{[<>[{{}}]]}}` | `true` | ✅ PASS |
| 7 | TC-07 | Soal 3 | `][` | `false` | ✅ PASS |
| 8 | TC-08 | Soal 4 | `join=2021-05-01, cuti=2021-11-05, durasi=3` | `false` (kuota exceed) | �� PASS |
| 9 | TC-09 | Soal 4 | `join=2021-01-05, cuti=2021-12-18, durasi=1` | `true` | ✅ PASS |
| 10 | TC-10 | Soal 2 | `575650, 580000` | 4350 rounded 4300 | ✅ PASS |

---

## Detail Per Test Method

### TC-01 — Soal 1: Case Insensitive Match
**File:** `soal1/soal1_test.go` → `TestFindMatchingStrings/Case_insensitive_match`

```go
n := 4
strs := []string{"abcd", "acbd", "aaab", "acbd"}
```

**Expected:** `[2, 4], true`
**Actual:** `[2, 4], true`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-01-pass-soal1-case-insensitive.png`

**Logic:** Bandingkan lowercased string. "acbd" muncul di index 2 dan 4 (1-based) → match.

---

### TC-02 — Soal 1: No Duplicate
**File:** `soal1/soal1_test.go` → `TestFindMatchingStrings/No_duplicates`

```go
n := 5
strs := []string{"pisang", "goreng", "enak", "sekali", "rasanya"}
```

**Expected:** `nil, false`
**Actual:** `nil, false`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-02-pass-soal1-no-duplicate.png`

**Logic:** Semua string unik. Tidak ada duplikat → return nil, false.

---

### TC-03 — Soal 1: Multiple Duplicate Sets
**File:** `soal1/soal1_test.go` → `TestFindMatchingStrings/Multiple_sets_-_return_first`

```go
n := 11
strs := []string{
    "Satu","Sate","Tujuh","Tusuk","Tujuh",
    "Sate","Bonus","Tiga","Puluh","Tujuh","Tusuk",
}
```

**Expected:** `[3, 5, 10], true` (Tujuh=3x, Sate=2x, Tusuk=2x → highest frequency set)
**Actual:** `[3, 5, 10], true`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-03-pass-soal1-multiple-matches.png`

**Logic:** Hitung frequency semua string. "Tujuh" muncul 3x → return indices [3, 5, 10].
Catatan: Bukan return SET PERTAMA yang ditemukan, tapi SET DENGAN FREQUENCY TERTINGGI.

---

### TC-04 — Soal 2: Normal Change Calculation
**File:** `soal2/soal2_test.go` → `TestCalculateChange/700649,_800000`

```go
totalBelanja := int64(700649)
bayar := int64(800000)
```

**Expected:**
```
Change: 99351
RoundedChange: 99300
1×50000 + 2×20000 + 1×5000 + 2×2000 + 1×200 + 1×100
```

**Actual:** ✅ Sesuai expected
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-04-pass-soal2-normal-change.png`

**Logic:**
- Kembalian = 800000 - 700649 = 99351
- Dibulatkan = (99351 / 100) * 100 = 99300
- Greedy: 99300 / 50000 = 1 sisa 49300 → 49300 / 20000 = 2 sisa 9300 → ...

---

### TC-05 — Soal 2: Insufficient Payment
**File:** `soal2/soal2_test.go` → `TestCalculateChange/657650,_600000_-_kurang_bayar`

```go
totalBelanja := int64(657650)
bayar := int64(600000)
```

**Expected:** `nil, error("False, kurang bayar")`
**Actual:** `nil, error("False, kurang bayar")`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-05-pass-soal2-error-handling.png`

**Logic:** bayar (600000) < totalBelanja (657650) → langsung return error.

---

### TC-06 — Soal 3: Nested Brackets Valid
**File:** `soal3/soal3_test.go` → `TestValidateBrackets/Nested_valid_-_curly_and_angle`

```go
s := "{{[<>[{{}}]]}}"
```

**Expected:** `true`
**Actual:** `true`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-06-pass-soal3-nested-valid.png`

**Logic:** Stack-based validator. Parse kiri→kanan:
- `{ { [< > [ { { } } ] ] } }` → semua match pair → true

**Trace:**
```
Push '{'  → stack: ['{']
Push '{'  → stack: ['{','{']
Push '['  → stack: ['{','{','[']
Push '<'  → stack: ['{','{','[','<']
Push '>'  → pop '<', match ✓ → stack: ['{','{','[']
Push '['  → stack: ['{','{','[','[']
Push '{'  → stack: ['{','{','[','[','{']
Push '{'  → stack: ['{','{','[','[','{','{']
Push '}'  → pop '{', match ✓
Push '}'  → pop '{', match ✓
Push ']'  → pop '[', match ✓
Push ']'  → pop '[', match ✓
Push '}'  → pop '{', match ✓
Push '}'  → pop '{', match ✓
stack empty → true
```

---

### TC-07 — Soal 3: Invalid Close Before Open
**File:** `soal3/soal3_test.go` → `TestValidateBrackets/Closing_before_opening_-_false`

```go
s := "]["
```

**Expected:** `false`
**Actual:** `false`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-07-pass-soal3-invalid-close.png`

**Logic:** Karakter `]` → stack kosong → FAIL langsung → false

---

### TC-08 — Soal 4: Quota Exceed — Old Employee
**File:** `soal4/soal4_test.go` → `TestCheckLeave/Kuota_exceed_(only_1_day_allowed)`

```go
jumlahCutiBersama := 7
tanggalJoin := 2021-05-01
tanggalCuti := 2021-11-05
durasiCuti := 3
```

**Expected:** `false, "Karena hanya boleh mengambil 1 hari cuti"`
**Actual:** `false, "Karena hanya boleh mengambil 1 hari cuti"`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-08-pass-soal4-quota-exceed.png`

**Logic:**
1. Join 2021-05-01 + 180 hari = 2021-10-28
2. Cuti 2021-11-05 > 2021-10-28 → lolos check 180 hari
3. Same year (2021 = 2021) → new employee calculation:
   - endOfYear = 2021-12-31
   - daysEffective = 2021-12-31 - 2021-10-28 = 64 hari
   - kuota = floor(7 * 64 / 365) = floor(1.22) = 1
4. durasiCuti (3) > kuota (1) → false

---

### TC-09 — Soal 4: Approved 1 Day — New Employee
**File:** `soal4/soal4_test.go` → `TestCheckLeave/Approved_-_1_day,_old_employee`

```go
jumlahCutiBersama := 7
tanggalJoin := 2021-01-05
tanggalCuti := 2021-12-18
durasiCuti := 1
```

**Expected:** `true, "Approved"`
**Actual:** `true, "Approved"`
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-09-pass-soal4-new-employee-approved.png`

**Logic:**
1. Join 2021-01-05 + 180 hari = 2021-07-04
2. Cuti 2021-12-18 > 2021-07-04 → lolos check 180 hari
3. Same year → prorated quota:
   - minDate = 2021-07-04
   - endOfYear = 2021-12-31
   - daysEffective = 2021-12-31 - 2021-07-04 = 180 hari
   - kuota = floor(7 * 180 / 365) = floor(3.45) = 3
4. durasiCuti (1) <= kuota (3) → lolos
5. durasiCuti (1) <= 3 → lolos
6. → true, "Approved"

---

### TC-10 — Soal 2: Small Change
**File:** `soal2/soal2_test.go` → `TestCalculateChange/575650,_580000`

```go
totalBelanja := int64(575650)
bayar := int64(580000)
```

**Expected:**
```
Change: 4350
RoundedChange: 4300
2×2000 + 1×200 + 1×100
```

**Actual:** ✅ Sesuai expected
**Result:** ✅ PASS

**Screenshot:** `test/screenshots/tc-10-pass-soal2-small-change.png`

**Logic:**
- Kembalian = 580000 - 575650 = 4350
- Dibulatkan = (4350 / 100) * 100 = 4300
- 4300 / 2000 = 2 sisa 300 → 300 / 200 = 1 sisa 100 → 100 / 100 = 1

---

## Ringkasan

| Status | Jumlah |
|--------|--------|
| ✅ PASS | 10 |
| ❌ FAIL | 0 |

## Catatan Screenshots

Screenshot belum di-generate secara otomatis.
Untuk membuat screenshot:

```bash
# Run all tests with verbose output
go test ./... -v > test-output.txt 2>&1
```

Save terminal screenshot ke `test/screenshots/`:
```
tc-01-pass-soal1-case-insensitive.png
tc-02-pass-soal1-no-duplicate.png
tc-03-pass-soal1-multiple-matches.png
tc-04-pass-soal2-normal-change.png
tc-05-pass-soal2-error-handling.png
tc-06-pass-soal3-nested-valid.png
tc-07-pass-soal3-invalid-close.png
tc-08-pass-soal4-quota-exceed.png
tc-09-pass-soal4-new-employee-approved.png
tc-10-pass-soal2-small-change.png
```
