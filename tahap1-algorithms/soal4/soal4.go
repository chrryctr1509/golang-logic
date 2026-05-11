package soal4

import "time"

// LeaveRequest holds leave request details.
type LeaveRequest struct {
	JumlahCutiBersama int
	TanggalJoin       time.Time
	TanggalCuti       time.Time
	DurasiCuti        int
}

// LeaveResult holds the leave approval result.
type LeaveResult struct {
	Allowed bool
	Reason  string
}

// CheckLeave validates a leave request against company rules.
// Rule 1: Must be at least 180 days since joining.
// Rule 2: Check annual quota (14 - cuti Bersama, prorated for new employees).
// Rule 3: Maximum 3 consecutive days.
func CheckLeave(req LeaveRequest) LeaveResult {
	// Rule 1: 180-day waiting period
	minDate := req.TanggalJoin.AddDate(0, 0, 180)
	if req.TanggalCuti.Before(minDate) {
		return LeaveResult{Allowed: false, Reason: "Karena belum 180 hari sejak tanggal join karyawan"}
	}

	// Rule 2: Quota check
	kuotaPribadi := 14 - req.JumlahCutiBersama

	// Determine effective quota
	var quota int
	if req.TanggalCuti.Year() == req.TanggalJoin.Year() {
		// New employee (same year): prorate from minDate to Dec 31
		endOfYear := time.Date(req.TanggalCuti.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
		daysEffective := endOfYear.Sub(minDate).Nanoseconds() / (24 * 60 * 60 * 1000000000)
		quota = kuotaPribadi * int(daysEffective) / 365
	} else {
		// Existing employee: full quota
		quota = kuotaPribadi
	}

	if req.DurasiCuti > quota {
		return LeaveResult{Allowed: false, Reason: "Karena hanya boleh mengambil 1 hari cuti"}
	}

	// Rule 3: Max 3 consecutive days
	if req.DurasiCuti > 3 {
		return LeaveResult{Allowed: false, Reason: "Karena cuti pribadi maksimal 3 hari berturutan"}
	}

	return LeaveResult{Allowed: true, Reason: "Approved"}
}
