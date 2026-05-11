package main

import (
	"fmt"
	"time"

	"tahap1-algorithms/soal1"
	"tahap1-algorithms/soal2"
	"tahap1-algorithms/soal3"
	"tahap1-algorithms/soal4"
)

func main() {
	fmt.Println("=== SOAL 1: String Matching ===")
	indices, ok := soal1.FindMatchingStrings(4, []string{"abcd", "acbd", "aaab", "acbd"})
	fmt.Printf("Input: 4, [\"abcd\",\"acbd\",\"aaab\",\"acbd\"]\n")
	fmt.Printf("Output: %v, %v\n\n", indices, ok)

	indices2, ok2 := soal1.FindMatchingStrings(11, []string{"Satu", "Sate", "Tujuh", "Tusuk", "Tujuh", "Sate", "Bonus", "Tiga", "Puluh", "Tujuh", "Tusuk"})
	fmt.Printf("Input: 11, [\"Satu\",\"Sate\",\"Tujuh\",\"Tusuk\",\"Tujuh\",\"Sate\",\"Bonus\",\"Tiga\",\"Puluh\",\"Tujuh\",\"Tusuk\"]\n")
	fmt.Printf("Output: %v, %v\n\n", indices2, ok2)

	indices3, ok3 := soal1.FindMatchingStrings(5, []string{"pisang", "goreng", "enak", "sekali", "rasanya"})
	fmt.Printf("Input: 5, [\"pisang\",\"goreng\",\"enak\",\"sekali\",\"rasanya\"]\n")
	fmt.Printf("Output: %v, %v\n\n", indices3, ok3)

	fmt.Println("=== SOAL 2: Kasir Kembalian ===")

	res1, _ := soal2.CalculateChange(700649, 800000)
	fmt.Printf("Input: 700649, 800000\n")
	fmt.Printf("Kembalian: %d, Dibulatkan: %d\n", res1.Change, res1.RoundedChange)
	fmt.Print("Pecahan:\n")
	for d := range res1.Denominations {
		if d >= 1000 {
			fmt.Printf("  %d lembar %d\n", res1.Denominations[d], d)
		} else {
			fmt.Printf("  %d koin %d\n", res1.Denominations[d], d)
		}
	}
	fmt.Println()

	res2, _ := soal2.CalculateChange(575650, 580000)
	fmt.Printf("Input: 575650, 580000\n")
	fmt.Printf("Kembalian: %d, Dibulatkan: %d\n", res2.Change, res2.RoundedChange)
	fmt.Print("Pecahan:\n")
	for d := range res2.Denominations {
		if d >= 1000 {
			fmt.Printf("  %d lembar %d\n", res2.Denominations[d], d)
		} else {
			fmt.Printf("  %d koin %d\n", res2.Denominations[d], d)
		}
	}
	fmt.Println()

	_, err := soal2.CalculateChange(657650, 600000)
	fmt.Printf("Input: 657650, 600000\n")
	fmt.Printf("Output: %v\n\n", err)

	fmt.Println("=== SOAL 3: Validasi Bracket ===")
	fmt.Printf("Input: \"{{[<>[{{}}]]}}\" → %v\n", soal3.ValidateBrackets("{{[<>[{{}}]]}}"))
	fmt.Printf("Input: \"[{}<>]\" → %v\n", soal3.ValidateBrackets("[{}<>]"))
	fmt.Printf("Input: \"]\" → %v\n", soal3.ValidateBrackets("]"))
	fmt.Printf("Input: \"][\" → %v\n", soal3.ValidateBrackets("]["))
	fmt.Printf("Input: \"[>]\" → %v\n\n", soal3.ValidateBrackets("[>]"))

	fmt.Println("=== SOAL 4: Cuti Karyawan ===")

	res4a := soal4.CheckLeave(soal4.LeaveRequest{
		JumlahCutiBersama: 7,
		TanggalJoin:       time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
		TanggalCuti:       time.Date(2021, 7, 5, 0, 0, 0, 0, time.UTC),
		DurasiCuti:        1,
	})
	fmt.Printf("Belum 180 hari: Allowed=%v, Reason=%q\n", res4a.Allowed, res4a.Reason)

	res4b := soal4.CheckLeave(soal4.LeaveRequest{
		JumlahCutiBersama: 7,
		TanggalJoin:       time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
		TanggalCuti:       time.Date(2021, 11, 5, 0, 0, 0, 0, time.UTC),
		DurasiCuti:        3,
	})
	fmt.Printf("Kuota exceed: Allowed=%v, Reason=%q\n", res4b.Allowed, res4b.Reason)

	res4c := soal4.CheckLeave(soal4.LeaveRequest{
		JumlahCutiBersama: 7,
		TanggalJoin:       time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
		TanggalCuti:       time.Date(2021, 12, 18, 0, 0, 0, 0, time.UTC),
		DurasiCuti:        1,
	})
	fmt.Printf("Approved 1 hari: Allowed=%v, Reason=%q\n", res4c.Allowed, res4c.Reason)

	res4d := soal4.CheckLeave(soal4.LeaveRequest{
		JumlahCutiBersama: 7,
		TanggalJoin:       time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
		TanggalCuti:       time.Date(2021, 12, 18, 0, 0, 0, 0, time.UTC),
		DurasiCuti:        3,
	})
	fmt.Printf("Approved 3 hari: Allowed=%v, Reason=%q\n", res4d.Allowed, res4d.Reason)
}
