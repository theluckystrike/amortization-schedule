amortization-schedule
=====================

``amortization-schedule`` is a small, dependency-free Go package for
fully-amortizing fixed-rate loans. It computes the level monthly payment, the
full per-period schedule (payment, interest, principal and remaining balance),
totals, the effective annual rate, and the closed-form balance after any number
of payments.

It is the calculation engine behind the asset loan calculators at
`AssetLoanCalculator.com <https://assetloancalculator.com/>`_, published as a
standalone MIT-licensed package so the same arithmetic can be checked and
reused outside the site.

.. code-block:: sh

   go get github.com/theluckystrike/amortization-schedule

.. toctree::
   :maxdepth: 2
   :caption: Contents

   quickstart
   formulas
   worked-examples
   api

* Source: https://github.com/theluckystrike/amortization-schedule
* Go reference: https://pkg.go.dev/github.com/theluckystrike/amortization-schedule
