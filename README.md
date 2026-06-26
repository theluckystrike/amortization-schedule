# amortization-schedule

[![Go Reference](https://pkg.go.dev/badge/github.com/theluckystrike/amortization-schedule.svg)](https://pkg.go.dev/github.com/theluckystrike/amortization-schedule)

Fully-amortizing fixed-rate loan math: the monthly payment and the full
per-period amortization schedule (interest, principal, remaining balance).
Pure, dependency-free Go — the same engine behind the asset-loan
[amortization calculator](https://assetloancalculator.com/) at
AssetLoanCalculator.com.

## Install

```sh
go get github.com/theluckystrike/amortization-schedule
```

## Quick example

```go
package main

import (
	"fmt"

	amortizationschedule "github.com/theluckystrike/amortization-schedule"
)

func main() {
	// $100k, 7%, 30-year loan
	pmt := amortizationschedule.MonthlyPayment(100_000.0, 0.07, 360)
	fmt.Println(pmt) // ~665.30

	sched := amortizationschedule.Schedule(100_000.0, 0.07, 360)
	fmt.Println(sched[0])              // period 1: interest ~583.33, principal ~81.97
	fmt.Println(amortizationschedule.TotalInterest(sched)) // ~139,508
	fmt.Println(amortizationschedule.EffectiveAPR(0.07))   // ~0.0723
}
```

## API

| Function | Description |
| --- | --- |
| `MonthlyPayment(principal, annualRate float64, months uint32) float64` | Fixed fully-amortizing payment (0% rate → straight-line). |
| `Schedule(principal, annualRate float64, months uint32) []PaymentRow` | Full per-period schedule; final balance clamped to 0. |
| `TotalInterest(rows []PaymentRow) float64` | Sum of the interest column. |
| `TotalPaid(rows []PaymentRow) float64` | Sum of all scheduled payments. |
| `EffectiveAPR(nominalAnnualRate float64) float64` | (1 + r/12)^12 − 1. |
| `RemainingBalanceAfter(principal, annualRate float64, months, periodsPaid uint32) float64` | Closed-form balance after k periods. |
| `InterestPortion(principal, annualRate float64, months, period uint32) float64` | Interest in a single period without building the schedule. |
| `MonthsToZero(...) uint32` | Periods to fully amortize (= term for a valid loan). |

`PaymentRow` holds `Period`, `Payment`, `Interest`, `Principal`, `Balance`.

## License

MIT.

## Links

- Asset loan amortization calculator: <https://assetloancalculator.com/>
- Package docs: <https://pkg.go.dev/github.com/theluckystrike/amortization-schedule>
