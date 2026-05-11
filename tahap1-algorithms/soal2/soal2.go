package soal2

import "errors"

// ChangeResult holds the change calculation result.
type ChangeResult struct {
	Change        int64
	RoundedChange int64
	Denominations map[int64]int
}

// CalculateChange computes the change using greedy denominations.
// Returns error if payment is less than total.
func CalculateChange(totalBelanja int64, bayar int64) (*ChangeResult, error) {
	if bayar < totalBelanja {
		return nil, errors.New("False, kurang bayar")
	}

	change := bayar - totalBelanja
	rounded := (change / 100) * 100

	denoms := []int64{100000, 50000, 20000, 10000, 5000, 2000, 1000, 500, 200, 100}
	result := make(map[int64]int)

	for _, d := range denoms {
		if rounded >= d {
			result[d] = int(rounded / d)
			rounded %= d
		}
	}

	return &ChangeResult{
		Change:        change,
		RoundedChange: (change / 100) * 100,
		Denominations: result,
	}, nil
}
