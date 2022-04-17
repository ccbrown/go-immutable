# Immutable [![codecov](https://codecov.io/gh/ccbrown/go-immutable/branch/master/graph/badge.svg)](https://codecov.io/gh/ccbrown/go-immutable) [![Documentation](https://godoc.org/github.com/ccbrown/go-immutable?status.svg)](https://godoc.org/github.com/ccbrown/go-immutable)

A collection of fast, general-purpose immutable data structures.

This package has been used in production at multiple companies, is stable, and is probably feature-complete.

## Data Structures

All data structures are fully persistent and safe for concurrent use. Unless otherwise noted, time complexities are worst-case (not amortized).

* Stack: Last in, first out. Constant time operations.
* Queue: First in, first out. Constant time operations.
* Ordered Map: Map with in-order iteration. Logarithmic time operations.
