package vos

import "strconv"

// Money represents a monetary amount
type Money int

// Float64 converts Money to float64
func (m Money) Float64() float64 {
	return float64(m) / 100
}

// Int converts Money to int
func (m Money) Int() int {
	return int(m)
}

// Int64 converts Money to int64
func (m Money) Int64() int64 {
	return int64(m)
}

// String converts Money to string
func (m Money) String() string {
	return strconv.FormatInt(m.Int64(), 10)
}
