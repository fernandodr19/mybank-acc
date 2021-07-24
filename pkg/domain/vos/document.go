package vos

import "strconv"

type Document int

func (d Document) Int() int {
	return int(d)
}

func (d Document) String() string {
	return strconv.Itoa(d.Int())
}
