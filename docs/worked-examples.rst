Worked examples
===============

Every figure on this page was produced by running the package itself
(``MonthlyPayment``, ``Schedule``, ``TotalInterest``, ``RemainingBalanceAfter``
and ``EffectiveAPR``) and rounding to the cent. The loan sizes and rates are
illustrative, not quotes.

Same machine, two terms
-----------------------

An $85,000 equipment loan at 8.25% nominal (8.5692% effective annual rate):

=============  ===============  ================  ====================  ====================
Term           Monthly payment  Total interest    Balance after 12 mo   Balance after 36 mo
=============  ===============  ================  ====================  ====================
60 months      1,733.68         19,020.88         70,674.69             38,236.17
84 months      1,335.44         27,176.96         75,638.49             54,440.11
=============  ===============  ================  ====================  ====================

The longer term lowers the payment by $398.24 a month and costs $8,156.08
more in interest. It also leaves $16,203.94 more owed after three years,
which matters if the machine may be traded in or sold before the note is
paid off.

To compare a finance agreement like this one against a one-dollar buyout lease
and a fair market value lease on the same machine, the
`equipment loan versus lease comparison on AssetLoanCalculator.com
<https://assetloancalculator.com/equipment-loan-calculator/>`_ takes all three
quotes side by side and recomputes each from the figures you enter.

When principal overtakes interest
---------------------------------

On short terms the principal portion exceeds the interest portion from the
first payment. On long terms it does not, and the crossover period is a useful
single number for how front-loaded a loan is:

=================================  ===========  ================  ===============================
Loan                               Payment      Total interest    Principal > interest from period
=================================  ===========  ================  ===============================
$40,000, 7.99%, 180 months         382.03       28,765.39         77
$60,000, 9.50%, 240 months         559.28       74,226.89         154
$85,000, 8.25%, 60 months          1,733.68     19,020.88         1
=================================  ===========  ================  ===============================

On the 20-year loan, the first payment is $475.00 interest and $84.28
principal, and after ten years (120 payments) $43,221.74 of the original
$60,000 is still owed.

.. code-block:: go

   sched := amortizationschedule.Schedule(60_000, 0.095, 240)
   for _, row := range sched {
       if row.Principal > row.Interest {
           fmt.Println(row.Period) // 154
           break
       }
   }

Zero-rate promotions
--------------------

Dealer 0% offers are handled as straight-line repayment. $12,000 over 48
months is $250.00 a month, zero interest, and $6,000.00 owed after 24
payments.
