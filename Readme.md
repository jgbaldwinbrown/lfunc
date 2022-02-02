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

## Creating your own compatible functions

You can easily wrap most mathematical functions so they can be used by this
library. You simply need to write a function that takes a set of `lfunc.Lfunc`
values, and returns a new `lfunc.Lfunc` value. `lfunc.Lfunc` is a function / closure
with the following signature:

```go
type Lfunc func(precision *big.Rat) Lret
```

And Lret is a struct with the following definition:

```go
type Lret struct {
	Min *big.Rat
	Max *big.Rat
}
```

For example, here is the source for the `lfunc.Add` function, annotated:

```go
func Add(x, y Lfunc) Lfunc {
	// We are taking in two Lfunc values, x and y, and returning a function
	// closure with the lfunc.Lfunc signature

	return func(precision *big.Rat) Lret {
		// min and max will be the return values of the closure
		min := new(big.Rat)
		max := new(big.Rat)

		// Here, we get the concrete minimum and maximum values of our
		// inputs (x and y), using the specified precision. "Unwrap"
		// gets the minimum and maximum values from the function's
		// output.
		xmin, xmax := Unwrap(x(precision))
		ymin, ymax := Unwrap(y(precision))

		// Since Add is an exact function, we generate our new minimum
		// values by adding the minimums of the inputs. If Add were an
		// approximate function, we would need to use the minimum of
		// the Add function as our output minimum.
		min.Add(xmin, ymin)
		// Same idea here, but for maximums.
		max.Add(xmax, ymax)

		return Lret{min, max}
	}
}
```
