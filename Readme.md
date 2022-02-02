# Lfunc

## Introduction

Lfunc is a library for doing arbitrary-precision, lazy arithmetic. In simple
terms, it combines simple operations like add, subtract, square root, and so
on, into a big function that takes one variable -- precision.

## Example usage

```go
package main

import (
	"fmt"
	"github.com/jgbaldwinbrown/lfunc"
)

func main() {
	// Our input variables are converted into functions that accept a precision
	two := lfunc.Lint64(2)
	thf := lfunc.Lint64(35)

	// Generate a function that calculates the square root of two
	sqrt2plus35 := lfunc.Add(lfunc.Sqrt(two), thf)

	// Calculate the answer to our function at low precision
	lowprec := sqrt2plus35(lfunc.NewPrec(10))

	// Extract the minimum and maximum values of our output (based on the precision specified)
	min, max := lfunc.Unwrap(lowprec)

	//min and max are *big.Rat values, so we can convert them to arbitrary-sized decimal strings for output.
	fmt.Println(min.FloatString(100), max.FloatString(100))

	// Extract the midpoint value and precision of our output -- this does not require recalculating the result.
	mid, prec := lfunc.MidPrec(lowprec)
	fmt.Println(mid.FloatString(100), prec.FloatString(100))

	// Re-run our calculation with higher precision, then print the new range
	hiprec := sqrt2plus35(lfunc.NewPrec(100))
	min, max = lfunc.Unwrap(hiprec)
	fmt.Println(min.FloatString(100), max.FloatString(100))
}
```

## Installation

```sh
go get github.com/jgbaldwinbrown/lfunc
```
