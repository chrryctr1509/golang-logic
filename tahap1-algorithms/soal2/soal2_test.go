package soal2

import (
	"testing"
)

func TestCalculateChange(t *testing.T) {
	tests := []struct {
		name          string
		totalBelanja  int64
		bayar         int64
		wantChange    int64
		wantRounded   int64
		wantDenoms    map[int64]int
		wantErr       bool
		errContains   string
	}{
		{
			name:         "700649, 800000",
			totalBelanja: 700649,
			bayar:        800000,
			wantChange:   99351,
			wantRounded:  99300,
			wantDenoms: map[int64]int{
				50000: 1,
				20000: 2,
				5000:  1,
				2000:  2,
				200:   1,
				100:   1,
			},
			wantErr: false,
		},
		{
			name:         "575650, 580000",
			totalBelanja: 575650,
			bayar:        580000,
			wantChange:   4350,
			wantRounded:  4300,
			wantDenoms: map[int64]int{
				2000: 2,
				200:  1,
				100:  1,
			},
			wantErr: false,
		},
		{
			name:        "657650, 600000 - kurang bayar",
			totalBelanja: 657650,
			bayar:        600000,
			wantErr:      true,
			errContains:  "False, kurang bayar",
		},
		{
			name:         "Exact payment",
			totalBelanja: 100000,
			bayar:        100000,
			wantChange:   0,
			wantRounded:  0,
			wantDenoms:   map[int64]int{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateChange(tt.totalBelanja, tt.bayar)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CalculateChange() expected error containing %q, got nil", tt.errContains)
				} else if err.Error() != tt.errContains {
					t.Errorf("CalculateChange() error = %q, want %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("CalculateChange() unexpected error: %v", err)
				return
			}
			if got.Change != tt.wantChange {
				t.Errorf("Change = %d, want %d", got.Change, tt.wantChange)
			}
			if got.RoundedChange != tt.wantRounded {
				t.Errorf("RoundedChange = %d, want %d", got.RoundedChange, tt.wantRounded)
			}
			for d, cnt := range tt.wantDenoms {
				if got.Denominations[d] != cnt {
					t.Errorf("Denominations[%d] = %d, want %d", d, got.Denominations[d], cnt)
				}
			}
		})
	}
}
