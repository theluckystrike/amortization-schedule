package amortizationschedule

import (
	"math"
	"testing"
)

func approx(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestMonthlyPayment(t *testing.T) {
	// $100k, 7%, 30y -> ~665.30
	m := MonthlyPayment(100_000.0, 0.07, 360)
	if !approx(m, 665.30, 0.01) {
		t.Fatalf("MonthlyPayment = %v, want ~665.30", m)
	}
}

func TestMonthlyPaymentZeroInterest(t *testing.T) {
	// 0% -> straight line
	m := MonthlyPayment(120_000.0, 0.0, 360)
	if !approx(m, 120_000.0/360.0, 1e-9) {
		t.Fatalf("0%% payment = %v, want %v", m, 120_000.0/360.0)
	}
}

func TestMonthlyPaymentGuards(t *testing.T) {
	if MonthlyPayment(0.0, 0.07, 360) != 0.0 {
		t.Fatal("zero principal -> 0")
	}
	if MonthlyPayment(100_000.0, 0.07, 0) != 0.0 {
		t.Fatal("zero months -> 0")
	}
}

func TestScheduleLength(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	if len(s) != 360 {
		t.Fatalf("len = %d, want 360", len(s))
	}
	// Periods are 1-based and sequential.
	if s[0].Period != 1 || s[359].Period != 360 {
		t.Fatalf("periods wrong: %d..%d", s[0].Period, s[359].Period)
	}
}

func TestScheduleFirstRow(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	// First payment: interest = 100000 * 0.07/12 = 583.33
	if !approx(s[0].Interest, 583.33, 0.01) {
		t.Fatalf("first interest = %v, want ~583.33", s[0].Interest)
	}
	// Principal = payment - interest ~ 665.30 - 583.33 = 81.97
	if !approx(s[0].Principal, 81.97, 0.01) {
		t.Fatalf("first principal = %v, want ~81.97", s[0].Principal)
	}
}

func TestScheduleFinalBalanceZero(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	if s[359].Balance != 0.0 {
		t.Fatalf("final balance = %v, want exactly 0.0", s[359].Balance)
	}
}

func TestScheduleInterestPrincipalSumToPayment(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	pmt := MonthlyPayment(100_000.0, 0.07, 360)
	for i, r := range s {
		if !approx(r.Interest+r.Principal, pmt, 0.01) {
			t.Fatalf("row %d: interest+principal = %v, want %v", i, r.Interest+r.Principal, pmt)
		}
	}
}

func TestScheduleBalanceDecreases(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	for i := 1; i < len(s); i++ {
		if s[i].Balance > s[i-1].Balance+1e-6 {
			t.Fatalf("balance increased at row %d: %v > %v", i, s[i].Balance, s[i-1].Balance)
		}
	}
}

func TestTotalInterest(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	// Total cost of a 100k 7% 30y loan is ~139,509 in interest.
	ti := TotalInterest(s)
	if !approx(ti, 139_508.0, 5.0) {
		t.Fatalf("TotalInterest = %v, want ~139508", ti)
	}
}

func TestTotalPaid(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	tp := TotalPaid(s)
	// payment * 360 ~ 665.30 * 360 = 239508
	if !approx(tp, 239_508.0, 5.0) {
		t.Fatalf("TotalPaid = %v, want ~239508", tp)
	}
}

func TestEffectiveAPR(t *testing.T) {
	// 12% nominal monthly -> (1.01)^12 - 1 = 0.12682
	apr := EffectiveAPR(0.12)
	if !approx(apr, 0.126825, 1e-5) {
		t.Fatalf("EffectiveAPR(0.12) = %v, want ~0.12682", apr)
	}
	// 0% -> 0
	if EffectiveAPR(0.0) != 0.0 {
		t.Fatal("EffectiveAPR(0) should be 0")
	}
}

func TestRemainingBalanceAfter(t *testing.T) {
	// After 0 periods -> full principal
	if !approx(RemainingBalanceAfter(100_000.0, 0.07, 360, 0), 100_000.0, 1e-9) {
		t.Fatal("balance after 0 periods should equal principal")
	}
	// After all periods -> 0
	if RemainingBalanceAfter(100_000.0, 0.07, 360, 360) != 0.0 {
		t.Fatal("balance after all periods should be 0")
	}
	// Mid-loan: after 1 period balance ~ 100000 - 81.97 = 99918
	bal := RemainingBalanceAfter(100_000.0, 0.07, 360, 1)
	if !approx(bal, 99_918.03, 0.05) {
		t.Fatalf("balance after 1 period = %v, want ~99918.03", bal)
	}
}

func TestRemainingBalanceMatchesSchedule(t *testing.T) {
	s := Schedule(100_000.0, 0.07, 360)
	// Check a few periods match the schedule balance.
	for _, n := range []uint32{1, 12, 60, 180, 359} {
		closed := RemainingBalanceAfter(100_000.0, 0.07, 360, n)
		if !approx(closed, s[n-1].Balance, 0.02) {
			t.Fatalf("period %d: closed-form %v vs schedule %v", n, closed, s[n-1].Balance)
		}
	}
}

func TestInterestPortionFirstPeriod(t *testing.T) {
	// First period interest = principal * monthly rate.
	ip := InterestPortion(100_000.0, 0.07, 360, 1)
	if !approx(ip, 583.33, 0.01) {
		t.Fatalf("InterestPortion(1) = %v, want ~583.33", ip)
	}
}

func TestInterestPortionGuards(t *testing.T) {
	if InterestPortion(100_000.0, 0.07, 360, 0) != 0.0 {
		t.Fatal("period 0 -> 0")
	}
	if InterestPortion(100_000.0, 0.07, 360, 400) != 0.0 {
		t.Fatal("period beyond term -> 0")
	}
}

func TestMonthsToZero(t *testing.T) {
	if MonthsToZero(100_000.0, 0.07, 360) != 360 {
		t.Fatal("MonthsToZero should equal term")
	}
	if MonthsToZero(0.0, 0.07, 360) != 0 {
		t.Fatal("invalid loan -> 0")
	}
}

func TestZeroRateSchedule(t *testing.T) {
	s := Schedule(120_000.0, 0.0, 360)
	// All interest 0, principal constant = 333.33, balance drops linearly.
	if s[0].Interest != 0.0 {
		t.Fatalf("0%% interest should be 0, got %v", s[0].Interest)
	}
	if !approx(s[0].Principal, 333.3333, 1e-3) {
		t.Fatalf("0%% principal = %v, want ~333.33", s[0].Principal)
	}
	if s[359].Balance != 0.0 {
		t.Fatalf("0%% final balance = %v, want 0", s[359].Balance)
	}
}

func TestScheduleNilOnInvalid(t *testing.T) {
	if Schedule(0.0, 0.07, 360) != nil {
		t.Fatal("zero principal -> nil schedule")
	}
	if Schedule(100_000.0, 0.07, 0) != nil {
		t.Fatal("zero months -> nil schedule")
	}
}
