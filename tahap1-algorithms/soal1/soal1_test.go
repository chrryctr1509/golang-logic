package soal1

import "testing"

func TestFindMatchingStrings(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		strings  []string
		want     []int
		wantBool bool
	}{
		{
			name:     "Basic duplicate - acbd",
			n:        4,
			strings:  []string{"abcd", "acbd", "aaab", "acbd"},
			want:     []int{2, 4},
			wantBool: true,
		},
		{
			name:     "Multiple sets - return first",
			n:        11,
			strings:  []string{"Satu", "Sate", "Tujuh", "Tusuk", "Tujuh", "Sate", "Bonus", "Tiga", "Puluh", "Tujuh", "Tusuk"},
			want:     []int{3, 5, 10},
			wantBool: true,
		},
		{
			name:     "No duplicates",
			n:        5,
			strings:  []string{"pisang", "goreng", "enak", "sekali", "rasanya"},
			want:     nil,
			wantBool: false,
		},
		{
			name:     "Case insensitive match",
			n:        3,
			strings:  []string{"Semua", "semua", "SEMUA"},
			want:     []int{1, 2, 3},
			wantBool: true,
		},
		{
			name:     "N mismatch",
			n:        3,
			strings:  []string{"a", "b"},
			want:     nil,
			wantBool: false,
		},
		{
			name:     "Single element",
			n:        1,
			strings:  []string{"alone"},
			want:     nil,
			wantBool: false,
		},
		{
			name:     "Empty strings",
			n:        2,
			strings:  []string{"", ""},
			want:     []int{1, 2},
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBool := FindMatchingStrings(tt.n, tt.strings)
			if gotBool != tt.wantBool {
				t.Errorf("FindMatchingStrings() bool = %v, want %v", gotBool, tt.wantBool)
				return
			}
			if gotBool {
				if len(got) != len(tt.want) {
					t.Errorf("FindMatchingStrings() indices = %v, want %v", got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("FindMatchingStrings() indices = %v, want %v", got, tt.want)
						return
					}
				}
			}
		})
	}
}
