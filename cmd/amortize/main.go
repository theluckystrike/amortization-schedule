// Command amortize prints a fixed-rate amortization schedule as CSV, or a
// one-line summary, using the amortizationschedule package.
//
//	amortize [-summary] PRINCIPAL ANNUAL_RATE_PERCENT MONTHS
//
// Example:
//
//	amortize 85000 8.25 60           # 60 CSV rows: period,payment,interest,principal,balance
//	amortize -summary 85000 8.25 60  # payment, total interest, total paid, effective annual rate
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	amortizationschedule "github.com/theluckystrike/amortization-schedule"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: amortize [-summary] PRINCIPAL ANNUAL_RATE_PERCENT MONTHS")
	fmt.Fprintln(os.Stderr, "example: amortize 85000 8.25 60")
}

func main() {
	summary := flag.Bool("summary", false, "print only payment, total interest, total paid and effective annual rate")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 3 {
		usage()
		os.Exit(2)
	}
	principal, err1 := strconv.ParseFloat(flag.Arg(0), 64)
	ratePct, err2 := strconv.ParseFloat(flag.Arg(1), 64)
	months, err3 := strconv.ParseUint(flag.Arg(2), 10, 32)
	if err1 != nil || err2 != nil || err3 != nil || principal <= 0 || months == 0 || ratePct < 0 {
		usage()
		os.Exit(2)
	}
	rate := ratePct / 100
	rows := amortizationschedule.Schedule(principal, rate, uint32(months))
	if *summary {
		fmt.Printf("payment=%.2f total_interest=%.2f total_paid=%.2f effective_annual_rate=%.4f%%\n",
			amortizationschedule.MonthlyPayment(principal, rate, uint32(months)),
			amortizationschedule.TotalInterest(rows),
			amortizationschedule.TotalPaid(rows),
			amortizationschedule.EffectiveAPR(rate)*100)
		return
	}
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"period", "payment", "interest", "principal", "balance"})
	for _, r := range rows {
		_ = w.Write([]string{
			strconv.Itoa(r.Period),
			strconv.FormatFloat(r.Payment, 'f', 2, 64),
			strconv.FormatFloat(r.Interest, 'f', 2, 64),
			strconv.FormatFloat(r.Principal, 'f', 2, 64),
			strconv.FormatFloat(r.Balance, 'f', 2, 64),
		})
	}
	w.Flush()
}
