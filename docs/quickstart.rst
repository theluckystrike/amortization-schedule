Quickstart
==========

Install the module:

.. code-block:: sh

   go get github.com/theluckystrike/amortization-schedule

The package name is ``amortizationschedule``. Rates are annual decimals
(``0.07`` for 7%) and terms are counted in months.

.. code-block:: go

   package main

   import (
       "fmt"

       amortizationschedule "github.com/theluckystrike/amortization-schedule"
   )

   func main() {
       // $100,000 at 7% over 30 years
       pmt := amortizationschedule.MonthlyPayment(100_000.0, 0.07, 360)
       fmt.Printf("%.2f\n", pmt) // 665.30

       sched := amortizationschedule.Schedule(100_000.0, 0.07, 360)
       first := sched[0]
       fmt.Printf("%.2f %.2f\n", first.Interest, first.Principal) // 583.33 81.97

       fmt.Printf("%.4f\n", amortizationschedule.EffectiveAPR(0.07)) // 0.0723
   }

Input handling
--------------

* A principal of zero or less, or a term of zero months, returns ``0`` from
  ``MonthlyPayment`` and ``nil`` from ``Schedule``.
* A zero rate is allowed and reduces to straight-line repayment:
  principal divided by months, with no interest.
* The last row of a schedule has its balance forced to exactly ``0`` so that
  floating-point residue does not show up as a stray fraction of a cent.
