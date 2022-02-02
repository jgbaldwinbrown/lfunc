package lfunc

import (
	"math/big"
)

type Lfunc func(precision *big.Rat) Lret

type Lret struct {
	Min *big.Rat
	Max *big.Rat
}

func Lint64(i int64) Lfunc {
	return func(precision *big.Rat) Lret {
		out := big.NewRat(i, 1)
		return Lret{out, out}
	}
}

func Add(x, y Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		min := new(big.Rat)
		max := new(big.Rat)
		xmin, xmax := Unwrap(x(precision))
		ymin, ymax := Unwrap(y(precision))
		min.Add(xmin, ymin)
		max.Add(xmax, ymax)
		return Lret{min, max}
	}
}

func Sub(x, y Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		min := new(big.Rat)
		max := new(big.Rat)
		xmin, xmax := Unwrap(x(precision))
		ymin, ymax := Unwrap(y(precision))
		min.Sub(xmin, ymin)
		max.Sub(xmax, ymax)
		return Lret{min, max}
	}
}

func Mul(x, y Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		min := new(big.Rat)
		max := new(big.Rat)
		xmin, xmax := Unwrap(x(precision))
		ymin, ymax := Unwrap(y(precision))
		min.Mul(xmin, ymin)
		max.Mul(xmax, ymax)
		return Lret{min, max}
	}
}

func Quo(x, y Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		min := new(big.Rat)
		max := new(big.Rat)
		xmin, xmax := Unwrap(x(precision))
		ymin, ymax := Unwrap(y(precision))
		min.Quo(xmin, ymin)
		max.Quo(xmax, ymax)
		return Lret{min, max}
	}
}

func avg(a *big.Rat, vals ...*big.Rat) {
	a.SetInt64(0)
	for _, val := range vals {
		a.Add(a, val)
	}
	count := big.NewRat(int64(len(vals)), 1)
	a.Quo(a, count)
}

func Avg(vals ...Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		min := new(big.Rat)
		max := new(big.Rat)
		mins := []*big.Rat{}
		maxs := []*big.Rat{}
		for _, val := range vals {
			vmin, vmax := Unwrap(val(precision))
			mins = append(mins, vmin)
			maxs = append(maxs, vmax)
		}
		avg(min, mins...)
		avg(max, maxs...)
		return Lret{min, max}
	}
}

func precisionCmp(min, max, precision *big.Rat) bool {
	temp := new(big.Rat).Sub(max, min)
	return temp.Cmp(precision) > 0
}

func sqrt(min, max, x *big.Rat, precision *big.Rat) {
	min.SetInt64(1)
	max.Set(x)
	mid := big.NewRat(0,1)
	minsq := new(big.Rat).Mul(min, min)
	midsq := new(big.Rat).Mul(mid, mid)
	maxsq := new(big.Rat).Mul(max, max)
	for avg(mid, min, max); precisionCmp(minsq.Mul(min, min), maxsq.Mul(max, max), precision); avg(mid, min, max) {
		if midsq.Mul(mid, mid).Cmp(x) < 0 {
			min.Set(mid)
		} else {
			max.Set(mid)
		}
	}
	return
}

func Sqrt(x Lfunc) Lfunc {
	return func(precision *big.Rat) Lret {
		minmin := new(big.Rat)
		minmax := new(big.Rat)
		maxmin := new(big.Rat)
		maxmax := new(big.Rat)
		xmin, xmax := Unwrap(x(precision))
		sqrt(minmin, minmax, xmin, precision)
		sqrt(maxmin, maxmax, xmax, precision)
		return Lret{minmin, maxmax}
	}
}

func Unwrap(l Lret) (*big.Rat, *big.Rat) {
	return l.Min, l.Max
}

func MidPrec(l Lret) (mid *big.Rat, prec *big.Rat) {
	mid = new(big.Rat)
	avg(mid, l.Min, l.Max)
	prec = new(big.Rat).Sub(l.Max, l.Min)
	return
}

func NewPrec(decplaces int64) *big.Rat {
	out := big.NewRat(1,1)
	tenth := big.NewRat(1,10)
	var i int64
	for i=0; i<decplaces; i++ {
		out.Mul(out, tenth)
	}
	return out
}
