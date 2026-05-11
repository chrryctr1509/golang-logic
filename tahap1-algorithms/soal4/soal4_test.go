package soal4

import (
	"testing"
	"time"
)

func TestCheckLeave(t *testing.T) {
	tests := []struct {
		name              string
		jumlahCutiBersama int
		tanggalJoin       time.Time
		tanggalCuti       time.Time
		durasiCuti        int
		wantAllowed       bool
		wantReason        string
	}{
		{
			name:              "Belum 180 hari",
			jumlahCutiBersama: 7,
			tanggalJoin:       time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
			tanggalCuti:       time.Date(2021, 7, 5, 0, 0, 0, 0, time.UTC),
			durasiCuti:        1,
			wantAllowed:       false,
			wantReason:        "Karena belum 180 hari sejak tanggal join karyawan",
		},
		{
			name:              "Kuota exceed (only 1 day allowed)",
			jumlahCutiBersama: 7,
			tanggalJoin:       time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
			tanggalCuti:       time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC),
			durasiCuti:        3,
			wantAllowed:       false,
			wantReason:        "Karena hanya boleh mengambil 1 hari cuti",
		},
		{
			name:              "Approved - 1 day, old employee",
			jumlahCutiBersama: 7,
			tanggalJoin:       time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
			tanggalCuti:       time.Date(2021, 12, 18, 0, 0, 0, 0, time.UTC),
			durasiCuti:        1,
			wantAllowed:       true,
			wantReason:        "Approved",
		},
		{
			name:              "Approved - 3 days, old employee",
			jumlahCutiBersama: 7,
			tanggalJoin:       time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
			tanggalCuti:       time.Date(2021, 12, 18, 0, 0, 0, 0, time.UTC),
			durasiCuti:        3,
			wantAllowed:       true,
			wantReason:        "Approved",
		},
		{
			name:              "Exceed 3 days",
			jumlahCutiBersama: 0,
			tanggalJoin:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			tanggalCuti:       time.Date(2021, 12, 18, 0, 0, 0, 0, time.UTC),
			durasiCuti:        4,
			wantAllowed:       false,
			wantReason:        "Karena cuti pribadi maksimal 3 hari berturutan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LeaveRequest{
				JumlahCutiBersama: tt.jumlahCutiBersama,
				TanggalJoin:       tt.tanggalJoin,
				TanggalCuti:       tt.tanggalCuti,
				DurasiCuti:        tt.durasiCuti,
			}
			got := CheckLeave(req)
			if got.Allowed != tt.wantAllowed {
				t.Errorf("CheckLeave() Allowed = %v, want %v", got.Allowed, tt.wantAllowed)
			}
			if got.Reason != tt.wantReason {
				t.Errorf("CheckLeave() Reason = %q, want %q", got.Reason, tt.wantReason)
			}
		})
	}
}
