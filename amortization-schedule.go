// Package amortizationschedule implements fully-amortizing fixed-rate loan math:
// the monthly payment and the full per-period amortization schedule (interest,
// principal, remaining balance). This is the same engine behind the asset loan
// amortization tools at https://assetloancalculator.com/.
//
// A worked instance of the same fixed-rate schedule, with the payment,
// interest, principal and balance columns rendered per period, is at
// https://assetloancalculator.com/rv-loan-calculator/
//
// Quick example:
//
//	pmt := amortizationschedule.MonthlyPayment(100_000.0, 0.07, 360) // ~665.30
//	sched := amortizationschedule.Schedule(100_000.0, 0.07, 360)
//	// sched[0].Interest ~= 583.33, sched[0].Principal ~= 81.97
//	_ = amortizationschedule.TotalInterest(sched)
package amortizationschedule

import "math"

// PaymentRow is one period of an amortization schedule.
type PaymentRow struct {
	Period   int     // 1-based period number
	Payment  float64 // total scheduled payment this period
	Interest float64 // interest portion
	Principal float64 // principal portion
	Balance  float64 // remaining balance after this payment
}

// MonthlyPayment returns the fixed monthly payment on a fully-amortizing
// fixed-rate loan (the standard amortization formula). principal <= 0 or
// months == 0 returns 0. A zero rate reduces to straight-line principal/months.
func MonthlyPayment(principal, annualRate float64, months uint32) float64 {
	if months == 0 || principal <= 0.0 {
		return 0.0
	}
	r := annualRate / 12.0
	n := float64(months)
	if r == 0.0 {
		return principal / n
	}
	pow := math.Pow(1.0+r, n)
	return principal * r * pow / (pow - 1.0)
}

// Schedule returns the full amortization schedule as a slice of PaymentRow, one
// row per period. The final row's Balance is clamped to 0 (no floating-point
// dust). For invalid inputs it returns nil.
func Schedule(principal, annualRate float64, months uint32) []PaymentRow {
	if months == 0 || principal <= 0.0 {
		return nil
	}
	pmt := MonthlyPayment(principal, annualRate, months)
	r := annualRate / 12.0
	bal := principal
	rows := make([]PaymentRow, 0, int(months))
	for i := uint32(1); i <= months; i++ {
		interest := bal * r
		princ := pmt - interest
		// Guard against overpayment on the final period.
		if princ > bal {
			princ = bal
		}
		bal -= princ
		if bal < 0.0 {
			bal = 0.0
		}
		// On the last period, force the balance to exactly 0.
		if i == months {
			bal = 0.0
		}
		row := PaymentRow{
			Period:    int(i),
			Payment:   pmt,
			Interest:  interest,
			Principal: princ,
			Balance:   bal,
		}
		rows = append(rows, row)
	}
	return rows
}

// TotalInterest sums the interest column of a schedule.
func TotalInterest(rows []PaymentRow) float64 {
	var sum float64
	for _, r := range rows {
		sum += r.Interest
	}
	return sum
}

// TotalPaid sums all scheduled payments (payment × periods).
func TotalPaid(rows []PaymentRow) float64 {
	var sum float64
	for _, r := range rows {
		sum += r.Payment
	}
	return sum
}

// EffectiveAPR returns the effective annual rate from a nominal annual rate
// compounded monthly: (1 + r/12)^12 - 1. nominalRate <= -1 returns 0.
func EffectiveAPR(nominalAnnualRate float64) float64 {
	if nominalAnnualRate <= -1.0 {
		return 0.0
	}
	m := nominalAnnualRate / 12.0
	return math.Pow(1.0+m, 12.0) - 1.0
}

// RemainingBalanceAfter returns the remaining loan balance after the given
// number of periods have been paid, using the closed-form amortization balance:
//
//	B(k) = P*(1+r)^k - PMT*(((1+r)^k - 1)/r)
//
// For a zero rate it returns principal - pmt*k. Returns the original principal
// for k <= 0 (nothing paid yet), and 0 for k >= months (loan paid off).
func RemainingBalanceAfter(principal, annualRate float64, months, periodsPaid uint32) float64 {
	if principal <= 0.0 || months == 0 {
		return 0.0
	}
	if periodsPaid == 0 {
		return principal
	}
	if periodsPaid >= months {
		return 0.0
	}
	r := annualRate / 12.0
	pmt := MonthlyPayment(principal, annualRate, months)
	k := float64(periodsPaid)
	if r == 0.0 {
		bal := principal - pmt*k
		if bal < 0.0 {
			return 0.0
		}
		return bal
	}
	growth := math.Pow(1.0+r, k)
	bal := principal*growth - pmt*((growth-1.0)/r)
	if bal < 0.0 {
		return 0.0
	}
	return bal
}

// InterestPortion returns the interest portion of the payment in period `n`
// (1-based) without building the whole schedule. It is the balance after n-1
// periods times the monthly rate.
func InterestPortion(principal, annualRate float64, months, period uint32) float64 {
	if period == 0 || period > months || principal <= 0.0 {
		return 0.0
	}
	balBefore := RemainingBalanceAfter(principal, annualRate, months, period-1)
	return balBefore * (annualRate / 12.0)
}

// MonthsToZero returns the number of full monthly periods required to fully
// amortize a loan — by definition equal to `months` for a valid loan. Provided
// for completeness; returns 0 for invalid input.
func MonthsToZero(principal, annualRate float64, months uint32) uint32 {
	if MonthlyPayment(principal, annualRate, months) <= 0.0 {
		return 0
	}
	return months
}
