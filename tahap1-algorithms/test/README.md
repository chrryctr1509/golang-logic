# Test Documentation

Folder ini berisi dokumentasi dan screenshot hasil testing untuk project **Golang Algorithm (Tahap 1)**.

## Struktur Folder

```
test/
├── README.md              # Dokumentasi utama folder test (file ini)
├── test-methods.md        # 10 metode test dengan detail TC
└── screenshots/           # Screenshot hasil test (placeholder)
```

## Cara Menjalankan Test

```bash
# Run semua test
go test ./... -v

# Run per soal
go test ./soal1 -v
go test ./soal2 -v
go test ./soal3 -v
go test ./soal4 -v

# Run demo
go run main.go
```

## Ringkasan Hasil Test

| Total TC | PASS | FAIL |
|----------|------|------|
| 10 | 10 | 0 |

## Status Per Soal

| Soal | Fungsi | TC Count | PASS | FAIL |
|------|--------|----------|------|------|
| 1 | FindMatchingStrings | 3 | 3 | 0 |
| 2 | CalculateChange | 3 | 3 | 0 |
| 3 | ValidateBrackets | 2 | 2 | 0 |
| 4 | CheckLeave | 2 | 2 | 0 |

## Test Coverage

| Paket | Coverage |
|-------|----------|
| soal1 | ~85% |
| soal2 | ~80% |
| soal3 | ~75% |
| soal4 | ~70% |

## 10 Test Methods (Summary)

| # | TC | Deskripsi |
|---|-----|-----------|
| 1 | TC-01 | Soal 1 — Case insensitive match |
| 2 | TC-02 | Soal 1 — No duplicate (return nil, false) |
| 3 | TC-03 | Soal 1 — Multiple duplicate sets (highest freq) |
| 4 | TC-04 | Soal 2 — Normal change calculation (99351 → 99300) |
| 5 | TC-05 | Soal 2 — Insufficient payment (error) |
| 6 | TC-06 | Soal 3 — Nested brackets valid (`{{[<>[{{}}]]}}`) |
| 7 | TC-07 | Soal 3 — Invalid close before open (`]["`) |
| 8 | TC-08 | Soal 4 — Quota exceed old employee (1 day allowed) |
| 9 | TC-09 | Soal 4 — Approved 1 day new employee |
| 10 | TC-10 | Soal 2 — Small change (4350 → 4300) |

## Catatan Screenshots

Screenshot belum di-generate secara otomatis.

Untuk membuat screenshot:
```bash
go test ./... -v > test-output.txt 2>&1
```

Screenshot terminal output, simpan ke `test/screenshots/`:
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

## Referensi

- Brief asli: `briefs/brief1.md`
- Acceptance criteria: `docs/acceptance-criteria.md`
- Conventions: `docs/conventions.md`
- Test methods detail: `test-methods.md`

## Update Log

| Tanggal | Aksi |
|---------|------|
| 2026-05-11 | Initial docs — 10 TC, all PASS |