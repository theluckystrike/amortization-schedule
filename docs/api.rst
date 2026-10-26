API reference
=============

The authoritative, always-current reference is generated from the source at
https://pkg.go.dev/github.com/theluckystrike/amortization-schedule. This page
summarises it.

.. list-table::
   :header-rows: 1

   * - Function
     - Returns
   * - ``MonthlyPayment(principal, annualRate float64, months uint32) float64``
     - Level payment; ``principal / months`` at a zero rate; ``0`` for invalid input.
   * - ``Schedule(principal, annualRate float64, months uint32) []PaymentRow``
     - One row per month; final balance forced to ``0``; ``nil`` for invalid input.
   * - ``TotalInterest(rows []PaymentRow) float64``
     - Sum of the interest column.
   * - ``TotalPaid(rows []PaymentRow) float64``
     - Sum of the payment column.
   * - ``EffectiveAPR(nominalAnnualRate float64) float64``
     - ``(1 + r/12)^12 - 1``; ``0`` for a rate of -100% or below.
   * - ``RemainingBalanceAfter(principal, annualRate float64, months, periodsPaid uint32) float64``
     - Closed-form balance after ``periodsPaid`` payments.
   * - ``InterestPortion(principal, annualRate float64, months, period uint32) float64``
     - Interest charged in a single 1-based period.
   * - ``MonthsToZero(principal, annualRate float64, months uint32) uint32``
     - Number of periods to amortize; equals the term for a valid loan, ``0`` otherwise.

``PaymentRow``
--------------

.. code-block:: go

   type PaymentRow struct {
       Period    int     // 1-based period number
       Payment   float64 // total scheduled payment this period
       Interest  float64 // interest portion
       Principal float64 // principal portion
       Balance   float64 // remaining balance after this payment
   }
