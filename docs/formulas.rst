Formulas
========

Every function in the package is a direct implementation of one of the
standard fixed-rate amortization identities below. ``P`` is the principal,
``r`` the monthly rate (annual rate divided by 12) and ``n`` the number of
monthly payments.

Level payment
-------------

.. math::

   \text{PMT} = P \cdot \frac{r (1 + r)^n}{(1 + r)^n - 1}

When ``r = 0`` the fraction is undefined, so the package returns ``P / n``.
This is ``MonthlyPayment``.

One period of the schedule
--------------------------

Each month, interest is charged on the balance carried in, and whatever is
left of the payment reduces the balance:

.. math::

   I_k = B_{k-1} \cdot r, \qquad
   \text{principal}_k = \text{PMT} - I_k, \qquad
   B_k = B_{k-1} - \text{principal}_k

with :math:`B_0 = P`. ``Schedule`` iterates this ``n`` times.
``InterestPortion`` returns :math:`I_k` for a single ``k`` without building
the whole table.

Balance after k payments
------------------------

The recurrence above has a closed form, used by ``RemainingBalanceAfter``:

.. math::

   B_k = P (1 + r)^k - \text{PMT} \cdot \frac{(1 + r)^k - 1}{r}

The first term is what the principal would grow to with no payments; the
second is the future value of the payments made so far. It agrees with the
iterated schedule to within floating-point rounding.

Effective annual rate
---------------------

A nominal annual rate compounded monthly is equivalent to an annual rate of

.. math::

   (1 + \text{annual}/12)^{12} - 1

which is what ``EffectiveAPR`` returns. At 7% nominal this is about 7.23%.
Note that this is the compounding effect of the note rate only; the APR in a
lender's disclosure also folds in certain fees, which this package does not
model.
